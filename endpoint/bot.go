package endpoint

import (
	"context"

	"tg.bot/telegram"
)

type BotService interface {
	Update(ctx context.Context, update *telegram.Update) error
}
