package bot

import (
	"context"
	"fmt"
	"unicode/utf16"

	"github.com/go-kit/kit/log/level"

	"cthulhu/telegram"
)

func (s *service) handleCrossposts(ctx context.Context, updateReq *telegram.Update) error {
	var (
		hashTags       map[string]struct{} = make(map[string]struct{})
		chatID         int64               = updateReq.Message.Chat.ID
		authorID       int                 = updateReq.Message.From.ID
		chatName, text string
		quoting        bool
	)

	if !s.Config.isMod(authorID) {
		return nil
	}

	if updateReq.Message.Chat.UserName != "" {
		chatName = updateReq.Message.Chat.UserName
		text = fmt.Sprintf("@%s >", chatName)
	} else {
		chatName = updateReq.Message.Chat.Title
		text = fmt.Sprintf("%s >", chatName)
	}

	quoting = false
	if updateReq.Message.ReplyToMessage == nil {
		text += fmt.Sprintf(" %s", updateReq.Message.Text)
	} else {
		quoting = true
		text += fmt.Sprintf(" %s", updateReq.Message.ReplyToMessage.Text)
	}

	if entities := updateReq.Message.Entities; entities != nil {
		for _, entity := range *entities {
			if entity.IsHashtag() {
				offset := entity.Offset
				utf16Text := utf16.Encode([]rune(updateReq.Message.Text))
				utf16HashTag := utf16Text[offset+1 : offset+entity.Length]
				utf8HashTag := string(utf16.Decode(utf16HashTag))
				hashTags[utf8HashTag] = struct{}{}
			}
		}
	}

	for _, g := range s.Config.Bot.AccessControl.Groups {
		for _, hashTag := range g.Group.CrossPostTags {
			if _, ok := hashTags[hashTag]; ok {
				if g.Group.ID != chatID {
					level.Info(s.Logger).Log("msg", "crossposting", "text", text, "to", g.Group.ID, "from", chatName)
					if quoting {
						s.Telegram.Reply(ctx, g.Group.ID, fmt.Sprintf("your message has been forwarded to %s", g.Group.URL), updateReq.Message.ReplyToMessage.MessageID)
						s.Telegram.DeleteMessage(ctx, g.Group.ID, updateReq.Message.MessageID)
					}
					s.Telegram.SendMessage(ctx, g.Group.ID, text)
				}
			}
		}
	}
	return nil
}
