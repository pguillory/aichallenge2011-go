/*
TODO: use PointQueue
*/

package main

type Army struct {
    time int64
    turn int
    terrain *Terrain
    values [MAX_ROWS][MAX_COLS]uint16
    counts [MAX_ROWS * MAX_COLS / 2]uint16
}

func NewArmy(terrain *Terrain) *Army {
    this := new(Army)
    this.terrain = terrain

    this.Calculate()
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

func (this *Army) IsScoutAt(p Point) bool {
    return (this.CountAt(p) < 15)
}

func (this *Army) IsSoldierAt(p Point) bool {
    return (this.CountAt(p) >= 15 && this.CountAt(p) < 30)
}

func (this *Army) IsBerzerkerAt(p Point) bool {
    return (this.CountAt(p) >= 30)
}

func (this *Army) Spread(p Point) {
    //ForEachNeighbor(p, func(p2 Point) {
    //ForEachPointWithinRadius2(p, 5, func(p2 Point) {
    v := this.At(p)
    if v > 0 {
        ForEachPointWithinManhattanDistance(p, 1, func(p2 Point) {
            if v < this.At(p2) {
                this.values[p2.row][p2.col] = v
                this.Spread(p2)
            }
        })
    }
}

func (this *Army) Calculate() {
    if this.turn == turn {
        return
    }
    startTime := now()

    var cohort uint16

    ForEachPoint(func(p Point) {
        if this.terrain.At(p).HasFriendlyAnt() {
            cohort += 1
            this.values[p.row][p.col] = cohort
        } else {
            this.values[p.row][p.col] = 0
        }
    })

    ForEachPoint(func(p Point) {
        this.Spread(p)
    })

    var counts [MAX_ROWS * MAX_COLS / 2]uint16
    ForEachPoint(func(p Point) {
        if this.At(p) > 0 {
            counts[this.At(p)] += 1
        }
    })
    this.counts = counts

    this.time = now() - startTime
    this.turn = turn
}

//func (this *Army) Berzerkers() *PointSet{
//    points := new(PointSet)
//
//    ForEachPoint(func(p Point) {
//        if this.IsBerzerkerAt(p) {
//            points.Include(p)
//        }
//    })
//
//    return points
//}

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
