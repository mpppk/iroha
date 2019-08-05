package storage

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/mpppk/iroha/ktkn"

	bolt "github.com/mpppk/bbolt"
	"github.com/pkg/errors"
)

type Bolt struct {
	*None
	db         *bolt.DB
	bucketName string
}

func NewBolt(dbPath string) (Storage, error) {
	db, err := bolt.Open(dbPath, 0600, nil)
	if err != nil {
		return nil, err
	}
	boltStorage := &Bolt{
		db:         db,
		bucketName: "main",
	}
	err = boltStorage.createBucketIfNotExists(boltStorage.bucketName)
	return boltStorage, err
}

func (b *Bolt) Get(ctx context.Context, indices []int) (wordsList [][]*ktkn.Word, ok bool, err error) {
	wordsList = make([][]*ktkn.Word, 0, 10)
	err = b.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(b.bucketName))
		if bucket == nil {
			return fmt.Errorf("failed to retrive bucket(%s)", b.bucketName)
		}
		v := bucket.Get(toStorageKey(indices))
		if v == nil {
			ok = false
			return nil
		}
		ok = true
		return json.Unmarshal(v, &wordsList)
	})
	return
}

func (b *Bolt) Set(ctx context.Context, indices []int, wordsList [][]*ktkn.Word) error {
	wl := wordsList
	if wl == nil {
		wl = make([][]*ktkn.Word, 0)
	}

	wordsListJsonBytes, err := json.Marshal(wl)

	if err != nil {
		return err
	}
	err = b.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(b.bucketName))
		err := b.Put(
			toStorageKey(indices),
			wordsListJsonBytes)
		return errors.Wrapf(err, "failed to put wordsList to bolt DB: indices:%s", indices)
	})
	if err != nil {
		return err
	}
	if err := b.deleteIndexChildren(indices); err != nil {
		return err
	}
	return nil
}

func (b *Bolt) deleteIndexChildren(indices []int) error {
	var deleteKeys [][]byte
	err := b.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(b.bucketName))
		c := bucket.Cursor()

		prefix := toStorageKey(indices)
		strPrefix := string(prefix)
		for k, _ := c.Seek(prefix); k != nil && bytes.HasPrefix(k, prefix); k, _ = c.Next() {
			if string(k) != strPrefix {
				deleteKeys = append(deleteKeys, k)
			}
		}
		return nil
	})
	err = b.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(b.bucketName))
		for _, deleteKey := range deleteKeys {
			if err := bucket.Delete(deleteKey); err != nil {
				return errors.Wrap(err, "failed to delete key: "+string(deleteKey))
			}
		}
		return nil
	})
	if err != nil {
		return errors.Wrapf(err, "failed to delete keys: %s", deleteKeys)
	}
	return nil
}
