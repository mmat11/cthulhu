package transport

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/go-cmp/cmp"

	"tg.bot/endpoint"
	"tg.bot/telegram"
)

func TestDecodeUpdateRequest(t *testing.T) {
	tests := []struct {
		name   string
		req    *http.Request
		expRes *telegram.Update
		expErr string
	}{
		{
			name: "test empty request",
			req: &http.Request{
				Body: ioutil.NopCloser(bytes.NewReader([]byte(`{}`))),
			},
			expRes: &telegram.Update{},
			expErr: "",
		},
		{
			name: "test request with message",
			req: &http.Request{
				Body: ioutil.NopCloser(bytes.NewReader([]byte(`{"message":{"text":"hello"}}`))),
			},
			expRes: &telegram.Update{
				Message: &telegram.Message{
					Text: "hello",
				},
			},
			expErr: "",
		},
		{
			name: "test invalid JSON",
			req: &http.Request{
				Body: ioutil.NopCloser(bytes.NewReader([]byte(``))),
			},
			expRes: nil,
			expErr: "unexpected end of JSON input",
		},
	}

	for _, c := range tests {
		t.Run(c.name, func(t *testing.T) {
			res, err := decodeUpdateRequest(context.Background(), c.req)
			if res != nil && !cmp.Equal(res, c.expRes) {
				t.Fatalf("unexpected result: want %v, got %v", c.expRes, res)
			}
			if err != nil && err.Error() != c.expErr {
				t.Fatalf("unexpected error: want %v, got %v", c.expErr, err)
			}
		})
	}
}

func TestEncodeUpdateResponse(t *testing.T) {
	tests := []struct {
		name    string
		resp    interface{}
		expErr  string
		expBody string
	}{
		{
			name:    "test ok response",
			resp:    endpoint.OkResponse,
			expErr:  "",
			expBody: "{\"code\":200}\n",
		},
		{
			name:    "test wrong response type",
			resp:    "wrong type",
			expErr:  "response is not endpoint.Response",
			expBody: "",
		},
	}

	for _, c := range tests {
		t.Run(c.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			if err := encodeUpdateResponse(context.Background(), w, c.resp); err != nil && err.Error() != c.expErr {
				t.Fatalf("unexpected error: want %v, got %v", c.expErr, err)
			}
			resp := w.Result()
			body, _ := ioutil.ReadAll(resp.Body)
			if string(body) != c.expBody {
				t.Fatalf("unexpected body: want %v, got %v", c.expBody, string(body))
			}
		})
	}
}
