package lib

import (
	"reflect"
	"testing"

	"github.com/k0kubun/pp"
)

func TestNewKatakana(t *testing.T) {
	katakanaBitsMap, _ := newKatakanaBitsMap()
	type args struct {
		words []string
	}
	tests := []struct {
		name string
		args args
		want *Katakana
	}{
		{
			name: "",
			args: args{
				words: []string{"アイウ", "イウエ"},
			},
			want: &Katakana{
				wordCountMap: WordCountMap{
					katakanaBitsMap['ア']: 1,
					katakanaBitsMap['イ']: 2,
					katakanaBitsMap['ウ']: 2,
					katakanaBitsMap['エ']: 1,
				},
				wordByKatakanaMap: WordByKatakanaMap{
					katakanaBitsMap['ア']: []*Word{
						{
							Bits: toWordBits(katakanaBitsMap, "アイウ"),
						},
					},
					katakanaBitsMap['エ']: []*Word{
						{
							Bits: toWordBits(katakanaBitsMap, "イウエ"),
						},
					},
				},
			},
		},
	}

	contains := func(wordByKatakanaMap WordByKatakanaMap, targetKatakanaBits KatakanaBits, targetWordBits WordBits) bool {
		for katakanaBits, words := range wordByKatakanaMap {
			if katakanaBits == targetKatakanaBits {
				for _, word := range words {
					if word.Bits == targetWordBits {
						return true
					}
				}
			}
		}
		return false
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			katakana := NewKatakana(tt.args.words)
			if !reflect.DeepEqual(katakana.wordCountMap, tt.want.wordCountMap) {
				t.Errorf("wordCountMap() = %v, want %v", katakana.wordCountMap, tt.want.wordCountMap)
			}
			for katakanaBits, words := range tt.want.wordByKatakanaMap {
				for _, word := range words {
					if !contains(katakana.wordByKatakanaMap, katakanaBits, word.Bits) {
						t.Errorf("wordByKatakanaMap() = %v, want %v", katakana.wordByKatakanaMap, tt.want.wordByKatakanaMap)
					}
				}
			}
		})
	}
}

func TestKatakana_ToSortedKatakanaAndWordBits(t *testing.T) {
	katakanaBitsMap, _ := newKatakanaBitsMap()
	type fields struct {
		katakanaBitsMap KatakanaBitsMap
		wordBitsMap     WordByKatakanaMap
		wordCountMap    WordCountMap
	}
	tests := []struct {
		name                        string
		fields                      fields
		wantKatakanaAndWordBitsList []*KatakanaBitsAndWords
	}{
		{
			name: "",
			fields: fields{
				wordCountMap: WordCountMap{
					katakanaBitsMap['ア']: 1,
					katakanaBitsMap['イ']: 2,
					katakanaBitsMap['ウ']: 2,
					katakanaBitsMap['エ']: 1,
				},
				wordBitsMap: WordByKatakanaMap{
					katakanaBitsMap['ア']: []*Word{
						{
							Bits: toWordBits(katakanaBitsMap, "アイウ"),
						},
					},
					katakanaBitsMap['エ']: []*Word{
						{
							Bits: toWordBits(katakanaBitsMap, "イウエ"),
						},
					},
				},
				katakanaBitsMap: katakanaBitsMap,
			},
			wantKatakanaAndWordBitsList: []*KatakanaBitsAndWords{
				{
					KatakanaBits: katakanaBitsMap['ア'],
					Words: []*Word{
						{
							Bits: toWordBits(katakanaBitsMap, "アイウ"),
						},
					},
				},
				{
					KatakanaBits: katakanaBitsMap['エ'],
					Words: []*Word{
						{
							Bits: toWordBits(katakanaBitsMap, "イウエ"),
						},
					},
				},
				{
					KatakanaBits: katakanaBitsMap['イ'],
					Words:        []*Word{},
				},
				{
					KatakanaBits: katakanaBitsMap['ウ'],
					Words:        []*Word{},
				},
			},
		},
	}

	contains := func(list []*KatakanaBitsAndWords, v *KatakanaBitsAndWords) bool {
		for _, nv := range list {
			if nv.KatakanaBits == v.KatakanaBits {
				// FIXME
				if len(v.Words) == 0 && len(nv.Words) == 0 {
					return true
				}
				return reflect.DeepEqual(*nv, *v)
			}
		}
		return false
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k := &Katakana{
				katakanaBitsMap:   tt.fields.katakanaBitsMap,
				wordByKatakanaMap: tt.fields.wordBitsMap,
				wordCountMap:      tt.fields.wordCountMap,
			}

			gotKatakanaAndWordBitsList := k.ListSortedKatakanaBitsAndWords()
			for _, want := range tt.wantKatakanaAndWordBitsList {
				if !contains(gotKatakanaAndWordBitsList, want) {
					pp.Println(gotKatakanaAndWordBitsList)
					t.Errorf("KatakanaAndWordBitsList = %v, should contains %v", gotKatakanaAndWordBitsList, want)
				}
			}
		})
	}
}
