package store

import (
	"context"
	"errors"

	"cthulhu/telegram"
)

var (
	ErrKeyAlreadyExists = errors.New("key already exists")
	ErrKeyNotFound      = errors.New("key not found")
)

type Service interface {
	Create(ctx context.Context, key string, value *telegram.Update) error
	Read(ctx context.Context, key string) (*telegram.Update, error)
	Update(ctx context.Context, key string, value *telegram.Update) error
	Delete(ctx context.Context, key string) (*telegram.Update, error)
	GetAll(ctx context.Context) map[string]*telegram.Update
}
