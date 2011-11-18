package main

type PointSet [MAX_ROWS][MAX_COLS]bool

func (this *PointSet) Include(p Point) {
    this[p.row][p.col] = true
}

func (this *PointSet) ForEach(f func(Point)) {
    ForEachPoint(func(p Point) {
        if this[p.row][p.col] {
            f(p)
        }
    })
}

func (this *PointSet) String() string {
    return GridToString(func(p Point) byte {
        if this[p.row][p.col] {
            return 'x'
        }
        return '.'
    })
}

/*
type PointSet struct {
    values [(MAX_ROWS * MAX_COLS + 31) / 32]uint32
}

func (this *PointSet) Include(p Point) {
    index := p.row * rows + p.col
    offset := index / 32
    mask := 1 << index % 32
    this.values[offset] |= mask
}

func (this *PointSet) ForEach(f func(Point)) {
    for index := 0; index < rows * cols; index++ {
        offset := index / 32
        mask := 1 << index % 32
        if this.values[offset] & mask > 0 {
            f(Point{index / rows, index % rows})
        }
    }
}
*/
