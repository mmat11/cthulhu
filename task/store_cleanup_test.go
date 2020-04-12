package task

import (
	"context"
	"testing"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/golang/mock/gomock"

	"cthulhu/bot"
	"cthulhu/mock"
	"cthulhu/telegram"
)

func TestStoreCleanupTask(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var (
		ctx      = context.Background()
		taskArgs = bot.TaskArgs{
			{
				Arg: bot.TaskArg{
					Name:  retentionKey,
					Value: "20000",
				},
			},
		}
		allUpdates map[string]interface{} = map[string]interface{}{
			"1": &telegram.Update{
				Message: &telegram.Message{
					Date: 1000000,
				},
			},
			"2": &telegram.Update{
				Message: &telegram.Message{
					Date: int(time.Now().Unix()),
				},
			},
		}
	)

	storeSvc := mock.NewStoreService(ctrl)
	gomock.InOrder(
		storeSvc.
			EXPECT().
			GetAll(ctx).
			Return(allUpdates),
		storeSvc.
			EXPECT().
			Delete(ctx, "1").
			Return(nil, nil),
	)

	telegramSvc := mock.NewTelegramService(ctrl)

	StoreCleanupTask(ctx, log.NewNopLogger(), bot.Config{}, storeSvc, telegramSvc, taskArgs)()
}
