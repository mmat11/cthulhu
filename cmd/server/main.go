package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"

	"tg.bot/bot"
	"tg.bot/cmd"
	"tg.bot/endpoint"
	"tg.bot/transport"
)

const (
	listenAddress = ":3000"
)

var (
	botToken       = bot.Token(os.Getenv("BOT_TOKEN"))
	webhookAddress = os.Getenv("WEBHOOK_ADDRESS")
)

func main() {
	var (
		logger log.Logger = cmd.MakeLogger()

		botService  *bot.Service  = bot.NewService(logger, botToken)
		endpointSet *endpoint.Set = endpoint.NewSet(botService)

		httpHandler http.Handler = transport.MakeHTTPHandler(
			*botService,
			*endpointSet,
			logger,
		)

		errs = make(chan error)
	)

	level.Info(logger).Log("msg", "starting")
	defer level.Info(logger).Log("msg", "stopped", "errs", <-errs)

	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errs <- fmt.Errorf("%s", <-c)
	}()
	go func() {
		errs <- http.ListenAndServe(listenAddress, httpHandler)
	}()
}
