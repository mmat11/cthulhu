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
		allUpdates map[string][]byte = map[string][]byte{
			"1": telegram.MarshalUpdate(&telegram.Update{
				Message: &telegram.Message{
					Date: 1000000,
				},
			}),
			"2": telegram.MarshalUpdate(&telegram.Update{
				Message: &telegram.Message{
					Date: int(time.Now().Unix()),
				},
			}),
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
	metricsSvc := mock.NewMetricsService(ctrl)
	metricsSvc.
		EXPECT().
		ObserveTasksDuration(gomock.Eq(storeCleanupTaskName), gomock.Any())

	StoreCleanupTask(ctx, log.NewNopLogger(), bot.Config{}, storeSvc, telegramSvc, metricsSvc, taskArgs)()
}
