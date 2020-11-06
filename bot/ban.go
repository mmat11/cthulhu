package bot

import (
	"context"
	"fmt"

	"github.com/go-kit/kit/log/level"

	"cthulhu/telegram"
)

const (
	banCommand   = "ban"
	unbanCommand = "unban"
)

func (s *service) handleBan(ctx context.Context, updateReq *telegram.Update) error {
	if updateReq.Message.ReplyToMessage == nil {
		level.Info(s.Logger).Log("msg", "no message quoted")
		return nil
	}

	var (
		chatID   = updateReq.Message.Chat.ID
		authorID = updateReq.Message.From.ID
		userID   = updateReq.Message.ReplyToMessage.From.ID
		userName = telegram.GetUserName(*updateReq.Message.ReplyToMessage.From)
	)

	level.Info(s.Logger).Log(
		"msg", "received new ban request",
		"chat_id", chatID,
		"author_id", authorID,
		"user_id", userID,
		"user_name", userName,
	)

	if !s.Config.hasPermissions(chatID, authorID, banCommand) {
		level.Info(s.Logger).Log("msg", "not enough privileges")
		return nil
	}
	if err := s.Telegram.KickChatMember(ctx, chatID, userID); err != nil {
		return err
	}
	s.Telegram.SendMessage(ctx, chatID, fmt.Sprintf("user %s banned", userName))
	return nil
}

func (s *service) handleUnban(ctx context.Context, updateReq *telegram.Update) error {
	if updateReq.Message.ReplyToMessage == nil {
		level.Info(s.Logger).Log("msg", "no message quoted")
		return nil
	}

	var (
		chatID   = updateReq.Message.Chat.ID
		authorID = updateReq.Message.From.ID
		userID   = updateReq.Message.ReplyToMessage.From.ID
		userName = telegram.GetUserName(*updateReq.Message.ReplyToMessage.From)
	)

	level.Info(s.Logger).Log(
		"msg", "received new unban request",
		"chat_id", chatID,
		"author_id", authorID,
		"user_id", userID,
		"user_name", userName,
	)

	if !s.Config.hasPermissions(chatID, authorID, unbanCommand) {
		level.Info(s.Logger).Log("msg", "not enough privileges")
		return nil
	}
	if err := s.Telegram.UnbanChatMember(ctx, chatID, userID); err != nil {
		return err
	}
	s.Telegram.SendMessage(ctx, chatID, fmt.Sprintf("user %s unbanned", userName))
	return nil
}
