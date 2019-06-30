package lib

import (
	"reflect"
	"testing"
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
