package lib

import (
	"math/bits"
)

type Iroha struct {
	katakanaBitsMap KatakanaBitsMap
	katakana        *Katakana
}

func NewIroha(words []string) *Iroha {
	km, _ := newKatakanaBitsMap()
	return &Iroha{
		katakanaBitsMap: km,
		katakana:        NewKatakana(words),
	}
}

func (i *Iroha) PrintWordCountMap() {
	i.katakana.PrintWordCountMap()
}

func (i *Iroha) PrintWordByKatakanaMap() {
	i.katakana.PrintWordByKatakanaMap()
}

func (i *Iroha) Search() (rowIndicesList [][]int) {
	katakanaBitsAndWordsList := i.katakana.ListSortedKatakanaBitsAndWords()
	wordsList, _ := i.searchByBits(katakanaBitsAndWordsList, WordBits(0))
	for _, words := range wordsList {
		var rowIndices []int
		for _, word := range words {
			rowIndices = append(rowIndices, int(word.Id))
		}
		rowIndicesList = append(rowIndicesList, rowIndices)
	}
	return
}

func (i *Iroha) searchByBits(katakanaBitsAndWords []*KatakanaBitsAndWords, remainKatakanaBits WordBits) ([][]*Word, bool) {
	if bits.OnesCount64(uint64(remainKatakanaBits)) == int(KatakanaLen) {
		return [][]*Word{{}}, true
	}

	if len(katakanaBitsAndWords) == 0 {
		return nil, false
	}

	katakanaAndWordBits := katakanaBitsAndWords[0]
	var irohaWordLists [][]*Word
	for _, word := range katakanaAndWordBits.Words {
		if remainKatakanaBits.HasDuplicatedKatakana(word.Bits) {
			continue
		}
		newRemainKatakanaBits := remainKatakanaBits.Merge(word.Bits)
		if newIrohaWordIdLists, ok := i.searchByBits(katakanaBitsAndWords[1:], newRemainKatakanaBits); ok {
			for _, newIrohaWordList := range newIrohaWordIdLists {
				newIrohaWordList = append(newIrohaWordList, word)
				irohaWordLists = append(irohaWordLists, newIrohaWordList)
			}
		}
	}

	// どれも入れない場合
	if remainKatakanaBits.has(katakanaAndWordBits.KatakanaBits) {
		if otherIrohaWordBitsLists, ok := i.searchByBits(katakanaBitsAndWords[1:], remainKatakanaBits); ok {
			irohaWordLists = append(irohaWordLists, otherIrohaWordBitsLists...)
		}
	}

	return irohaWordLists, len(irohaWordLists) > 0
}
