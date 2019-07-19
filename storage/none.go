package storage

import (
	"github.com/mpppk/iroha/katakana"
)

type None struct{}

func (e *None) Get(indices []int) ([][]*katakana.Word, bool, error) {
	return nil, false, nil
}

func (e *None) Set(indices []int, wordsList [][]*katakana.Word) error {
	return nil
}
