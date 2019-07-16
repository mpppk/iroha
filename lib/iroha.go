package lib

import (
	"fmt"
	"math/bits"

	"golang.org/x/sync/errgroup"
)

var logDepthThreshold = 1
var minParallelDepth = 1
var maxParallelDepth = 2

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
	i.log = NewLog(katakanaBitsAndWordsList, logDepthThreshold, minParallelDepth)
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

func (i *Iroha) f(word *Word, katakanaBitsAndWords []*KatakanaBitsAndWords, remainKatakanaBits WordBits) ([][]*Word, bool, error) {
	var results [][]*Word
	if remainKatakanaBits.HasDuplicatedKatakana(word.Bits) {
		return nil, false, nil
	}
	newRemainKatakanaBits := remainKatakanaBits.Merge(word.Bits)
	newIrohaWordIdLists, ok, err := i.searchByBits(katakanaBitsAndWords[1:], newRemainKatakanaBits)
	if err != nil {
		return nil, false, err
	}
	if ok {
		for _, newIrohaWordList := range newIrohaWordIdLists {
			newIrohaWordList = append(newIrohaWordList, word)
			results = append(results, newIrohaWordList)
		}
	}
	return results, true, nil
}

func (i *Iroha) gf(word *Word, katakanaBitsAndWords []*KatakanaBitsAndWords, remainKatakanaBits WordBits, wordListChan chan<- []*Word) error {
	wordLists, ok, err := i.f(word, katakanaBitsAndWords, remainKatakanaBits)
	if err != nil {
		return err
	}
	if ok {
		for _, wordList := range wordLists {
			wordListChan <- wordList
		}
	}
	return nil
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
	if len(katakanaAndWordBits.Words) == 0 {
		return nil, false, nil
	}

	depth := int(KatakanaLen) - len(katakanaBitsAndWords)
	var irohaWordLists [][]*Word

	goroutineMode := depth >= minParallelDepth && depth <= maxParallelDepth

	if goroutineMode {
		eg := errgroup.Group{}
		wordListChan := make(chan []*Word, 100)
		for _, word := range katakanaAndWordBits.Words {
			gWord := word
			eg.Go(func() error {
				if err := i.gf(gWord, katakanaBitsAndWords, remainKatakanaBits, wordListChan); err != nil {
					return err
				}
				i.log.PrintProgressLog(depth, "")
				return nil
			})
		}
		errChan := make(chan error)
		go func() {
			if err := eg.Wait(); err != nil {
				errChan <- err
				return
			}
			close(wordListChan)
		}()

	L:
		for {
			select {
			case wordList, ok := <-wordListChan:
				if !ok {
					break L
				}
				irohaWordLists = append(irohaWordLists, wordList)
			case err, ok := <-errChan:
				if !ok {
					return nil, false, fmt.Errorf("unexpected error channel closing")
				}
				return nil, false, err
			}
		}
	} else {
		for _, word := range katakanaAndWordBits.Words {
			wordList, ok, err := i.f(word, katakanaBitsAndWords, remainKatakanaBits)
			if err != nil {
				return nil, false, err
			}
			if ok {
				irohaWordLists = append(irohaWordLists, wordList...)
			}
			i.log.PrintProgressLog(depth, "")
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
