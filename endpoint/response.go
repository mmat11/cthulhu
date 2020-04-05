package endpoint

import "net/http"

type Response struct {
	Status int
}

var (
	OkResponse Response = Response{Status: http.StatusOK}
)
