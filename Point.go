package main

import "math"
import "fmt"

type Point struct {
    row, col int
}

func (this Point) Normalize() (result Point) {
    result.row = normalize_row(this.row)
    result.col = normalize_col(this.col)
    return
}

func (this Point) Plus(p Point) (result Point) {
    result.row = normalize_row(this.row + p.row)
    result.col = normalize_col(this.col + p.col)
    return
}

func (this Point) Equals(p Point) bool {
    return (this.row == p.row && this.col == p.col)
}

func (this Point) IsOrigin() bool {
    return (this.row == 0 && this.col == 0)
}

func (this Point) Distance2(p Point) int {
    var abs1, abs2, dr, dc int

    abs1 = this.row - p.row
    if (abs1 < 0) {
        abs1 = -abs1
    }

    abs2 = rows - abs1
    if (abs2 < 0) {
        abs2 = -abs2
    }

    if abs1 < abs2 {
        dr = abs1
    } else {
        dr = abs2
    }

    abs1 = this.col - p.col
    if (abs1 < 0) {
        abs1 = -abs1
    }

    abs2 = cols - abs1
    if (abs2 < 0) {
        abs2 = -abs2
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
        this.row = normalize_row(this.row - 1)
    case EAST:
        this.col = normalize_col(this.col + 1)
    case SOUTH:
        this.row = normalize_row(this.row + 1)
    case WEST:
        this.col = normalize_col(this.col - 1)
    }
    return this
}

func (this Point) String() string {
    return fmt.Sprintf("Point{%v, %v}", this.row, this.col)
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
            p2 := p.Plus(d)
            f(p2)
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
/*
    var dir Direction
    for dir = 1; dir <= STAY; dir *= 2 {
        p2 := p.Neighbor(dir)
        f(p2)
    }
*/
    ForEachDirection(func(dir Direction) {
        f(p.Neighbor(dir))
    })
}
