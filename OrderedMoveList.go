package main

import "fmt"
import "sort"

type MoveValue struct {
    move Move
    value float32
}

func (this MoveValue) String() string {
    return fmt.Sprintf("\n{%v:%v %c %v}", this.move.from.row, this.move.from.col, this.move.dir.Char(), this.value)
}



type OrderedMoveList struct {
    slice []MoveValue
}

func NewOrderedMoveList(capacity int) *OrderedMoveList {
    this := new(OrderedMoveList)
    this.slice = make([]MoveValue, 0, capacity)
    return this
}

// implements sort.Interface
func (this *OrderedMoveList) Len() int {
    return len(this.slice)
}

func (this *OrderedMoveList) Less(i, j int) bool {
    return this.slice[i].value < this.slice[j].value
}

func (this *OrderedMoveList) Swap(i, j int) {
    this.slice[i], this.slice[j] = this.slice[j], this.slice[i]
}

func (this *OrderedMoveList) Add(move Move, value float32) {
    this.slice = append(this.slice, MoveValue{move, value})
}

func (this *OrderedMoveList) ForBestWorst(valid func(Move) bool, bestFunc func(Move), worstFunc func(Move)) {
    sort.Sort(this)

    i := 0
    j := len(this.slice) - 1
    for i < j {
        for i < j && !valid(this.slice[j].move) {
            j--
        }
        if i < j {
            bestFunc(this.slice[j].move)
            j--
        } else {
            break
        }

        for i < j && !valid(this.slice[i].move) {
            i++
        }
        if i < j {
            worstFunc(this.slice[i].move)
            i++
        }
    }
}
