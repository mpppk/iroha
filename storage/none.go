package storage

import (
	"context"

	"github.com/mpppk/iroha/ktkn"
)

type None struct{}

func (e *None) Get(ctx context.Context, indices []int) ([][]*ktkn.Word, bool, error) {
	return nil, false, nil
}

func (e *None) Set(ctx context.Context, indices []int, wordsList [][]*ktkn.Word) error {
	return nil
}
