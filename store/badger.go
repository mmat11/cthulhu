package store

import (
	"context"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"

	badgerv2 "github.com/dgraph-io/badger/v2"
)

type badger struct {
	db     *badgerv2.DB
	Logger log.Logger
}

func NewBadger(logger log.Logger, badgerPath string) (*badger, error) {
	db, err := badgerv2.Open(badgerv2.DefaultOptions(badgerPath))
	if err != nil {
		return nil, err
	}

	return &badger{
		db:     db,
		Logger: logger,
	}, nil
}

func (b *badger) Create(ctx context.Context, key string, value []byte) error {
	level.Info(b.Logger).Log("msg", "store: create", "key", key)

	var bkey []byte = []byte(key)

	err := b.db.Update(func(txn *badgerv2.Txn) error {
		return txn.Set(bkey, value)
	})
	if err != nil {
		return err
	}
	return nil
}

func (b *badger) Read(ctx context.Context, key string) ([]byte, error) {
	txn := b.db.NewTransaction(true)
	defer txn.Discard()

	var bkey []byte = []byte(key)

	item, err := txn.Get(bkey)
	if err != nil {
		if err == badgerv2.ErrKeyNotFound {
			return nil, ErrKeyNotFound
		}
		return nil, err
	}

	val, err := item.ValueCopy(nil)
	if err != nil {
		return nil, err
	}

	if err := txn.Commit(); err != nil {
		return nil, err
	}

	level.Info(b.Logger).Log("msg", "store: read", "key", key)

	return val, nil
}

func (b *badger) Update(ctx context.Context, key string, value []byte) error {
	level.Info(b.Logger).Log("msg", "store: update", "key", key)

	txn := b.db.NewTransaction(true)
	defer txn.Discard()

	var bkey []byte = []byte(key)

	_, err := txn.Get(bkey)
	if err != nil {
		if err == badgerv2.ErrKeyNotFound {
			return ErrKeyNotFound
		}
		return err
	}

	err = b.db.Update(func(txn *badgerv2.Txn) error {
		return txn.Set(bkey, value)
	})
	if err != nil {
		return err
	}

	if err := txn.Commit(); err != nil {
		return err
	}

	return nil
}

func (b *badger) Delete(ctx context.Context, key string) ([]byte, error) {
	level.Info(b.Logger).Log("msg", "store: delete", "key", key)

	var bkey []byte = []byte(key)

	txn := b.db.NewTransaction(true)
	defer txn.Discard()

	item, err := txn.Get(bkey)
	if err != nil {
		if err == badgerv2.ErrKeyNotFound {
			return nil, ErrKeyNotFound
		}
		return nil, err
	}

	val, err := item.ValueCopy(nil)
	if err != nil {
		return nil, err
	}

	if err := txn.Delete(bkey); err != nil {
		return nil, err
	}

	if err := txn.Commit(); err != nil {
		return nil, err
	}

	return val, nil
}

func (b *badger) GetAll(ctx context.Context) map[string][]byte {
	var kv map[string][]byte = make(map[string][]byte)

	err := b.db.View(func(txn *badgerv2.Txn) error {
		it := txn.NewIterator(badgerv2.DefaultIteratorOptions)
		defer it.Close()
		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			val, err := item.ValueCopy(nil)
			if err != nil {
				return err
			}
			kv[string(item.Key())] = val
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
	return kv
}

func (b *badger) Close() {
	b.db.Close()
}
