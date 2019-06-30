package lib

import (
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
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			katakana := NewKatakana(tt.args.words)
			for katakanaBits, count := range katakana.wordCountMap {
				expectedCount := tt.want.wordCountMap[katakanaBits]
				if expectedCount != count {
					t.Errorf("NewKatakana.wordCountMap[%v] = %v, want %v", katakanaBits, count, expectedCount)
				}
			}
		})
	}
}
