package transport

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/transport"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"

	"tg.bot/bot"
	"tg.bot/endpoint"
)

func MakeHTTPHandler(s bot.Service, e endpoint.Set, logger log.Logger) http.Handler {
	var (
		r                  = mux.NewRouter()
		apiPrefix          = "/v1"
		apiPrefixWithToken = fmt.Sprintf("%s/%s", apiPrefix, s.Token)

		updatePath = fmt.Sprintf("%s/update", apiPrefixWithToken)
	)

	r.Methods("POST").Path(updatePath).Handler(httptransport.NewServer(
		e.Update,
		decodeUpdateRequest,
		encodeUpdateResponse,
		[]httptransport.ServerOption{
			httptransport.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
			httptransport.ServerErrorEncoder(
				func(_ context.Context, err error, w http.ResponseWriter) {
					if err == nil {
						panic("endpoint error is nil")
					}
					w.Header().Set("Content-Type", "application/json; charset=utf-8")
					w.WriteHeader(http.StatusInternalServerError)
					json.NewEncoder(w).Encode(map[string]interface{}{
						"error": err.Error(),
					})
				},
			),
		}...,
	))

	return r
}
