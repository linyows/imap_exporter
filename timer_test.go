package main

import (
	"math"
	"testing"
	"time"
)

func TestTimer(t *testing.T) {
	timer := StartTimer()
	time.Sleep(time.Second * 3 / 10)
	duration1 := int(math.Floor(timer.Lap() * 10))
	if 3 != duration1 {
		t.Errorf("Timer 1st rap failure")
	}
	time.Sleep(time.Second * 5 / 10)
	duration2 := int(math.Floor(timer.Lap() * 10))
	if 5 != duration2 {
		t.Errorf("Timer 2nd rap failure")
	}
}
