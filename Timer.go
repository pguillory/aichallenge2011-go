package main

import "time"
import "fmt"
import "bytes"

type Timer struct {
    start int64
    currentName string
    names []string
    times []int64
}

func now() int64 {
    return time.Nanoseconds() / 1000000
}

func NewTimer() *Timer {
    this := new(Timer)
    this.names = make([]string, 0)
    this.times = make([]int64, 0)
    return this
}

func (this *Timer) Start(name string) {
    this.currentName = name
    this.start = now()
}

func (this *Timer) Stop() {
    end := now()
    this.names = append(this.names, this.currentName)
    this.times = append(this.times, end - this.start)
    //this.times[this.currentName] = end - this.start
}

func (this *Timer) ForEach(f func(string, int64)) {
    for i, name := range this.names {
        f(name, this.times[i])
    }
}

func (this *Timer) String() string {
    b := new(bytes.Buffer)
    this.ForEach(func(name string, runtime int64) {
        b.WriteString(fmt.Sprintf("%v:%v ", name, runtime))
    })
    return b.String()
}
