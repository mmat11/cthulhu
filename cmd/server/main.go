package main

import (
	"io/ioutil"
	"net/http"
	"os"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"gopkg.in/yaml.v2"

	"cthulhu/bot"
	"cthulhu/cmd"
	"cthulhu/endpoint"
	"cthulhu/transport"
)

const (
	// server
	listenAddress = ":443"
	certFile      = "./bot.pem"
	keyFile       = "./server.key"

	// bot
	configFile = "config.yaml"
)

var (
	botToken = bot.Token(os.Getenv("BOT_TOKEN"))
)

func main() {
	var (
		logger log.Logger = cmd.MakeLogger()

		config      bot.Config    = readConfig(logger, configFile)
		botService  bot.Service   = bot.NewService(logger, config, botToken)
		endpointSet *endpoint.Set = endpoint.NewSet(botService)

		httpHandler http.Handler = transport.MakeHTTPHandler(
			botService,
			*endpointSet,
			logger,
		)
	)

	level.Info(logger).Log("msg", "start")
	defer level.Info(logger).Log("msg", "stop")

	if err := http.ListenAndServeTLS(listenAddress, certFile, keyFile, httpHandler); err != nil {
		level.Error(logger).Log(
			"msg", "failed to start server",
			"err", err,
		)
		os.Exit(1)
	}
}

func readConfig(logger log.Logger, path string) bot.Config {
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
	return cfg
}
