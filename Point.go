package main

import "math"
import "fmt"

type Point struct {
    row, col int
}

func (this Point) Normalize() Point {
    this.row = normalizeRow(this.row)
    this.col = normalizeCol(this.col)
    return this
}

func (this Point) Plus(p Point) Point {
    this.row = (this.row + p.row + rows) % rows
    this.col = (this.col + p.col + cols) % cols
    return this
}

func (this Point) Equals(p Point) bool {
    return (this.row == p.row && this.col == p.col)
}

func (this Point) Distance2(p Point) int {
    var abs1, abs2, dr, dc int

    if (this.row > p.row) {
        abs1 = this.row - p.row
    } else {
        abs1 = p.row - this.row
    }

    if (rows > abs1) {
        abs2 = rows - abs1
    } else {
        abs2 = abs1 - rows
    }

    if abs1 < abs2 {
        dr = abs1
    } else {
        dr = abs2
    }

    if (this.col > p.col) {
        abs1 = this.col - p.col
    } else {
        abs1 = p.col - this.col
    }

    if (cols > abs1) {
        abs2 = cols - abs1
    } else {
        abs2 = abs1 - cols
    }

    if abs1 < abs2 {
        dc = abs1
    } else {
        dc = abs2
    }

    return dr * dr + dc * dc
}

func (this Point) Neighbor(dir Direction) Point {
    switch dir {
    case NORTH:
        this.row = (this.row - 1 + rows) % rows
    case EAST:
        this.col = (this.col + 1       ) % cols
    case SOUTH:
        this.row = (this.row + 1       ) % rows
    case WEST:
        this.col = (this.col - 1 + cols) % cols
    }
    return this
}

func (this Point) String() string {
    return fmt.Sprintf("%v:%v", this.row, this.col)
}

func ForEachPoint(f func(Point)) {
    var p Point
    for p.row = 0; p.row < rows; p.row++ {
        for p.col = 0; p.col < cols; p.col++ {
            f(p)
        }
    }
}

func ForEachPointWithinManhattanDistance(p Point, distance int, f func(Point)) {
    var d, s Point
    for d.row, s.row = -distance, 0; d.row <= distance && s.row < rows; d.row, s.row = d.row + 1, s.row + 1 {
        for d.col, s.col = -distance, 0; d.col <= distance && s.col < cols; d.col, s.col = d.col + 1, s.col + 1 {
            f(p.Plus(d))
        }
    }
}

func ForEachPointWithinRadius2(p Point, radius2 int, f func(Point)) {
    distance := int(math.Ceil(math.Sqrt(float64(radius2))))
    ForEachPointWithinManhattanDistance(p, distance, func(p2 Point) {
        if p.Distance2(p2) <= radius2 {
            f(p2)
        }
    })
}

func ForEachNeighbor(p Point, f func(Point)) {
    f(p.Neighbor(NORTH))
    f(p.Neighbor(EAST))
    f(p.Neighbor(SOUTH))
    f(p.Neighbor(WEST))
    f(p.Neighbor(STAY))
}
