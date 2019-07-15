package lib

import (
	"math/bits"
)

type Iroha struct {
	katakanaBitsMap KatakanaBitsMap
	katakana        *Katakana
	log             *Log
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

func (i *Iroha) Search() (rowIndicesList [][]int, err error) {
	katakanaBitsAndWordsList := i.katakana.ListSortedKatakanaBitsAndWords()
	i.log = NewLog(katakanaBitsAndWordsList)
	wordsList, _, err := i.searchByBits(katakanaBitsAndWordsList, WordBits(0))
	if err != nil {
		return nil, err
	}
	for _, words := range wordsList {
		var rowIndices []int
		for _, word := range words {
			rowIndices = append(rowIndices, int(word.Id))
		}
		rowIndicesList = append(rowIndicesList, rowIndices)
	}
	return
}

func (i *Iroha) searchByBits(katakanaBitsAndWords []*KatakanaBitsAndWords, remainKatakanaBits WordBits) ([][]*Word, bool, error) {
	remainKatakanaNum := bits.OnesCount64(uint64(remainKatakanaBits))
	if remainKatakanaNum == int(KatakanaLen) {
		return [][]*Word{{}}, true, nil
	}

	if len(katakanaBitsAndWords) == 0 {
		return nil, false, nil
	}

	katakanaAndWordBits := katakanaBitsAndWords[0]
	depth := int(KatakanaLen) - len(katakanaBitsAndWords)
	var irohaWordLists [][]*Word
	for cur, word := range katakanaAndWordBits.Words {
		measurer := NewTimeMeasurerAndStart()
		if remainKatakanaBits.HasDuplicatedKatakana(word.Bits) {
			continue
		}
		newRemainKatakanaBits := remainKatakanaBits.Merge(word.Bits)
		newIrohaWordIdLists, ok, err := i.searchByBits(katakanaBitsAndWords[1:], newRemainKatakanaBits)
		if err != nil {
			return nil, false, err
		}
		if ok {
			for _, newIrohaWordList := range newIrohaWordIdLists {
				newIrohaWordList = append(newIrohaWordList, word)
				irohaWordLists = append(irohaWordLists, newIrohaWordList)
			}
		}
		if t := measurer.GetElapsedTimeSec(); t > 5 {
			i.log.PrintProgressLog(depth, cur, t)
		}
	}

	// どれも入れない場合
	if remainKatakanaBits.has(katakanaAndWordBits.KatakanaBits) {
		otherIrohaWordBitsLists, ok, err := i.searchByBits(katakanaBitsAndWords[1:], remainKatakanaBits)
		if err != nil {
			return nil, false, err
		}
		if ok {
			irohaWordLists = append(irohaWordLists, otherIrohaWordBitsLists...)
		}
	}

	return irohaWordLists, len(irohaWordLists) > 0, nil
}
