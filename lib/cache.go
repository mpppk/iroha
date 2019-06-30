package lib

import "sort"

type IrohaCache map[int]map[string][][]string

func (i IrohaCache) Set(words []string, iroha string, value [][]string) {
	sortedIroha := sortIroha(iroha)
	index := len(words)
	_, ok := i[index]
	if !ok {
		i[index] = map[string][][]string{}
	}
	i[index][sortedIroha] = value
}

func (i IrohaCache) Get(words []string, iroha string) ([][]string, bool) {
	m, ok := i[len(words)]
	if !ok {
		return nil, false
	}
	sortedIroha := sortIroha(iroha)
	value, ok := m[sortedIroha]
	return value, ok
}

func sortIroha(iroha string) string {
	irohaRunes := []rune(iroha)
	sort.Slice(irohaRunes, func(i, j int) bool {
		return irohaRunes[i] > irohaRunes[j]
	})
	return string(irohaRunes)
}
