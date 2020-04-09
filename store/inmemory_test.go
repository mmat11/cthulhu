package store

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"

	"cthulhu/mock"
	"cthulhu/telegram"
)

func TestCreateOk(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var (
		ctx = context.Background()
		key = "test"
		val = &telegram.Update{
			Message: &telegram.Message{
				Text: "test",
			},
		}
	)

	svc := mock.NewStoreService(ctrl)
	svc.
		EXPECT().
		Create(ctx, key, val).
		Return(nil)

	if err := svc.Create(ctx, key, val); err != nil {
		t.Fatal(err)
	}
}

func TestCreateAlreadyExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var (
		ctx = context.Background()
		key = "test"
		val = &telegram.Update{
			Message: &telegram.Message{
				Text: "test",
			},
		}
	)

	svc := mock.NewStoreService(ctrl)
	gomock.InOrder(
		svc.
			EXPECT().
			Create(ctx, key, val).
			Return(nil),
		svc.
			EXPECT().
			Create(ctx, key, val).
			Return(ErrKeyAlreadyExists),
	)

	if err := svc.Create(ctx, key, val); err != nil {
		t.Fatal(err)
	}
	if err := svc.Create(ctx, key, val); err == nil {
		t.Fatal("expected error")
	}
}

func TestReadOk(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var (
		ctx = context.Background()
		key = "test"
		val = &telegram.Update{
			Message: &telegram.Message{
				Text: "test",
			},
		}
	)

	svc := mock.NewStoreService(ctrl)
	gomock.InOrder(
		svc.
			EXPECT().
			Create(ctx, key, val).
			Return(nil),
		svc.
			EXPECT().
			Read(ctx, key).
			Return(val, nil),
	)

	if err := svc.Create(ctx, key, val); err != nil {
		t.Fatal(err)
	}
	if v, err := svc.Read(ctx, key); err != nil {
		t.Fatal(err)
	} else {
		if !cmp.Equal(v, val) {
			t.Fatalf("expected %v, got %v", val, v)
		}
	}
}

func TestReadNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var (
		ctx = context.Background()
		key = "test"
	)

	svc := mock.NewStoreService(ctrl)
	svc.
		EXPECT().
		Read(ctx, key).
		Return(nil, ErrKeyNotFound)

	if _, err := svc.Read(ctx, key); err != ErrKeyNotFound {
		t.Fatal(err)
	}
}

func TestUpdateOk(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var (
		ctx = context.Background()
		key = "test"
		val = &telegram.Update{
			Message: &telegram.Message{
				Text: "test",
			},
		}
		newVal = &telegram.Update{
			Message: &telegram.Message{
				Text: "new",
			},
		}
	)

	svc := mock.NewStoreService(ctrl)
	gomock.InOrder(
		svc.
			EXPECT().
			Create(ctx, key, val).
			Return(nil),
		svc.
			EXPECT().
			Update(ctx, key, newVal).
			Return(nil),
		svc.
			EXPECT().
			Read(ctx, key).
			Return(newVal, nil),
	)

	if err := svc.Create(ctx, key, val); err != nil {
		t.Fatal(err)
	}
	if err := svc.Update(ctx, key, newVal); err != nil {
		t.Fatal(err)
	}
	if v, err := svc.Read(ctx, key); err != nil {
		t.Fatal(err)
	} else {
		if !cmp.Equal(v, newVal) {
			t.Fatalf("expected %v, got %v", newVal, v)
		}
	}
}
func TestUpdateNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var (
		ctx    = context.Background()
		key    = "test"
		newVal = &telegram.Update{
			Message: &telegram.Message{
				Text: "new",
			},
		}
	)

	svc := mock.NewStoreService(ctrl)
	svc.
		EXPECT().
		Update(ctx, key, newVal).
		Return(ErrKeyNotFound)

	if err := svc.Update(ctx, key, newVal); err != ErrKeyNotFound {
		t.Fatal(err)
	}
}

func TestDeleteOk(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var (
		ctx = context.Background()
		key = "test"
		val = &telegram.Update{
			Message: &telegram.Message{
				Text: "test",
			},
		}
	)

	svc := mock.NewStoreService(ctrl)
	gomock.InOrder(
		svc.
			EXPECT().
			Create(ctx, key, val).
			Return(nil),
		svc.
			EXPECT().
			Delete(ctx, key).
			Return(val, nil),
		svc.
			EXPECT().
			Read(ctx, key).
			Return(nil, ErrKeyNotFound),
	)

	if err := svc.Create(ctx, key, val); err != nil {
		t.Fatal(err)
	}
	if v, err := svc.Delete(ctx, key); err != nil {
		t.Fatal(err)
	} else {
		if !cmp.Equal(v, val) {
			t.Fatalf("expected %v, got %v", val, v)
		}
	}
	if _, err := svc.Read(ctx, key); err != ErrKeyNotFound {
		t.Fatal(err)
	}
}

func TestDeleteNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var (
		ctx = context.Background()
		key = "test"
	)

	svc := mock.NewStoreService(ctrl)
	svc.
		EXPECT().
		Delete(ctx, key).
		Return(nil, ErrKeyNotFound)

	if _, err := svc.Delete(ctx, key); err != ErrKeyNotFound {
		t.Fatal(err)
	}
}
