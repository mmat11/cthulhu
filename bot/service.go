package bot

import (
	"context"

	"github.com/go-kit/kit/log"

	"tg.bot/telegram"
)

type Token string

type Service struct {
	Logger log.Logger
	Token  Token
}

func NewService(logger log.Logger, token Token) *Service {
	return &Service{
		Logger: logger,
		Token:  token,
	}
}

func (s *Service) Update(ctx context.Context, req *telegram.Update) error {
	return nil
}
