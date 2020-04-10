package task

import (
	"context"
	"strconv"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"

	"cthulhu/bot"
	"cthulhu/store"
	"cthulhu/telegram"
)

func init() {
	Register(storeCleanupTaskName, StoreCleanupTask)
}

const (
	storeCleanupTaskName string = "StoreCleanup"
	retentionKey         string = "retention"
)

var defaultRetention = 7200

func StoreCleanupTask(
	ctx context.Context,
	logger log.Logger,
	config bot.Config,
	store store.Service,
	args bot.TaskArgs) func() {
	return func() {
		level.Info(logger).Log("msg", "running task")

		var retention int = defaultRetention

		for _, arg := range args {
			if arg.Arg.Name == retentionKey {
				retention, _ = strconv.Atoi(arg.Arg.Value)
			}
		}
		for k, v := range store.GetAll(ctx) {
			if u, ok := v.(*telegram.Update); ok {
				if u.Message.Date < int(time.Now().Unix())-retention {
					if _, err := store.Delete(ctx, k); err != nil {
						level.Error(logger).Log("msg", "error deleting key", "key", k, "err", err)
					}
				}
			}
		}
	}
}
