package bot

import (
	"context"

	"github.com/go-kit/kit/log/level"

	"cthulhu/telegram"
)

const broadcastCommand = "broadcast"

func (s *service) handleBroadcast(ctx context.Context, updateReq *telegram.Update) error {
	var (
		chatID   = updateReq.Message.Chat.ID
		authorID = updateReq.Message.From.ID
		message  = updateReq.Message.CommandArguments()
	)

	level.Info(s.Logger).Log(
		"msg", "received new broadcast request",
		"chat_id", chatID,
		"author_id", authorID,
		"message", message,
	)

	if !s.Config.hasPermissions(chatID, authorID, broadcastCommand) {
		level.Info(s.Logger).Log("msg", "not enough privileges")
		return nil
	}

	for _, g := range s.Config.Bot.AccessControl.Groups {
		if g.Group.ID != chatID {
			if message != "" {
				level.Info(s.Logger).Log(
					"msg", "writing",
					"chat_id", g.Group.ID,
					"message", message,
				)
				telegram.SendMessage(ctx, string(s.Token), g.Group.ID, message)
			}
		}
	}
	return nil
}
