package lib

import (
	"fmt"
	"math/bits"
)

type Iroha struct {
	katakanaBitsMap KatakanaBitsMap
	katakana        *Katakana
}

func NewIroha(words []string) *Iroha {
	return &Iroha{
		katakanaBitsMap: newKatakanaBitsMap(),
		katakana:        NewKatakana(words),
	}
}

func (i *Iroha) Search() {
	katakanaAndWordBitsList := i.katakana.ListSortedKatakanaAndWordBits()
	res, _ := i.searchByBits(katakanaAndWordBitsList, WordBits(0))
	fmt.Println(len(res))
}

func (i *Iroha) searchByBits(katakanaAndWordBitsList []*KatakanaAndWordBits, remainKatakanaBits WordBits) ([][]WordBits, bool) {
	if bits.OnesCount64(uint64(remainKatakanaBits)) == int(katakanaLen) {
		return [][]WordBits{{}}, true
	}

	if len(katakanaAndWordBitsList) == 0 {
		return nil, false
	}

	katakanaAndWordBits := katakanaAndWordBitsList[0]
	var irohaWordBitsLists [][]WordBits
	for _, wordBits := range katakanaAndWordBits.WordBitsList {
		if remainKatakanaBits.HasDuplicatedKatakana(wordBits) {
			continue
		}
		newRemainKatakanaBits := remainKatakanaBits.Merge(wordBits)
		if newIrohaWordBitsLists, ok := i.searchByBits(katakanaAndWordBitsList[1:], newRemainKatakanaBits); ok {
			for _, newIrohaWordBitsList := range newIrohaWordBitsLists {
				newIrohaWordBitsList = append(newIrohaWordBitsList, wordBits)
				irohaWordBitsLists = append(irohaWordBitsLists, newIrohaWordBitsList)
			}
		}
	}

	// どれも入れない場合
	if remainKatakanaBits.has(katakanaAndWordBits.KatakanaBits) {
		if otherIrohaWordBitsLists, ok := i.searchByBits(katakanaAndWordBitsList[1:], remainKatakanaBits); ok {
			irohaWordBitsLists = append(irohaWordBitsLists, otherIrohaWordBitsLists...)
		}
	}

	return irohaWordBitsLists, len(irohaWordBitsLists) > 0
}
