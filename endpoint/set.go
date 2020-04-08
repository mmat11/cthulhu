package endpoint

import (
	"github.com/go-kit/kit/endpoint"

	"tg.bot/bot"
)

type Set struct {
	Update endpoint.Endpoint
}

func NewSet(s bot.Service) *Set {
	return &Set{
		Update: MakeUpdateEndpoint(s),
	}
}
