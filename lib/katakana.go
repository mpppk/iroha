package lib

import (
	"fmt"
	"sort"
	"strings"
)

type WordId uint16
type WordBits uint64
type Word struct {
	Id   WordId
	Bits WordBits
}

func (w WordBits) has(katakanaBits KatakanaBits) bool {
	return w&WordBits(katakanaBits) != 0
}

func (w WordBits) HasDuplicatedKatakana(otherWordBits WordBits) bool {
	return w&otherWordBits != 0
}

func (w WordBits) Merge(otherWordBits WordBits) WordBits {
	return w | otherWordBits
}

type KatakanaBits uint64
type KatakanaBitsMap map[rune]KatakanaBits
type RKatakanaBitsMap map[KatakanaBits]rune
type WordByKatakanaMap map[KatakanaBits][]*Word

func (w WordByKatakanaMap) print(wordCountMap WordCountMap, rKatakanaBitsMap RKatakanaBitsMap, wordMap WordMap) {
	for _, katakanaBits := range wordCountMap.toSortedKatakanaBitsList() {
		fmt.Print(string(rKatakanaBitsMap[katakanaBits]), ": ")
		for _, word := range w[katakanaBits] {
			fmt.Print(wordMap[word.Id] + ",")
		}
		fmt.Println("")
	}

}

type WordMap map[WordId]string
type WordCountMap map[KatakanaBits]int

type KatakanaCount struct {
	katakanaBits KatakanaBits
	count        int
}
type KatakanaBitsAndWords struct {
	KatakanaBits KatakanaBits
	Words        []*Word
}

func (w WordCountMap) print(rKatakanaBitsMap RKatakanaBitsMap) {
	katakanaBitsList := w.toSortedKatakanaBitsList()
	for _, katakanaBits := range katakanaBitsList {
		fmt.Print(string(rKatakanaBitsMap[katakanaBits]) + ": ")
		fmt.Println(w[katakanaBits])
	}
}

func (w WordCountMap) toSortedKatakanaBitsList() (katakanaBits []KatakanaBits) {
	sortedKatakanaBitsCounts := w.toSortedList()
	for _, katakanaBitsCount := range sortedKatakanaBitsCounts {
		katakanaBits = append(katakanaBits, katakanaBitsCount.katakanaBits)
	}
	return
}

func (w WordCountMap) toSortedList() []*KatakanaCount {
	katakanaBitsCounts := w.toList()
	sort.Slice(katakanaBitsCounts, func(i, j int) bool {
		a := katakanaBitsCounts[i]
		b := katakanaBitsCounts[j]
		if a.count == b.count {
			return a.katakanaBits < b.katakanaBits
		}
		return a.count < b.count
	})
	return katakanaBitsCounts
}

func (w WordCountMap) toList() []*KatakanaCount {
	var katakanaBitsCounts []*KatakanaCount
	for katakanaBits, count := range w {
		katakanaBitsCounts = append(katakanaBitsCounts, &KatakanaCount{
			katakanaBits: katakanaBits,
			count:        count,
		})
	}
	return katakanaBitsCounts
}

type Katakana struct {
	katakanaBitsMap   KatakanaBitsMap
	RKatakanaBitsMap  RKatakanaBitsMap
	wordByKatakanaMap WordByKatakanaMap
	wordMap           WordMap
	wordCountMap      WordCountMap
}

var KatakanaLen = uint64(45)

func NewKatakana(words []string) *Katakana {
	normalizedWords, orgWords, wordIds := NormalizeAndFilterKatakanaWords(words)
	km, rkm := newKatakanaBitsMap()
	katakana := &Katakana{
		katakanaBitsMap:  km,
		RKatakanaBitsMap: rkm,
		wordMap:          toWordMap(orgWords, wordIds),
	}

	wordBitsList := katakana.loadWords(normalizedWords)
	wordCountMap := countWordBitsFrequency(wordBitsList)
	katakana.wordCountMap = wordCountMap
	katakana.wordByKatakanaMap = katakana.createWordBitsMap(wordBitsList, wordIds)
	return katakana
}

func (k *Katakana) ListSortedKatakanaBitsAndWords() (katakanaAndWordBitsList []*KatakanaBitsAndWords) {
	katakanaBitsList := k.wordCountMap.toSortedKatakanaBitsList()
	for _, katakanaBits := range katakanaBitsList {
		katakanaAndWordBitsList = append(katakanaAndWordBitsList, &KatakanaBitsAndWords{
			KatakanaBits: katakanaBits,
			Words:        k.wordByKatakanaMap[katakanaBits],
		})
	}
	return katakanaAndWordBitsList
}

func (k *Katakana) ToWord(wordId WordId) string {
	return k.wordMap[wordId]
}

func (k *Katakana) loadWords(words []string) (wordBits []WordBits) {
	for _, word := range words {
		wordBits = append(wordBits, k.toWordBits(word))
	}
	return wordBits
}

func (k *Katakana) PrintWordByKatakanaMap() {
	k.wordByKatakanaMap.print(k.wordCountMap, k.RKatakanaBitsMap, k.wordMap)
}

func (k *Katakana) PrintWordCountMap() {
	k.wordCountMap.print(k.RKatakanaBitsMap)
}

func toWordMap(words []string, wordIds []WordId) WordMap {
	wordMap := WordMap{}
	for i, word := range words {
		wordMap[wordIds[i]] = word
	}
	return wordMap
}

func (k *Katakana) toWordBits(word string) WordBits {
	return toWordBits(k.katakanaBitsMap, word)
}

func (k *Katakana) createWordBitsMap(wordBitsList []WordBits, wordIds []WordId) WordByKatakanaMap {
	sortedKatakanaBitsList := k.wordCountMap.toSortedKatakanaBitsList()
	return newWordBitsMap(sortedKatakanaBitsList, wordBitsList, wordIds)
}

func newWordBitsMap(sortedKatakanaBits []KatakanaBits, wordBitsList []WordBits, wordIds []WordId) WordByKatakanaMap {
	var newWordBitsList []WordBits
	copy(newWordBitsList, wordBitsList)

	wordBitsMap := WordByKatakanaMap{}
	for i, wordBits := range wordBitsList {
		for _, katakanaBits := range sortedKatakanaBits {
			if wordBits.has(katakanaBits) {
				wordBitsMap[katakanaBits] = append(wordBitsMap[katakanaBits], &Word{
					Id:   WordId(wordIds[i]),
					Bits: wordBits,
				})
				break
			}
		}
	}
	return wordBitsMap
}

func toWordBits(bitsMap KatakanaBitsMap, word string) WordBits {
	bits := WordBits(0)
	for _, w := range word {
		bits |= WordBits(bitsMap[w])
	}
	return bits
}

func countWordBitsFrequency(wordBitsList []WordBits) WordCountMap {
	wordCountMaps := WordCountMap{}
	for _, wb := range wordBitsList {
		for i := uint64(0); i < KatakanaLen; i++ {
			katakanaBits := KatakanaBits(1 << i)
			if wb.has(katakanaBits) {
				wordCountMaps[katakanaBits]++
			}
		}
	}
	return wordCountMaps
}

func newKatakanaList() []rune {
	return []rune{
		'ア', 'イ', 'ウ', 'エ', 'オ',
		'カ', 'キ', 'ク', 'ケ', 'コ',
		'サ', 'シ', 'ス', 'セ', 'ソ',
		'タ', 'チ', 'ツ', 'テ', 'ト',
		'ナ', 'ニ', 'ヌ', 'ネ', 'ノ',
		'ハ', 'ヒ', 'フ', 'ヘ', 'ホ',
		'マ', 'ミ', 'ム', 'メ', 'モ',
		'ヤ', 'ユ', 'ヨ',
		'ラ', 'リ', 'ル', 'レ', 'ロ',
		'ワ', 'ン',
	}
}

func newKatakanaBitsMap() (KatakanaBitsMap, RKatakanaBitsMap) {
	m := KatakanaBitsMap{}
	rm := RKatakanaBitsMap{}
	for i, katakana := range newKatakanaList() {
		m[katakana] = 1 << uint64(i)
		rm[1<<uint64(i)] = katakana
	}
	return m, rm
}

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

func NormalizeAndFilterKatakanaWords(words []string) (normalizedWords, orgWords []string, wordIds []WordId) {
	wordMap := map[string]struct{}{}
	for wordId, word := range words {
		normalizedWord := NormalizeKatakanaWord(word)
		_, ok := wordMap[word]
		if !HasDuplicatedRune(normalizedWord) && !ok {
			normalizedWords = append(normalizedWords, normalizedWord)
			orgWords = append(orgWords, word)
			wordIds = append(wordIds, WordId(wordId))
			wordMap[word] = struct{}{}
		}
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

func HasDuplicatedRune(word string) bool {
	m := map[rune]struct{}{}
	for _, r := range word {
		if _, ok := m[r]; ok {
			return true
		}
		m[r] = struct{}{}
	}
	return false
}
