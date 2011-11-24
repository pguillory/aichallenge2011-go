package main

type Army struct {
    terrain *Terrain
    values [MAX_ROWS][MAX_COLS]uint16
    counts [MAX_ROWS * MAX_COLS / 2]uint16
}

func NewArmy(terrain *Terrain) *Army {
    this := new(Army)
    this.terrain = terrain
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

func (this *Army) IsBerzerkerAt(p Point) bool {
    return (this.CountAt(p) >= 15)
    //return (this.CountAt(p) >= 20 && this.terrain.EnemiesVisibleFrom(p) < this.terrain.AlliesVisibleFrom(p))
}

func (this *Army) Spread(p Point) {
    ForEachPointWithinManhattanDistance(p, 1, func(p2 Point) {
        v, v2 := this.At(p), this.At(p2)
        if 0 < v && v < v2 {
            this.values[p2.row][p2.col] = v
            this.Spread(p2)
        }
    })
}

func (this *Army) Calculate() {
    var a uint16

    ForEachPoint(func(p Point) {
        if this.terrain.At(p).HasFriendlyAnt() {
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

func (this *Army) Berzerkers() *PointSet{
    points := new(PointSet)

    ForEachPoint(func(p Point) {
        if this.IsBerzerkerAt(p) {
            points.Include(p)
        }
    })

    return points
}

func (this *Army) String() string {
    return GridToString(func(p Point) byte {
        v := this.At(p)
        switch {
        case v == 0:
            return '.'
        case v <= 26:
            if this.IsSoldierAt(p) {
                return 'A' + byte(v) - 1
            } else {
                return 'a' + byte(v) - 1
            }
        }
        return '+'
    })
}
