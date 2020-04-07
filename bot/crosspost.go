package bot

import (
	"context"
	"fmt"
	"unicode/utf16"

	"tg.bot/telegram"
)

func (s *Service) handleCrossposts(ctx context.Context, updateReq *telegram.Update) error {
	var (
		text     string              = fmt.Sprintf("from: t.me/%s // %s", updateReq.Message.Chat.UserName, updateReq.Message.Text)
		hashTags map[string]struct{} = make(map[string]struct{}, 0)
		originID int64               = updateReq.Message.Chat.ID
	)

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
					telegram.SendMessage(ctx, string(s.Token), g.Group.ID, text)
				}
			}
		}
	}
	return nil
}
