package storage

import (
	"context"
	"sync"

	"github.com/mpppk/iroha/ktkn"
)

type MemoryCache struct {
	m   map[string][][]*ktkn.Word
	mut *sync.RWMutex
}

func newMemoryCache() *MemoryCache {
	return &MemoryCache{
		m:   map[string][][]*ktkn.Word{},
		mut: new(sync.RWMutex),
	}
}

func (mc MemoryCache) Get(key string) ([][]*ktkn.Word, bool) {
	mc.mut.RLock()
	wordsList, ok := mc.m[key]
	mc.mut.RUnlock()
	return wordsList, ok
}

func (mc MemoryCache) Set(key string, wordsList [][]*ktkn.Word) {
	mc.mut.Lock()
	mc.m[key] = wordsList
	mc.mut.Unlock()
}

type Memory struct {
	cache        *MemoryCache
	otherStorage Storage
}

func NewMemory() *Memory {
	return &Memory{
		otherStorage: &None{},
		cache:        newMemoryCache(),
	}
}

func NewMemoryWithOtherStorage(storage Storage) *Memory {
	m := NewMemory()
	m.otherStorage = storage
	return m
}

func (m *Memory) Set(ctx context.Context, indices []int, wordsList [][]*ktkn.Word) error {
	m.cache.Set(toStorageStrKey(indices), wordsList)
	return m.otherStorage.Set(ctx, indices, wordsList)
}

func (m *Memory) Get(ctx context.Context, indices []int) ([][]*ktkn.Word, bool, error) {
	wordsList, ok := m.cache.Get(toStorageStrKey(indices))
	if ok {
		return wordsList, true, nil
	}
	return m.otherStorage.Get(ctx, indices)
}
