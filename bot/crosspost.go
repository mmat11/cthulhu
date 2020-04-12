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
		hashTags   map[string]struct{} = make(map[string]struct{})
		originID   int64               = updateReq.Message.Chat.ID
		authorID   int                 = updateReq.Message.From.ID
		originName string
		text       string
	)

	if !s.Config.isModerator(originID, authorID) && !s.Config.isAdmin(originID, authorID) {
		return nil
	}

	if updateReq.Message.Chat.UserName != "" {
		originName = updateReq.Message.Chat.UserName
		text = fmt.Sprintf("@%s >", originName)
	} else {
		originName = updateReq.Message.Chat.Title
		text = fmt.Sprintf("%s >", originName)
	}

	if updateReq.Message.ReplyToMessage == nil {
		text += fmt.Sprintf(" %s", updateReq.Message.Text)
	} else {
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
				if g.Group.ID != originID {
					level.Info(s.Logger).Log("msg", "crossposting", "text", text, "to", g.Group.ID, "from", updateReq.Message.Chat.UserName)
					s.Telegram.SendMessage(ctx, g.Group.ID, text)
				}
			}
		}
	}
	return nil
}
