package bot

import (
	"context"
	"fmt"
	"strings"
	"unicode/utf16"

	"github.com/go-kit/kit/log/level"

	"cthulhu/telegram"
)

func (s *service) handleCrossposts(ctx context.Context, updateReq *telegram.Update) error {
	var (
		hashTags             map[string]struct{} = make(map[string]struct{})
		chatID               int64               = updateReq.Message.Chat.ID
		authorID             int                 = updateReq.Message.From.ID
		chatName, text       string
		isCrosspost, quoting bool
		fwdText              string = "your message has been forwarded to"
	)

	if !s.Config.isMod(authorID) {
		level.Info(s.Logger).Log("msg", "user is not mod", "chat_id", chatID, "author", authorID)
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
					isCrosspost = true
					s.Telegram.SendMessage(ctx, g.Group.ID, text)
					if quoting {
						fwdText = fmt.Sprintf("%s %s,", fwdText, g.Group.URL)
					}
				}
			}
		}
	}

	if quoting && isCrosspost {
		if err := s.Telegram.DeleteMessage(ctx, chatID, updateReq.Message.MessageID); err != nil {
			level.Info(s.Logger).Log("msg", "crosspost delete", "err", err)
		}
		if err := s.Telegram.Reply(ctx, chatID, strings.TrimRight(fwdText, ","), updateReq.Message.ReplyToMessage.MessageID); err != nil {
			level.Info(s.Logger).Log("msg", "crosspost reply", "err", err)
		}
	}
	return nil
}
