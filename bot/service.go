package bot

import (
	"context"

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

	if command := updateReq.Message.Command(); command != "" {
		level.Info(s.Logger).Log("msg", "received new command", "command", command)
		switch command {
		case banCommand:
			return s.handleBan(ctx, updateReq)
		case unbanCommand:
			return s.handleUnban(ctx, updateReq)
		}
	}

	return s.handleCrossposts(ctx, updateReq)
}
