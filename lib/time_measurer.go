package lib

import (
	"time"
)

type TimeMeasurer struct {
	start time.Time
}

func NewTimeMeasurer() *TimeMeasurer {
	return &TimeMeasurer{}
}

func NewTimeMeasurerAndStart() *TimeMeasurer {
	t := NewTimeMeasurer()
	t.Start()
	return t
}

func (t *TimeMeasurer) Start() {
	t.start = time.Now()
}

func (t *TimeMeasurer) GetElapsedTimeSec() float64 {
	end := time.Now()
	return end.Sub(t.start).Seconds()
}
