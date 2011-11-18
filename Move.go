package main

type Move struct {
    from Point
    dir Direction
}

func (this Move) Destination() Point {
    return this.from.Neighbor(this.dir)
}

