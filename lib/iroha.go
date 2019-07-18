package lib

import (
	"fmt"
	"math/bits"

	"golang.org/x/sync/errgroup"
)

type Iroha struct {
	katakanaBitsMap KatakanaBitsMap
	katakana        *Katakana
	log             *Log
	depths          *DepthOptions
	storage         Storage
}

type DepthOptions struct {
	MaxLog      int
	MinParallel int
	MaxParallel int
	MaxStorage  int
}

func NewIroha(words []string, storage Storage, options *DepthOptions) *Iroha {
	km, _ := newKatakanaBitsMap()
	return &Iroha{
		katakanaBitsMap: km,
		katakana:        NewKatakana(words),
		depths:          options,
		storage:         storage,
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
	i.log = NewLog(katakanaBitsAndWordsList, i.depths.MaxLog, i.depths.MinParallel)
	wordsList, _, _, err := i.searchByBits([]int{}, katakanaBitsAndWordsList, WordBits(0))
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

func (i *Iroha) f(word *Word, usedIndices []int, katakanaBitsAndWords []*KatakanaBitsAndWords, remainKatakanaBits WordBits) ([][]*Word, bool, bool, error) {
	var results [][]*Word
	if remainKatakanaBits.HasDuplicatedKatakana(word.Bits) {
		return nil, false, false, nil
	}
	newRemainKatakanaBits := remainKatakanaBits.Merge(word.Bits)
	newIrohaWordIdLists, ok, cacheUsed, err := i.searchByBits(usedIndices, katakanaBitsAndWords[1:], newRemainKatakanaBits)
	if err != nil {
		return nil, false, false, err
	}
	if ok {
		for _, newIrohaWordList := range newIrohaWordIdLists {
			newIrohaWordList = append(newIrohaWordList, word)
			results = append(results, newIrohaWordList)
		}
	}
	return results, true, cacheUsed, nil
}

func (i *Iroha) gf(word *Word, usedIndices []int, katakanaBitsAndWords []*KatakanaBitsAndWords, remainKatakanaBits WordBits, wordListChan chan<- []*Word) error {
	// FIXME: check cache
	wordLists, ok, _, err := i.f(word, usedIndices, katakanaBitsAndWords, remainKatakanaBits)
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

func (i *Iroha) searchByBits(usedIndices []int, katakanaBitsAndWords []*KatakanaBitsAndWords, remainKatakanaBits WordBits) ([][]*Word, bool, bool, error) {
	remainKatakanaNum := bits.OnesCount64(uint64(remainKatakanaBits))
	if remainKatakanaNum == int(KatakanaLen) {
		return [][]*Word{{}}, true, false, nil
	}

	if len(katakanaBitsAndWords) == 0 {
		return nil, false, false, nil
	}

	katakanaAndWordBits := katakanaBitsAndWords[0]
	if len(katakanaAndWordBits.Words) == 0 {
		return nil, false, false, nil
	}

	depth := int(KatakanaLen) - len(katakanaBitsAndWords)

	if depth <= i.depths.MaxStorage {
		if results, ok, err := i.storage.Get(usedIndices); err != nil {
			return nil, false, false, err
		} else if ok {
			return results, true, true, nil
		}
	}

	var irohaWordLists [][]*Word
	goroutineMode := depth >= i.depths.MinParallel && depth <= i.depths.MaxParallel
	if goroutineMode {
		eg := errgroup.Group{}
		wordListChan := make(chan []*Word, 100)
		for index, word := range katakanaAndWordBits.Words {
			gWord := word
			newIndices := generateNewUsedIndices(usedIndices, index)
			eg.Go(func() error {
				if err := i.gf(gWord, newIndices, katakanaBitsAndWords, remainKatakanaBits, wordListChan); err != nil {
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
					return nil, false, false, fmt.Errorf("unexpected error channel closing")
				}
				return nil, false, false, err
			}
		}
	} else {
		for index, word := range katakanaAndWordBits.Words {
			newIndices := generateNewUsedIndices(usedIndices, index)
			wordList, ok, cacheUsed, err := i.f(word, newIndices, katakanaBitsAndWords, remainKatakanaBits)
			if err != nil {
				return nil, false, false, err
			}
			if ok {
				irohaWordLists = append(irohaWordLists, wordList...)
			}
			msg := ""
			if cacheUsed {
				msg = "cache used"
			} else if depth <= i.depths.MaxStorage {
				msg = "cache saved"
			}
			i.log.PrintProgressLog(depth, msg)
		}
	}

	// どれも入れない場合
	if remainKatakanaBits.has(katakanaAndWordBits.KatakanaBits) {
		otherIrohaWordBitsLists, ok, cacheUsed, err := i.searchByBits(append(usedIndices, -1), katakanaBitsAndWords[1:], remainKatakanaBits)
		if err != nil {
			return nil, false, false, err
		}
		if ok {
			irohaWordLists = append(irohaWordLists, otherIrohaWordBitsLists...)
		}
		msg := "no-add"
		if cacheUsed {
			msg += " / cache used"
		} else if depth <= i.depths.MaxStorage {
			msg += " / cache saved"
		}

		i.log.PrintProgressLog(depth, msg)
	}

	if depth <= i.depths.MaxStorage {
		if err := i.storage.Set(usedIndices, irohaWordLists); err != nil {
			return nil, false, false, err
		}
	}
	return irohaWordLists, len(irohaWordLists) > 0, false, nil
}

func generateNewUsedIndices(usedIndices []int, newIndex int) []int {
	newIndices := make([]int, len(usedIndices), len(usedIndices)+1)
	copy(newIndices, usedIndices)
	newIndices = append(newIndices, newIndex)
	return newIndices
}
