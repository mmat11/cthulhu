package task

import (
	"context"

	"github.com/go-kit/kit/log"

	"cthulhu/bot"
	"cthulhu/metrics"
	"cthulhu/store"
	"cthulhu/telegram"
)

type Task func(
	ctx context.Context,
	logger log.Logger,
	config bot.Config,
	storeSvc store.Service,
	tgSvc telegram.Service,
	metricsSvc metrics.Service,
	args bot.TaskArgs,
) func()

var Registry = map[string]Task{}

func Register(name string, t Task) {
	Registry[name] = t
}
