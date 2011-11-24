package main

//import "fmt"

const POINTQUEUE_CAPACITY = 64 * 1024

type PointQueue struct {
    points [POINTQUEUE_CAPACITY]Point
    start, end int
}

func (this *PointQueue) Push(p Point) {
    this.points[this.end] = p
    if this.end < POINTQUEUE_CAPACITY {
        this.end += 1
    } else {
        this.end = 0
    }
    if this.end == this.start {
        panic("PointQueue ran out of capacity")
    }
}

func (this *PointQueue) ForEach(f func(Point)) {
    for i := this.start; i != this.end; i = (i + 1) % POINTQUEUE_CAPACITY {
        f(this.points[i])
    }
}
