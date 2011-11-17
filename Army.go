package main

type Army struct {
    m *Map
    values [MAX_ROWS][MAX_COLS]uint16
    counts [MAX_ROWS * MAX_COLS / 2]uint16
}

func NewArmy(m *Map) *Army {
    this := new(Army)
    this.m = m
    return this
}

func (this *Army) At(p Point) uint16 {
    return this.values[p.row][p.col]
}

func (this *Army) Count(a uint16) uint16 {
    return this.counts[a]
}

func (this *Army) CountAt(p Point) uint16 {
    return this.Count(this.At(p))
}

func (this *Army) IsSoldierAt(p Point) bool {
    return (this.CountAt(p) >= 3)
}

func (this *Army) Spread(p Point) {
    ForEachPointWithinManhattanDistance(p, 1, func(p2 Point) {
        v, v2 := this.At(p), this.At(p2)
        if v < v2 && v > 0 {
            this.values[p2.row][p2.col] = v
            this.Spread(p2)
        }
    })
}

func (this *Army) Iterate() {
    var a uint16

    ForEachPoint(func(p Point) {
        if this.m.At(p).HasFriendlyAnt() {
            a += 1
            this.values[p.row][p.col] = a
        } else {
            this.values[p.row][p.col] = 0
        }
    })

    ForEachPoint(func(p Point) {
        this.Spread(p)
    })

    var counts [MAX_ROWS * MAX_COLS / 2]uint16
    ForEachPoint(func(p Point) {
        counts[this.At(p)] += 1
    })
    this.counts = counts
}

func (this *Army) String() string {
    return GridToString(func(p Point) byte {
        v := this.At(p)
        switch {
        case v == 0:
        case v <= 26:
            if this.IsSoldierAt(p) {
                return 'A' + byte(v) - 1
            } else {
                return 'a' + byte(v) - 1
            }
        default:
            return '+'
        }
        return '.'
    })
}
