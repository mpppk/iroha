package lib

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/pkg/errors"

	bolt "github.com/mpppk/bbolt"
)

type Storage interface {
	Set(indices []int, wordsList [][]*Word) error
	Get(indices []int) ([][]*Word, bool, error)
}

type NoStorage struct{}

func (e *NoStorage) Get(indices []int) ([][]*Word, bool, error) {
	return nil, false, nil
}

func (e *NoStorage) Set(indices []int, wordsList [][]*Word) error {
	return nil
}

type BoltStorage struct {
	db         *bolt.DB
	bucketName string
}

func NewBoltStorage(dbPath string) (*BoltStorage, error) {
	db, err := bolt.Open(dbPath, 0600, nil)
	if err != nil {
		return nil, err
	}
	boltStorage := &BoltStorage{
		db:         db,
		bucketName: "main",
	}
	err = boltStorage.createBucketIfNotExists(boltStorage.bucketName)
	return boltStorage, err
}

func (b *BoltStorage) Get(indices []int) (wordsList [][]*Word, ok bool, err error) {
	wordsList = make([][]*Word, 0, 10)
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

func (b *BoltStorage) Set(indices []int, wordsList [][]*Word) error {
	wl := wordsList
	if wl == nil {
		wl = make([][]*Word, 0)
	}

	wordsListJsonBytes, err := json.Marshal(wl)

	if err != nil {
		return err
	}
	return b.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(b.bucketName))
		err := b.Put(
			toStorageKey(indices),
			wordsListJsonBytes)
		return errors.Wrapf(err, "failed to put wordsList to bolt DB: indices:%s", indices)
	})
}

func toStorageKey(indices []int) []byte {
	strKey := ""
	if len(indices) == 0 {
		return []byte("no-index")
	}
	for _, index := range indices {
		strKey += strconv.Itoa(index) + ":"
	}
	return []byte(strKey)
}

func (b *BoltStorage) createBucketIfNotExists(bucketName string) error {
	return b.db.Update(func(tx *bolt.Tx) error {
		if _, err := tx.CreateBucketIfNotExists([]byte(bucketName)); err != nil {
			return fmt.Errorf("failed to create %s bucket: %s", bucketName, err)
		}
		return nil
	})
}
