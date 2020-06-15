package main

import "time"

type Timer struct {
	mark time.Time
}

func StartTimer() {
	return &Timer{mark: time.Now()}
}

func (t *Timer) Lap() float64 {
	now := time.Now()
	lap := (now.Sub(t.mark)).Seconds()
	t.mark = now
	return lap
}
