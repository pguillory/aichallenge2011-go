package main

import "math"
import "fmt"

type Point struct {
    row, col int
}

/*
func (this Point) AssertValid() {
    if this.row < 0 {
        fmt.Printf("row: %v or %#v\n", this.row, this.row)
        panic("Row too low")
    }
    if this.row >= MAX_ROWS {
        panic("Row too high")
    }
    if this.col < 0 {
        panic("Col too low")
    }
    if this.col >= MAX_COLS {
        panic("Col too high")
    }
}
*/

func (this Point) Normalize() (result Point) {
    result.row = normalizeRow(this.row)
    result.col = normalizeCol(this.col)
    return
}

func (this Point) Plus(p Point) (result Point) {
    result.row = (this.row + p.row + rows * MAX_ROWS) % rows
    result.col = (this.col + p.col + cols * MAX_COLS) % cols
    return
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
    case STAY:
    default:
        panic("Can't find neighbor to " + dir.String())
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

/*
// Radius lookup tables. Only slightly faster -- 23.8s vs 24.3s.  Not worth the complexity.
var radiusTables [200][]Point

func PrepareRadius2Tables(radius2 int) {
    distance := int(math.Ceil(math.Sqrt(float64(radius2))))

    radiusTables[radius2] = make([]Point, rows * cols, rows * cols)
    count := 0

    var d Point
    for d.row = -distance; d.row <= distance; d.row++ {
        for d.col = -distance; d.col <= distance; d.col++ {
            if d.row * d.row + d.col * d.col <= radius2 {
                radiusTables[radius2][count] = d
                count += 1
            }
        }
    }

    radiusTables[radius2] = radiusTables[radius2][0:count]
}

func ForEachPointWithinRadius2(p Point, radius2 int, f func(Point)) {
    if len(radiusTables[radius2]) == 0 {
        PrepareRadius2Tables(radius2)
    }

    for _, d := range radiusTables[radius2] {
        f(p.Plus(d))
    }
}
*/

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
