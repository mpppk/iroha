package lib

import (
	"reflect"
	"testing"

	"github.com/k0kubun/pp"
)

func TestNewKatakana(t *testing.T) {
	katakanaBitsMap := newKatakanaBitMap()
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
				wordBitsMap: WordBitsMap{
					katakanaBitsMap['ア']: []WordBits{
						toWordBits(katakanaBitsMap, "アイウ"),
					},
					katakanaBitsMap['エ']: []WordBits{
						toWordBits(katakanaBitsMap, "イウエ"),
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			katakana := NewKatakana(tt.args.words)
			if !reflect.DeepEqual(katakana.wordCountMap, tt.want.wordCountMap) {
				t.Errorf("wordCountMap() = %v, want %v", katakana.wordCountMap, tt.want.wordCountMap)
			}
			if !reflect.DeepEqual(katakana.wordBitsMap, tt.want.wordBitsMap) {
				t.Errorf("wordBitsMap() = %v, want %v", katakana.wordBitsMap, tt.want.wordBitsMap)
			}
		})
	}
}

func TestKatakana_ToSortedKatakanaAndWordBits(t *testing.T) {
	katakanaBitsMap := newKatakanaBitMap()
	type fields struct {
		katakanaBitsMap KatakanaBitsMap
		wordBitsMap     WordBitsMap
		wordCountMap    WordCountMap
	}
	tests := []struct {
		name                        string
		fields                      fields
		wantKatakanaAndWordBitsList []*KatakanaAndWordBits
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
				wordBitsMap: WordBitsMap{
					katakanaBitsMap['ア']: []WordBits{
						toWordBits(katakanaBitsMap, "アイウ"),
					},
					katakanaBitsMap['エ']: []WordBits{
						toWordBits(katakanaBitsMap, "イウエ"),
					},
				},
				katakanaBitsMap: katakanaBitsMap,
			},
			wantKatakanaAndWordBitsList: []*KatakanaAndWordBits{
				{
					KatakanaBits: katakanaBitsMap['ア'],
					WordBitsList: []WordBits{
						toWordBits(katakanaBitsMap, "アイウ"),
					},
				},
				{
					KatakanaBits: katakanaBitsMap['エ'],
					WordBitsList: []WordBits{
						toWordBits(katakanaBitsMap, "イウエ"),
					},
				},
				{
					KatakanaBits: katakanaBitsMap['イ'],
					WordBitsList: []WordBits{},
				},
				{
					KatakanaBits: katakanaBitsMap['ウ'],
					WordBitsList: []WordBits{},
				},
			},
		},
	}

	contains := func(list []*KatakanaAndWordBits, v *KatakanaAndWordBits) bool {
		for _, nv := range list {
			if nv.KatakanaBits == v.KatakanaBits {
				// FIXME
				if len(v.WordBitsList) == 0 && len(nv.WordBitsList) == 0 {
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
				katakanaBitsMap: tt.fields.katakanaBitsMap,
				wordBitsMap:     tt.fields.wordBitsMap,
				wordCountMap:    tt.fields.wordCountMap,
			}

			gotKatakanaAndWordBitsList := k.ListSortedKatakanaAndWordBits()
			for _, want := range tt.wantKatakanaAndWordBitsList {
				if !contains(gotKatakanaAndWordBitsList, want) {
					pp.Println(gotKatakanaAndWordBitsList)
					t.Errorf("KatakanaAndWordBitsList = %v, should contains %v", gotKatakanaAndWordBitsList, want)
				}
			}
		})
	}
}
