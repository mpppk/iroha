package storage

import (
	"context"
	"fmt"
	"strconv"

	bolt "github.com/mpppk/bbolt"
	"github.com/mpppk/iroha/ktkn"
)

var progressNotStarted = 0
var progressProcessing = 1
var progressDone = 2

type Storage interface {
	Start(ctx context.Context, indices []int) error
	Set(ctx context.Context, indices []int, wordsList [][]*ktkn.Word) error
	Get(ctx context.Context, indices []int) ([][]*ktkn.Word, bool, error)
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
