package lib

import "strings"

func newNormalizeKatakanaMap() map[rune]rune {
	m := map[rune]rune{}
	m['ァ'] = 'ア'
	m['ィ'] = 'イ'
	m['ゥ'] = 'ウ'
	m['ェ'] = 'エ'
	m['ォ'] = 'オ'
	m['ッ'] = 'ツ'
	m['ャ'] = 'ヤ'
	m['ュ'] = 'ユ'
	m['ョ'] = 'ヨ'
	m['ガ'] = 'カ'
	m['ギ'] = 'キ'
	m['グ'] = 'ク'
	m['ゲ'] = 'ケ'
	m['ゴ'] = 'コ'
	m['ザ'] = 'サ'
	m['ジ'] = 'シ'
	m['ズ'] = 'ス'
	m['ゼ'] = 'セ'
	m['ゾ'] = 'ソ'
	m['ダ'] = 'タ'
	m['ヂ'] = 'チ'
	m['ヅ'] = 'ツ'
	m['デ'] = 'テ'
	m['ド'] = 'ト'
	m['バ'] = 'ハ'
	m['ビ'] = 'ヒ'
	m['ブ'] = 'フ'
	m['ベ'] = 'ヘ'
	m['ボ'] = 'ホ'
	m['パ'] = 'ハ'
	m['ピ'] = 'ヒ'
	m['プ'] = 'フ'
	m['ペ'] = 'ヘ'
	m['ポ'] = 'ホ'
	m['ヴ'] = 'ウ'
	return m
}

func NormalizeKatakanaWords(words []string) (newWords []string) {
	for _, word := range words {
		newWords = append(newWords, NormalizeKatakanaWord(word))
	}
	return
}

func NormalizeKatakanaWord(word string) string {
	m := newNormalizeKatakanaMap()
	var runes []rune
	for _, w := range word {
		if newW, ok := m[w]; ok {
			runes = append(runes, newW)
			continue
		}
		runes = append(runes, w)
	}
	newWord := string(runes)
	return strings.Replace(newWord, "ー", "", -1)
}
