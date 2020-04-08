package endpoint

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"

	"tg.bot/mock"
	"tg.bot/telegram"
)

func TestMakeUpdateEndpointOk(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var (
		ctx = context.Background()
		req = &telegram.Update{
			Message: &telegram.Message{
				Text: "hello",
			},
		}
	)

	svc := mock.NewBotService(ctrl)
	svc.
		EXPECT().
		Update(ctx, req).
		Return(nil)

	endpoint := MakeUpdateEndpoint(svc)

	resp, err := endpoint(ctx, req)
	if err != nil {
		t.Fatal(err)
	}

	response := resp.(Response)
	expectedResponse := OkResponse
	if !cmp.Equal(response, expectedResponse) {
		t.Log(cmp.Diff(expectedResponse, response))
		t.Fatal("unexpected response")
	}
}

func TestMakeUpdateEndpointErr(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var (
		ctx = context.Background()
		req = &telegram.Update{
			Message: &telegram.Message{
				Text: "hello",
			},
		}
	)

	svc := mock.NewBotService(ctrl)
	svc.
		EXPECT().
		Update(ctx, req).
		Return(errors.New("boom"))

	endpoint := MakeUpdateEndpoint(svc)

	resp, err := endpoint(ctx, req)
	if resp != nil {
		t.Fatalf("expected nil resp, got %v", resp)
	}
	if err == nil {
		t.Fatalf("expected error")
	}
}
