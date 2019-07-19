package storage

import (
	"github.com/mpppk/iroha/katakana"
)

type Memory struct {
	cache        map[string][][]*katakana.Word
	otherStorage Storage
}

func NewMemory(storage Storage) *Memory {
	return &Memory{
		otherStorage: &None{},
	}

}

func NewMemoryWithOtherStorage(storage Storage) *Memory {
	return &Memory{
		otherStorage: storage,
	}
}

func (m *Memory) Set(indices []int, wordsList [][]*katakana.Word) error {
	m.cache[toStorageStrKey(indices)] = wordsList
	return m.otherStorage.Set(indices, wordsList)
}

func (m *Memory) Get(indices []int) ([][]*katakana.Word, bool, error) {
	wordsList, ok := m.cache[toStorageStrKey(indices)]
	if ok {
		return wordsList, true, nil
	}
	return m.otherStorage.Get(indices)
}
