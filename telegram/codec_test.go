package telegram_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"cthulhu/telegram"
)

func TestMarshalUnmarshalUpdate(t *testing.T) {
	var upd *telegram.Update = &telegram.Update{
		UpdateID: 1,
		Message: &telegram.Message{
			MessageID: 2,
		},
	}

	encoded := telegram.MarshalUpdate(upd)
	res, err := telegram.UnmarshalUpdate(encoded)
	if err != nil {
		t.Fatal(err)
	}

	if !cmp.Equal(upd, res) {
		t.Fatalf("error during marshal/unmarshal: want %v, got %v", upd, res)
	}
}
