package cmd

import (
	"os"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

func MakeLogger() log.Logger {
	var logger log.Logger
	logger = log.NewJSONLogger(log.NewSyncWriter(os.Stdout))
	logger = level.NewFilter(logger, level.AllowInfo())
	logger = level.NewInjector(logger, level.InfoValue())
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)

	return logger
}
