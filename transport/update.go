package transport

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"tg.bot/endpoint"
	"tg.bot/telegram"
)

func decodeUpdateRequest(_ context.Context, r *http.Request) (interface{}, error) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, errors.New("body is empty")
	}

	var updateReq telegram.Update

	if err := json.Unmarshal(body, &updateReq); err != nil {
		return nil, errors.New("failed unmarshaling request")
	}
	return &updateReq, nil
}

func encodeUpdateResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	resp, ok := response.(endpoint.Response)
	if !ok {
		return errors.New("response is not endpoint.Response")
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(resp.Status)

	json.NewEncoder(w).Encode(map[string]interface{}{
		"code": resp.Status,
	})
	return nil
}
