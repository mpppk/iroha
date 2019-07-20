package storage

import (
	"github.com/mpppk/iroha/ktkn"
)

type None struct{}

func (e *None) Get(indices []int) ([][]*ktkn.Word, bool, error) {
	return nil, false, nil
}

func (e *None) Set(indices []int, wordsList [][]*ktkn.Word) error {
	return nil
}
