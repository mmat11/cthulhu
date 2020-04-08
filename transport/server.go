package transport

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/go-kit/kit/transport"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"

	"cthulhu/bot"
	"cthulhu/endpoint"
)

func MakeHTTPHandler(s bot.Service, e endpoint.Set, logger log.Logger) http.Handler {
	var (
		r                  = mux.NewRouter()
		apiPrefix          = "/v1"
		apiPrefixWithToken = fmt.Sprintf("%s/%s", apiPrefix, s.GetToken())

		updatePath = fmt.Sprintf("%s/update", apiPrefixWithToken)
	)

	level.Info(logger).Log("msg", "add handler", "path", updatePath)

	r.Methods(http.MethodPost).Path(updatePath).Handler(httptransport.NewServer(
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
					errorEncoder(w, err.Error())
				},
			),
		}...,
	))

	return r
}
