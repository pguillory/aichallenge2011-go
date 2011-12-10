package main

import "fmt"

const SEARCHQUEUEBLOCK_CAPACITY = 100000

type GoalPoint struct {
    goal *Goal
    point Point
}

type SearchQueueBlock struct {
    members [SEARCHQUEUEBLOCK_CAPACITY]GoalPoint
    start, end int
}

func (this *SearchQueueBlock) Push(p Point, goal *Goal) {
    this.members[this.end].goal = goal
    this.members[this.end].point = p
    this.end += 1
}

func (this *SearchQueueBlock) Pop() (Point, *Goal) {
    member := this.members[this.start]
    this.start += 1
    return member.point, member.goal
}

func (this *SearchQueueBlock) Empty() bool {
    return (this.start == this.end)
}


type SearchQueue struct {
    queues [MAX_TRAVEL_DISTANCE]*SearchQueueBlock
    nearest Distance
    length int
}

func (this *SearchQueue) Push(p Point, distance Distance, goal *Goal) {
    if distance < this.nearest {
        fmt.Println("Decreasing distance")
        panic("Decreasing distance")
    }
    if this.queues[distance] == nil {
        this.queues[distance] = new(SearchQueueBlock)
    }
    this.queues[distance].Push(p, goal)
    this.length += 1
}

func (this *SearchQueue) Pop() (Point, Distance, *Goal) {
    for this.queues[this.nearest] == nil || this.queues[this.nearest].Empty() {
        this.queues[this.nearest] = nil
        this.nearest++
    }
    p, goal := this.queues[this.nearest].Pop()
    this.length -= 1
    return p, this.nearest, goal
}

func (this *SearchQueue) Empty() bool {
    return this.length == 0
}
