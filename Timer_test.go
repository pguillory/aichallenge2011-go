package main

import "testing"
import "time"

func TestTimer(t *testing.T) {
    timer := NewTimer()
    timer.Start("test")
    time.Sleep(10 * 1000000)
    timer.Stop()
    if timer.times["test"] < 10 || timer.times["test"] > 11 {
        t.Error(timer)
    }

    timer.Start("test")
    time.Sleep(10 * 1000000)
    timer.Stop()
    if timer.times["test"] < 20 || timer.times["test"] > 21 {
        t.Error(timer)
    }
}
