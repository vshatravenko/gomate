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
		Timeout: 1 * time.Second,
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

		if v := bucket.Get([]byte(key)); v != nil {
			d := gob.NewDecoder(bytes.NewReader(v))
			return d.Decode(value)
		}

		return ErrNotFound
	})
}

func (store *KVStore) Put(key string, value interface{}) error {
	if value == nil {
		return ErrBadValue
	}

	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(value); err != nil {
		return err
	}

	return store.db.Update(func(tx *bolt.Tx) error {
		return tx.Bucket(bucketName).Put([]byte(key), buf.Bytes())
	})
}

func (store *KVStore) Delete(key string) error {
	return store.db.Update(func(tx *bolt.Tx) error {
		c := tx.Bucket(bucketName).Cursor()

		if k, _ := c.Seek([]byte(key)); k == nil || string(k) != key {
			return ErrNotFound
		}

		return c.Delete()
	})
}
