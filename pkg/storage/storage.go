package storage

import (
	"bytes"
	"encoding/gob"
	"errors"
	"time"

	"github.com/boltdb/bolt"
)

type KVStore struct {
	db *bolt.DB
}

var (
	ErrNotFound = errors.New("storage: key not found")
	ErrBadValue = errors.New("storage: bad value")

	bucketName = []byte("main")
)

func Open(path string) (*KVStore, error) {
	opts := &bolt.Options{
		Timeout: 50 * time.Millisecond,
	}

	db, err := bolt.Open(path, 0640, opts)

	if err != nil {
		return nil, err
	}

	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(bucketName)
		return err
	})

	if err != nil {
		return nil, err
	}

	return &KVStore{db: db}, nil
}

func (store *KVStore) Close() error {
	return store.db.Close()
}

func (store *KVStore) Get(key string, value interface{}) error {
	return store.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(bucketName)

		if v := bucket.Get([]byte("key")); v != nil {
			d := gob.NewDecoder(bytes.NewReader(v))
			return d.Decode(v)
		}

		return nil
	})

}

func (store *KVStore) Put(key string, value interface{}) error {
	return nil
}

func (store *KVStore) Delete(key string) error {
	return nil
}
