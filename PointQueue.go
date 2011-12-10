package main

import "fmt"

const POINTQUEUE_CAPACITY = 40000

type PointQueue struct {
    points [POINTQUEUE_CAPACITY]Point
    start, end int
}

func (this *PointQueue) Push(p Point) {
    if this.end >= POINTQUEUE_CAPACITY {
        panic(fmt.Sprintf("PointQueue ran out of capacity, start=%v, end=%v, capi", this.start, this.end))
    }
    this.points[this.end] = p
    this.end += 1
}

func (this *PointQueue) Pop() (Point) {
    p := this.points[this.start]
    this.start += 1
    return p
}

func (this *PointQueue) Empty() bool {
    return (this.start == this.end)
}

func (this *PointQueue) Size() int {
    return (this.end - this.start)
}

func (this *PointQueue) ForEach(f func(Point)) {
    for !this.Empty() {
        f(this.Pop())
    }
}

func (this *PointQueue) String() string {
    return fmt.Sprintf("PointQueue{size=%v}", this.Size())
}

/*
const POINTQUEUE_CAPACITY = MAX_ROWS * MAX_COLS

type PointQueue struct {
    points [POINTQUEUE_CAPACITY]Point
    start, end int
}

func (this *PointQueue) Push(p Point) {
    this.points[this.end] = p
    this.end += 1
    this.end %= POINTQUEUE_CAPACITY
    if this.end == this.start {
        panic("PointQueue ran out of capacity")
    }
}

func (this *PointQueue) Unshift(p Point) {
    this.start -= 1
    this.start += POINTQUEUE_CAPACITY
    this.start %= POINTQUEUE_CAPACITY
    if this.end == this.start {
        panic("PointQueue ran out of capacity")
    }
    this.points[this.start] = p
}

func (this *PointQueue) ForEach(f func(Point)) {
    for this.start != this.end {
        f(this.points[this.start])
        if this.start == this.end {
            break
        }
        this.start += 1
        this.start %= POINTQUEUE_CAPACITY
    }
}

func (this *PointQueue) Clear() {
    this.start = this.end
}

func (this *PointQueue) Empty() bool {
    return (this.end == this.start)
}

func (this *PointQueue) Size() int {
    return (this.end - this.start + POINTQUEUE_CAPACITY) % POINTQUEUE_CAPACITY
}
*/
