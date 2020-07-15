package task

import (
	"context"
	"strconv"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"

	"cthulhu/bot"
	"cthulhu/metrics"
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
	tg telegram.Service,
	metrics metrics.Service,
	args bot.TaskArgs) func() {
	return func() {
		level.Info(logger).Log("msg", "running task")

		var (
			startTime     = time.Now()
			retention int = defaultRetention
		)

		for _, arg := range args {
			if arg.Arg.Name == retentionKey {
				retention, _ = strconv.Atoi(arg.Arg.Value)
			}
		}
		for k, v := range store.GetAll(ctx) {
			u, err := telegram.UnmarshalUpdate(v)
			if err == nil {
				if u.Message.Date < int(time.Now().Unix())-retention {
					if _, err := store.Delete(ctx, k); err != nil {
						level.Error(logger).Log("msg", "error deleting key", "key", k, "err", err)
					}
				}
			}
		}

		metrics.ObserveTasksDuration(
			storeCleanupTaskName,
			float64(time.Since(startTime).Seconds()),
		)
	}
}
