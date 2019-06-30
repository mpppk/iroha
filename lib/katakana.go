package lib

import "sort"

type WordBits uint64
type KatakanaBits uint64
type KatakanaBitsMap map[rune]KatakanaBits
type WordBitsMap map[KatakanaBits][]WordBits
type WordCountMap map[KatakanaBits]int
type KatakanaCount struct {
	katakanaBits KatakanaBits
	count        int
}

func (w WordCountMap) getSortedKatakanaBitsCounts() []*KatakanaCount {
	katakanaBitsCounts := w.toList()
	sort.Slice(katakanaBitsCounts, func(i, j int) bool {
		return katakanaBitsCounts[i].count > katakanaBitsCounts[j].count
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
	katakanaBitsMap KatakanaBitsMap
	wordBitsMap     WordBitsMap
	wordCountMap    WordCountMap
}

var katakanaLen = uint64(45)

func NewKatakana(words []string) *Katakana {
	katakana := &Katakana{
		katakanaBitsMap: newKatakanaBitMap(),
	}

	wordBitsList := katakana.loadWords(words)
	wordCountMap := countWordBitsFrequency(wordBitsList)
	katakana.wordCountMap = wordCountMap
	return katakana
}

func (k *Katakana) loadWords(words []string) (wordBits []WordBits) {
	for _, word := range words {
		wordBits = append(wordBits, k.toWordBits(word))
	}
	return wordBits
}

func (k *Katakana) toWordBits(word string) WordBits {
	return toWordBits(k.katakanaBitsMap, word)
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
		for i := uint64(0); i < katakanaLen; i++ {
			katakanaBits := KatakanaBits(1 << i)
			if wb&WordBits(katakanaBits) == 1 {
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

func newKatakanaBitMap() KatakanaBitsMap {
	m := KatakanaBitsMap{}
	for i, katakana := range newKatakanaList() {
		m[katakana] = 1 << uint64(i)
	}
	return m
}
