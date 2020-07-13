package store

import (
	"context"
	"errors"
)

var (
	ErrKeyAlreadyExists = errors.New("key already exists")
	ErrKeyNotFound      = errors.New("key not found")
)

type Service interface {
	Create(ctx context.Context, key string, value []byte) error
	Read(ctx context.Context, key string) ([]byte, error)
	Update(ctx context.Context, key string, value []byte) error
	Delete(ctx context.Context, key string) ([]byte, error)
	GetAll(ctx context.Context) map[string][]byte
	Close()
}
