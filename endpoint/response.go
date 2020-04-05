package endpoint

type Response struct {
	Code uint32
}

var (
	OkResponse Response = Response{Code: 200}
)
