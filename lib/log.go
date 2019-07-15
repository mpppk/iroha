package lib

import (
	"log"
	"runtime"
)

type Log struct {
	counts   []int
	curs     []int
	measurer *TimeMeasurer
}

func NewLog(katakanaBitsAndWordsList []*KatakanaBitsAndWords) *Log {
	counts := make([]int, len(katakanaBitsAndWordsList), len(katakanaBitsAndWordsList))
	for index, katakanaBitsAndWords := range katakanaBitsAndWordsList {
		counts[index] = len(katakanaBitsAndWords.Words)
	}
	countsSum := 0
	for _, count := range counts {
		countsSum += count
	}
	log.Printf("counts: %.2v %d\n", counts, countsSum)

	return &Log{
		counts:   counts,
		curs:     make([]int, len(counts), len(counts)),
		measurer: NewTimeMeasurer(),
	}
}

func (l *Log) updateProgress(depth, cur int) {
	l.curs[depth] = cur
	for i := range l.curs {
		if i <= depth {
			continue
		}
		l.curs[i] = 0
	}
}

func (l *Log) getProgress(depth int) float64 {
	max := float64(l.counts[0])
	cur := float64(l.curs[0])
	for index, count := range l.counts[1 : depth+1] {
		max *= float64(count)
		cur *= float64(l.curs[index])
	}
	return (cur / max) * 100
}

func (l *Log) PrintProgressLog(depth, current int, sec float64, msg string) {
	l.updateProgress(depth, current)
	m := ""
	if msg != "" {
		m = " (" + msg + ")"
	}
	percents := make([]int, len(l.curs), len(l.curs))
	for i, cur := range l.curs {
		if l.counts[i] == 0 {
			percents[i] = 0
			continue
		}
		percents[i] = int((float64(cur) / float64(l.counts[i])) * 100)
	}
	log.Printf("d:%02v %04v/%04v %05.2fsec %03.6f%s gr=%06d %02d%s",
		depth,
		current,
		l.counts[depth],
		sec,
		l.getProgress(depth),
		"%",
		runtime.NumGoroutine(),
		percents,
		m,
	)
}
