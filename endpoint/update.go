package endpoint

import (
	"context"

	"github.com/go-kit/kit/endpoint"

	"tg.bot/telegram"
)

func MakeUpdateEndpoint(s BotService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(*telegram.Update)
		err := s.Update(ctx, req)
		if err != nil {
			return nil, err
		}
		return OkResponse, nil
	}
}
