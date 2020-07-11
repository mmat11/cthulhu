package store

import (
	"context"
	"sync"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

type inMemory struct {
	KV     map[string][]byte
	mu     sync.Mutex
	Logger log.Logger
}

func NewInMemory(logger log.Logger) *inMemory {
	return &inMemory{
		KV:     make(map[string][]byte),
		Logger: logger,
	}
}

func (s *inMemory) Create(ctx context.Context, key string, value []byte) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	level.Info(s.Logger).Log("msg", "store: create", "key", key)

	if _, ok := s.KV[key]; ok {
		return ErrKeyAlreadyExists
	}
	s.KV[key] = value
	return nil
}

func (s *inMemory) Read(ctx context.Context, key string) ([]byte, error) {
	if v, ok := s.KV[key]; ok {
		return v, nil
	}
	return nil, ErrKeyNotFound
}

func (s *inMemory) Update(ctx context.Context, key string, value []byte) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	level.Info(s.Logger).Log("msg", "store: update", "key", key)

	if _, ok := s.KV[key]; ok {
		s.KV[key] = value
		return nil
	}
	return ErrKeyNotFound
}

func (s *inMemory) Delete(ctx context.Context, key string) ([]byte, error) {
	var v []byte

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

func (s *inMemory) GetAll(ctx context.Context) map[string][]byte {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.KV
}

func (s *inMemory) Close() {
}
