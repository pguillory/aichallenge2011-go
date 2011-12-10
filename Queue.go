package main

import "unsafe"

const QUEUEBLOCK_CAPACITY = 100000

type Queue struct {
    members [QUEUEBLOCK_CAPACITY]unsafe.Pointer
    start, end int
}

func (this *Queue) Push(pointer unsafe.Pointer) {
    this.members[this.end] = pointer
    this.end += 1
}

func (this *Queue) Pop() unsafe.Pointer {
    result := this.members[this.start]
    this.start += 1
    return result
}

func (this *Queue) Empty() bool {
    return (this.start == this.end)
}
