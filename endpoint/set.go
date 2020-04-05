package endpoint

import (
	"github.com/go-kit/kit/endpoint"
)

type Set struct {
	Update endpoint.Endpoint
}

func NewSet(s BotService) *Set {
	return &Set{
		Update: MakeUpdateEndpoint(s),
	}
}
