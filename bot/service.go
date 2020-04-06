package bot

import (
	"context"
	"fmt"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"

	"tg.bot/telegram"
)

type Service struct {
	Logger log.Logger
	Token  Token
	Config Config
}

func NewService(logger log.Logger, config Config, token Token) *Service {
	return &Service{
		Logger: logger,
		Token:  token,
		Config: config,
	}
}

func (s *Service) Update(ctx context.Context, updateReq *telegram.Update) error {
	if updateReq.Message == nil {
		return nil
	}

	command := updateReq.Message.Command()
	if command == "" {
		return nil
	}

	level.Info(s.Logger).Log("msg", "received new command", "command", command)
	switch command {
	case "ban":
		if updateReq.Message.ReplyToMessage == nil {
			level.Info(s.Logger).Log("msg", "no message quoted")
			return nil
		}

		var (
			chatID   = updateReq.Message.Chat.ID
			authorID = updateReq.Message.From.ID
			userID   = updateReq.Message.ReplyToMessage.From.ID
		)

		level.Info(s.Logger).Log(
			"msg", "received new ban request",
			"chat_id", chatID,
			"author_id", authorID,
			"user_id", userID,
		)

		if !s.Config.CheckAdminPermissions(chatID, authorID, command) {
			level.Info(s.Logger).Log("msg", "not enough privileges")
			return nil
		}
		if err := telegram.KickChatMember(ctx, string(s.Token), chatID, userID); err != nil {
			return err
		}
		telegram.SendMessage(ctx, string(s.Token), chatID, fmt.Sprintf("user %s banned", updateReq.Message.ReplyToMessage.From.UserName))
		return nil
	}
	return nil
}
