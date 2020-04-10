package task

import (
	"context"
	"strconv"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"

	"cthulhu/bot"
	"cthulhu/store"
)

func init() {
	Register("StoreCleanup", StoreCleanup)
}

func StoreCleanup(
	ctx context.Context,
	logger log.Logger,
	config bot.Config,
	store store.Service,
	args bot.TaskArgs) func() {
	return func() {
		level.Info(logger).Log("msg", "running task")

		var olderThan int = 3600 // defaults to 1h

		for _, arg := range args {
			if arg.Arg.Name == "olderThan" {
				olderThan, _ = strconv.Atoi(arg.Arg.Value)
			}
		}
		for uID, update := range store.GetAll(ctx) {
			if update.Message.Date < int(time.Now().Unix())-olderThan {
				if _, err := store.Delete(ctx, uID); err != nil {
					level.Error(logger).Log("msg", "error deleting key", "key", uID)
				}
			}
		}
	}
}
