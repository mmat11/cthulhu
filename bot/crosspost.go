package bot

import (
	"context"

	"tg.bot/telegram"
)

func (s *Service) handleCrossposts(ctx context.Context, updateReq *telegram.Update) error {
	var (
		text     string              = updateReq.Message.Text
		hashTags map[string]struct{} = make(map[string]struct{}, 0)
	)

	if entities := updateReq.Message.Entities; entities != nil {
		for _, entity := range *entities {
			if entity.IsHashtag() {
				hashTags[updateReq.Message.Text[entity.Offset:entity.Length]] = struct{}{}
			}
		}
	}

	for _, g := range s.Config.Bot.AccessControl.Groups {
		for _, hashTag := range g.Group.CrossPostTags {
			if _, ok := hashTags[hashTag]; ok {
				telegram.SendMessage(ctx, string(s.Token), g.Group.ID, text)
			}
		}
	}
	return nil
}
