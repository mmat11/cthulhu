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

	svc := mock.NewStoreService(ctrl)
	gomock.InOrder(
		svc.
			EXPECT().
			GetAll(ctx).
			Return(allUpdates),
		svc.
			EXPECT().
			Delete(ctx, "1").
			Return(nil, nil),
	)

	StoreCleanupTask(ctx, log.NewNopLogger(), bot.Config{}, svc, taskArgs)()
}
