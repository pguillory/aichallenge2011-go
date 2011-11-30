package main

//import "fmt"

type PointSet [MAX_ROWS][MAX_COLS]bool

func (this *PointSet) Include(p Point) {
    this[p.row][p.col] = true
}

func (this *PointSet) Exclude(p Point) {
    this[p.row][p.col] = false
}

func (this *PointSet) Includes(p Point) bool {
    return this[p.row][p.col]
}

func (this *PointSet) Minus(other *PointSet) *PointSet {
    result := new(PointSet)
    this.ForEach(func(p Point) {
        if !other.Includes(p) {
            result.Include(p)
        }
    })
    return result
}

func (this *PointSet) Intersection(other *PointSet) *PointSet {
    result := new(PointSet)
    ForEachPoint(func(p Point) {
        if this.Includes(p) && other.Includes(p) {
            result.Include(p)
        }
    })
    return result
}

func (this *PointSet) Union(other *PointSet) *PointSet {
    result := new(PointSet)
    ForEachPoint(func(p Point) {
        if this.Includes(p) || other.Includes(p) {
            result.Include(p)
        }
    })
    return result
}

func (this *PointSet) Visibility() *PointSet {
    result := new(PointSet)
    this.ForEach(func(p Point) {
        ForEachPointWithinRadius2(p, viewradius2, func(p2 Point) {
            result.Include(p2)
        })
    })
    return result
}

func (this *PointSet) ForEach(f func(Point)) {
    ForEachPoint(func(p Point) {
        if this.Includes(p) {
            f(p)
        }
    })
}

func (this *PointSet) Cardinality() int {
    result := 0

    ForEachPoint(func(p Point) {
        if this[p.row][p.col] {
            result += 1
        }
    })

    return result
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
// TODO: make this optimized version work
const POINTSET_BITSIZE = MAX_ROWS * MAX_COLS
const POINTSET_INTSIZE = (POINTSET_BITSIZE + 31) / 32

type PointSet struct {
    values [POINTSET_INTSIZE]uint32
}

func (this *PointSet) Include(p Point) {
    index := uint32(p.row * rows + p.col)
    offset := index / 32
    mask := uint32(1) << (index % 32)
    this.values[offset] |= mask
}

func (this *PointSet) Includes(p Point) bool {
    index := uint32(p.row * rows + p.col)
    offset := index / 32
    mask := uint32(1) << (index % 32)
    return this.values[offset] & mask > 0
}

func (this *PointSet) Minus(other *PointSet) *PointSet {
    result := new(PointSet)
    for offset := 0; offset < POINTSET_INTSIZE; offset++ {
        result.values[offset] = this.values[offset] & ^other.values[offset]
    }
    return result
}

func (this *PointSet) Intersection(other *PointSet) *PointSet {
    result := new(PointSet)
    for offset := 0; offset < POINTSET_INTSIZE; offset++ {
        result.values[offset] = this.values[offset] & other.values[offset]
    }
    return result
}

func (this *PointSet) Union(other *PointSet) *PointSet {
    result := new(PointSet)
    for offset := 0; offset < POINTSET_INTSIZE; offset++ {
        result.values[offset] = this.values[offset] | other.values[offset]
    }
    return result
}

func (this *PointSet) Visibility() *PointSet {
    result := new(PointSet)
    this.ForEach(func(p Point) {
        ForEachPointWithinRadius2(p, viewradius2, func(p2 Point) {
            result.Include(p2)
        })
    })
    return result
}

func (this *PointSet) ForEach(f func(Point)) {
    size := uint32(rows * cols)
    for index := uint32(0); index < size; index++ {
        offset := index / 32
        mask := uint32(1) << (index % 32)
        if this.values[offset] & mask > 0 {
            f(Point{int(index / uint32(rows)), int(index % uint32(rows))})
        }
    }
}

func (this *PointSet) Cardinality() int {
    result := 0

    this.ForEach(func(p Point) {
        result += 1
    })

    return result
}

func (this *PointSet) String() string {
    return GridToString(func(p Point) byte {
        if this.Includes(p) {
            return 'x'
        }
        return '.'
    })
}
*/
