package lib

import (
	"context"
	"fmt"
	"math/bits"

	"github.com/mpppk/iroha/storage"

	"github.com/mpppk/iroha/ktkn"

	"golang.org/x/sync/errgroup"
)

type Iroha struct {
	katakanaBitsMap ktkn.KatakanaBitsMap
	katakana        *ktkn.Katakana
	log             *Log
	depths          *DepthOptions
	storage         storage.Storage
}

type DepthOptions struct {
	MaxLog      int
	MinParallel int
	MaxParallel int
	MaxStorage  int
}

func NewIroha(words []string, storage storage.Storage, options *DepthOptions) *Iroha {
	km, _ := ktkn.NewKatakanaBitsMap()
	return &Iroha{
		katakanaBitsMap: km,
		katakana:        ktkn.NewKatakana(words),
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
	wordsList, _, err := i.searchByBits([]int{}, katakanaBitsAndWordsList, ktkn.WordBits(0))
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

func (i *Iroha) f(word *ktkn.Word, usedIndices []int, katakanaBitsAndWords []*ktkn.KatakanaBitsAndWords, remainKatakanaBits ktkn.WordBits) ([][]*ktkn.Word, bool, error) {
	var results [][]*ktkn.Word
	if remainKatakanaBits.HasDuplicatedKatakana(word.Bits) {
		return nil, false, nil
	}
	newRemainKatakanaBits := remainKatakanaBits.Merge(word.Bits)
	newIrohaWordIdLists, ok, err := i.searchByBits(usedIndices, katakanaBitsAndWords[1:], newRemainKatakanaBits)
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

func (i *Iroha) gf(word *ktkn.Word, usedIndices []int, katakanaBitsAndWords []*ktkn.KatakanaBitsAndWords, remainKatakanaBits ktkn.WordBits, wordListChan chan<- []*ktkn.Word) error {
	// FIXME: check cache
	wordLists, ok, err := i.f(word, usedIndices, katakanaBitsAndWords, remainKatakanaBits)
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

func (i *Iroha) searchByBits(usedIndices []int, katakanaBitsAndWords []*ktkn.KatakanaBitsAndWords, remainKatakanaBits ktkn.WordBits) ([][]*ktkn.Word, bool, error) {
	ctx := context.Background()
	remainKatakanaNum := bits.OnesCount64(uint64(remainKatakanaBits))
	if remainKatakanaNum == int(ktkn.KatakanaLen) {
		return [][]*ktkn.Word{{}}, true, nil
	}

	if len(katakanaBitsAndWords) == 0 {
		return nil, false, nil
	}

	katakanaAndWordBits := katakanaBitsAndWords[0]
	if len(katakanaAndWordBits.Words) == 0 {
		return nil, false, nil
	}

	depth := int(ktkn.KatakanaLen) - len(katakanaBitsAndWords) - 1

	if depth <= i.depths.MaxStorage {
		if results, ok, err := i.storage.Get(ctx, usedIndices); err != nil {
			return nil, false, err
		} else if ok {
			i.log.PrintProgressLog(depth, "cache used")
			return results, true, nil
		}
	}

	if depth == i.depths.MaxStorage {
		if err := i.storage.Start(ctx, usedIndices); err != nil {
			return nil, false, err
		}
	}

	var irohaWordLists [][]*ktkn.Word
	goroutineMode := depth >= i.depths.MinParallel && depth <= i.depths.MaxParallel
	if goroutineMode {
		eg := errgroup.Group{}
		wordListChan := make(chan []*ktkn.Word, 1000000)
		for index, word := range katakanaAndWordBits.Words {
			gWord := word
			newIndices := generateNewUsedIndices(usedIndices, index)
			eg.Go(func() error {
				if err := i.gf(gWord, newIndices, katakanaBitsAndWords, remainKatakanaBits, wordListChan); err != nil {
					return err
				}
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
		for index, word := range katakanaAndWordBits.Words {
			newIndices := generateNewUsedIndices(usedIndices, index)
			wordList, ok, err := i.f(word, newIndices, katakanaBitsAndWords, remainKatakanaBits)
			if err != nil {
				return nil, false, err
			}
			if ok {
				irohaWordLists = append(irohaWordLists, wordList...)
			}
		}
	}

	// どれも入れない場合
	if remainKatakanaBits.Has(katakanaAndWordBits.KatakanaBits) {
		otherIrohaWordBitsLists, ok, err := i.searchByBits(append(usedIndices, -1), katakanaBitsAndWords[1:], remainKatakanaBits)
		if err != nil {
			return nil, false, err
		}
		if ok {
			irohaWordLists = append(irohaWordLists, otherIrohaWordBitsLists...)
		}
	}

	msg := ""
	if depth <= i.depths.MaxStorage {
		if err := i.storage.Set(ctx, usedIndices, irohaWordLists); err != nil {
			return nil, false, err
		}
		msg += "cache saved"
	}
	i.log.PrintProgressLog(depth, msg)
	return irohaWordLists, len(irohaWordLists) > 0, nil
}

func generateNewUsedIndices(usedIndices []int, newIndex int) []int {
	newIndices := make([]int, len(usedIndices), len(usedIndices)+1)
	copy(newIndices, usedIndices)
	newIndices = append(newIndices, newIndex)
	return newIndices
}
