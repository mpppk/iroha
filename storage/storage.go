package storage

import (
	"fmt"
	"strconv"

	"github.com/mpppk/iroha/katakana"

	bolt "github.com/mpppk/bbolt"
)

type Storage interface {
	Set(indices []int, wordsList [][]*katakana.Word) error
	Get(indices []int) ([][]*katakana.Word, bool, error)
}

func toStorageStrKey(indices []int) string {
	strKey := ""
	if len(indices) == 0 {
		return "no-index"
	}
	for _, index := range indices {
		strKey += strconv.Itoa(index) + ":"
	}
	return strKey
}

func toStorageKey(indices []int) []byte {
	return []byte(toStorageStrKey(indices))
}

func (b *Bolt) createBucketIfNotExists(bucketName string) error {
	return b.db.Update(func(tx *bolt.Tx) error {
		if _, err := tx.CreateBucketIfNotExists([]byte(bucketName)); err != nil {
			return fmt.Errorf("failed to create %s bucket: %s", bucketName, err)
		}
		return nil
	})
}
