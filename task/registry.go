package task

import (
	"context"

	"github.com/go-kit/kit/log"

	"cthulhu/bot"
	"cthulhu/store"
)

type Task func(
	ctx context.Context,
	logger log.Logger,
	config bot.Config,
	store store.Service,
	args bot.TaskArgs,
) func()

var Registry = map[string]Task{}

func Register(name string, t Task) {
	Registry[name] = t
}
