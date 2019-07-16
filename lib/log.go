package lib

import (
	"fmt"
	"log"
	"runtime"
	"strings"
)

type Log struct {
	counts           []int
	curs             []int
	measurer         *TimeMeasurer
	depthThreshold   int
	minParallelDepth int
}

func NewLog(katakanaBitsAndWordsList []*KatakanaBitsAndWords, depthThreshold, maxParallelDepth int) *Log {
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
		counts:           counts,
		curs:             make([]int, len(counts), len(counts)),
		measurer:         NewTimeMeasurer(),
		depthThreshold:   depthThreshold,
		minParallelDepth: maxParallelDepth,
	}
}

func (l *Log) updateProgress(depth int) {
	if depth > l.minParallelDepth {
		return
	}
	l.curs[depth]++

	if depth >= len(l.curs) {
		return
	}

	for i := range l.curs {
		if i <= depth {
			continue
		}
		l.curs[i] = 0
	}
}

func (l *Log) getProgress() float64 {
	basePercent := 1.0
	progress := 0.0
	for index, count := range l.counts[:l.minParallelDepth+1] {
		progress += float64(l.curs[index]) / float64(count) * basePercent * 100
		basePercent *= 1.0 / float64(count)
	}
	return progress
}

func (l *Log) PrintProgressLog(depth int, sec float64, msg string) {
	if depth > l.depthThreshold {
		return
	}
	l.updateProgress(depth)
	m := ""
	if msg != "" {
		m = " (" + msg + ")"
	}
	var logs []string
	logs = append(logs, fmt.Sprintf("d:%02v", depth))
	if depth < l.minParallelDepth+1 {
		logs = append(logs, fmt.Sprintf("%04v/%04v", l.curs[depth], l.counts[depth]))
	} else {
		logs = append(logs, "----/----")
	}
	logs = append(logs, fmt.Sprintf("%05.2fsec", sec))
	logs = append(logs, fmt.Sprintf("p:%02.6f", l.getProgress()))
	logs = append(logs, fmt.Sprintf("g:%06d", runtime.NumGoroutine()))
	logs = append(logs, fmt.Sprintf("%02v", l.getPercentSlice()))
	logs = append(logs, m)
	log.Printf(strings.Join(logs, " "))
}

func (l *Log) getPercentSlice() []int {
	percentsLen := l.depthThreshold + 1
	if percentsLen > l.minParallelDepth+1 {
		percentsLen = l.minParallelDepth + 1
	}
	if percentsLen <= 0 {
		return []int{}
	}
	percents := make([]int, percentsLen, percentsLen)
	for i, cur := range l.curs[:percentsLen] {
		if l.counts[i] == 0 {
			percents[i] = 0
			continue
		}
		percents[i] = int((float64(cur) / float64(l.counts[i])) * 100)
	}
	return percents
}
