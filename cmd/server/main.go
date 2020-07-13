package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/robfig/cron/v3"
	"gopkg.in/yaml.v2"

	"cthulhu/bot"
	"cthulhu/cmd"
	"cthulhu/endpoint"
	"cthulhu/metrics"
	"cthulhu/store"
	"cthulhu/task"
	"cthulhu/telegram"
	"cthulhu/transport"
)

const (
	// server
	listenAddress = ":443"
	certFile      = "./bot.pem"
	keyFile       = "./server.key"

	// prometheus exporter
	metricsListenAddress = ":2112"

	// bot
	configFile = "config.yaml"

	// telegram
	apiEndpoint = "https://api.telegram.org"
)

var (
	botToken   = bot.Token(os.Getenv("BOT_TOKEN"))
	badgerPath = os.Getenv("BADGER_PATH")
)

func main() {
	var (
		logger          log.Logger       = cmd.MakeLogger()
		config          bot.Config       = readConfig(logger, configFile, botToken)
		metricsService  metrics.Service  = metrics.NewService()
		telegramService telegram.Service = telegram.NewService(logger, apiEndpoint, string(botToken))
		storeService    store.Service
		botService      bot.Service
		endpointSet     *endpoint.Set
		httpHandler     http.Handler
	)

	switch config.Bot.Database.Type {
	case "badger":
		s, err := store.NewBadger(logger, badgerPath)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		storeService = s
	default:
		storeService = store.NewInMemory(logger)
	}

	botService = bot.NewService(logger, config, storeService, telegramService, metricsService)
	endpointSet = endpoint.NewSet(botService)
	httpHandler = transport.MakeHTTPHandler(
		botService,
		*endpointSet,
		logger,
	)

	level.Info(logger).Log("msg", "start")
	defer level.Info(logger).Log("msg", "stop")

	level.Info(logger).Log("msg", "tasks registered", "tasks", task.Registry)
	go registerTasks(logger, config, storeService, telegramService)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	s := &http.Server{
		Addr:    listenAddress,
		Handler: httpHandler,
	}

	promexp := &http.Server{
		Addr:    metricsListenAddress,
		Handler: promhttp.Handler(),
	}

	// webhook handler
	go func() {
		level.Info(logger).Log("msg", "starting server", "addr", listenAddress)
		if err := s.ListenAndServeTLS(certFile, keyFile); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}()

	// prometheus metrics
	go func() {
		level.Info(logger).Log("msg", "starting metrics server", "addr", metricsListenAddress)
		if err := promexp.ListenAndServe(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}()

	<-c

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	storeService.Close()

	if err := promexp.Shutdown(ctx); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if err := s.Shutdown(ctx); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func registerTasks(logger log.Logger, cfg bot.Config, st store.Service, tg telegram.Service) {
	c := cron.New()
	for _, t := range cfg.Bot.Tasks {
		f, ok := task.Registry[t.Task.Name]
		if !ok {
			level.Error(logger).Log(
				"msg", "task not implemented",
				"task_name", t.Task.Name,
			)
			os.Exit(1)
		}
		c.AddFunc(
			t.Task.Cron,
			f(
				context.Background(),
				log.With(logger, "task", t.Task.Name),
				cfg,
				st,
				tg,
				t.Task.Args,
			),
		)
	}
	c.Start()
}

func readConfig(logger log.Logger, path string, token bot.Token) bot.Config {
	var cfg bot.Config

	data, err := ioutil.ReadFile(path)
	if err != nil {
		level.Error(logger).Log(
			"msg", "failed to open config file",
			"err", err,
		)
		os.Exit(1)
	}

	err = yaml.Unmarshal([]byte(data), &cfg)

	if err != nil {
		level.Error(logger).Log(
			"msg", "config file data is invalid",
			"err", err,
		)
		os.Exit(1)
	}

	cfg.Bot.Token = token

	return cfg
}
