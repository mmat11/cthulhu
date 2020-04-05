package transport

import (
	"context"
	"net/http"
)

func decodeUpdateRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return nil, nil
}

func encodeUpdateResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return nil
}
