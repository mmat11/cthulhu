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
	Create(ctx context.Context, key string, value interface{}) error
	Read(ctx context.Context, key string) (interface{}, error)
	Update(ctx context.Context, key string, value interface{}) error
	Delete(ctx context.Context, key string) (interface{}, error)
	GetAll(ctx context.Context) map[string]interface{}
}
