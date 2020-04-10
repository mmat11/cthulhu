package store

import (
	"context"
	"sync"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"

	"cthulhu/telegram"
)

type service struct {
	KV     map[string]*telegram.Update
	mu     sync.Mutex
	Logger log.Logger
}

func NewInMemory(logger log.Logger) *service {
	return &service{KV: make(map[string]*telegram.Update), Logger: logger}
}

func (s *service) Create(ctx context.Context, key string, value *telegram.Update) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	level.Info(s.Logger).Log("msg", "store: create", "key", key)

	if _, ok := s.KV[key]; ok {
		return ErrKeyAlreadyExists
	}
	s.KV[key] = value
	return nil
}

func (s *service) Read(ctx context.Context, key string) (*telegram.Update, error) {
	if v, ok := s.KV[key]; ok {
		return v, nil
	}
	return nil, ErrKeyNotFound
}

func (s *service) Update(ctx context.Context, key string, value *telegram.Update) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	level.Info(s.Logger).Log("msg", "store: update", "key", key)

	if _, ok := s.KV[key]; ok {
		s.KV[key] = value
		return nil
	}
	return ErrKeyNotFound
}

func (s *service) Delete(ctx context.Context, key string) (*telegram.Update, error) {
	var v *telegram.Update

	s.mu.Lock()
	defer s.mu.Unlock()

	level.Info(s.Logger).Log("msg", "store: delete", "key", key)

	if _, ok := s.KV[key]; !ok {
		return nil, ErrKeyNotFound
	}
	v = s.KV[key]
	delete(s.KV, key)
	return v, nil
}

func (s *service) GetAll(ctx context.Context) map[string]*telegram.Update {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.KV
}
