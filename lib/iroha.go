package lib

import (
	"fmt"
	"strings"
)

func NewIroha() string {
	return "アイウエオカキクケコサシスセソタチツテトナニヌネノハヒフヘホマミムメモヤユヨラリルレロワオン"
}

//func IsIroha()

func CanUseIroha(words []string, remainIroha string, cache IrohaCache) [][]string {
	//fmt.Println("words", words, remainIroha)
	if len(words) == 0 {
		return [][]string{}
	}

	if cachedResult, ok := cache.Get(words, remainIroha); ok {
		fmt.Println("cache used")
		return cachedResult
	}

	// 使わない場合
	wordsList := CanUseIroha(words[1:], remainIroha, cache)
	cache.Set(words[1:], remainIroha, wordsList)

	// 使う場合
	// 現在着目している単語がまだ使っていない文字で表現できるかチェック
	// できない場合はerrが返ってくる
	newRemainIroha, err := removeIrohaByWord(words[0], remainIroha)
	if err == nil && newRemainIroha == "" {
		fmt.Println("found", words[0], remainIroha, newRemainIroha)
		return [][]string{{words[0]}}
	}

	if err != nil {
		return wordsList
	}

	wordsList2 := CanUseIroha(words[1:], newRemainIroha, cache)
	cache.Set(words[1:], newRemainIroha, wordsList)

	var newWordsList [][]string
	for _, words2 := range wordsList2 {
		newWords := append([]string{words[0]}, words2...)
		newWordsList = append(newWordsList, newWords)
	}
	ret := append(wordsList, newWordsList...)
	return ret
}

func removeIrohaByWord(word string, remainIroha string) (newIroha string, err error) {
	newIroha = remainIroha
	for _, r := range word {
		if strings.ContainsRune(newIroha, r) {
			newIroha = strings.Replace(newIroha, string(r), "", 1)
			continue
		}
		return "", fmt.Errorf("rune %q does not found in word %q", r, word)

		//rIndex := strings.IndexRune(newIroha, r)
		//if rIndex < 0 {
		//	return "", fmt.Errorf("rune %q does not found in word %q", r, word)
		//}
		//wordSlice := []rune(word)
		//if rIndex == len(wordSlice)-1 {
		//	newIroha = string(wordSlice[:rIndex])
		//	continue
		//}
		//fmt.Println(word, wordSlice, rIndex)
		//newIroha = string(append(wordSlice[:rIndex], wordSlice[rIndex+1:]...))
	}
	return newIroha, nil
}

func uniq(word string) string {
	m := map[rune]bool{}
	for _, r := range word {
		m[r] = true
	}

	var uniqRunes []rune
	for key := range m {
		uniqRunes = append([]rune(uniqRunes), key)
	}
	return string(uniqRunes)
}
