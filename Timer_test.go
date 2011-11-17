package main

import "testing"
import "time"

func TestTimer(t *testing.T) {
    timer := NewTimer()
    timer.Start("test")
    time.Sleep(10 * 1000000)
    timer.Stop()
    if timer.times[0] < 10 || timer.times[0] > 11 {
        t.Error(timer)
    }
}
