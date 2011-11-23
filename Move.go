package main

import "fmt"

type Move struct {
    from Point
    dir Direction
}

func (this Move) Destination() Point {
    return this.from.Neighbor(this.dir)
}

func (this Move) String() string {
    return fmt.Sprintf("%v %c", this.from, this.dir.Char())
}
