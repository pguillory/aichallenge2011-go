package main

type Scent struct {
    terrain *Terrain
    mystery *Mystery
    value [MAX_ROWS][MAX_COLS]float32
}

func NewScent(terrain *Terrain, mystery *Mystery) *Scent {
    this := new(Scent)
    this.terrain = terrain
    this.mystery = mystery
    return this
}

func (this *Scent) At(p Point) float32 {
    return this.value[p.row][p.col]
}

func (this *Scent) Calculate() {
    var newValue [MAX_ROWS][MAX_COLS]float32

    ForEachPoint(func(p Point) {
        var v float32

        s := this.terrain.At(p)
        switch {
        case s.HasWater():
            v = 0.0
        case s.HasFriendlyAnt():
            v = 0.0
        case s.HasFriendlyHill():
            v = 0.0
        default:
            v = (this.value[(p.row - 1 + rows) % rows][(p.col           )       ] +
                 this.value[(p.row           )       ][(p.col - 1 + cols) % cols] +
                 this.value[(p.row           )       ][(p.col           )       ] +
                 this.value[(p.row           )       ][(p.col + 1       ) % cols] +
                 this.value[(p.row + 1       ) % rows][(p.col           )       ]) / 5.0 * 0.95

            v += this.mystery.At(p) * 10.0

            if s.HasFood() {
                v += 100.0
            } else if s.HasEnemyHill() {
                v += 500.0
            } else if s.HasAnt() {
                if s.IsEnemy() {
                    v += 5.0
                } else {
                    //ForEachNeighbor(p, func(p2 Point) {
                    //    
                    //})
                }
            }
        }

        newValue[p.row][p.col] = v
    })

    this.value = newValue
}

func (this *Scent) CalculateSoldier() {
    var newValue [MAX_ROWS][MAX_COLS]float32

    ForEachPoint(func(p Point) {
        var v float32

        s := this.terrain.At(p)
        switch {
        case s.HasWater():
            v = 0.0
        case s.HasFriendlyHill():
            v = 0.0
        default:
            v = (this.value[(p.row - 1 + rows) % rows][(p.col           )       ] +
                 this.value[(p.row           )       ][(p.col - 1 + cols) % cols] +
                 this.value[(p.row           )       ][(p.col           )       ] +
                 this.value[(p.row           )       ][(p.col + 1       ) % cols] +
                 this.value[(p.row + 1       ) % rows][(p.col           )       ]) / 5.0 * 0.95

            v += this.mystery.At(p) * 10.0

            if s.HasEnemyHill() {
                v += 500.0
            } else if s.HasAnt() {
                if s.IsEnemy() {
                    v += 5.0
                } else {
                    //ForEachNeighbor(p, func(p2 Point) {
                    //    
                    //})
                }
            }
        }

        newValue[p.row][p.col] = v
    })

    this.value = newValue
}

func (this *Scent) String() string {
    return GridToString(func(p Point) byte {
        square := this.terrain.At(p)
        switch {
        case square.HasFood():
            return '*'
        case square.HasLand():
            switch {
            case this.At(p) <     0: return '-'
            case this.At(p) <     1: return '0'
            case this.At(p) <     2: return '1'
            case this.At(p) <     4: return '2'
            case this.At(p) <     8: return '3'
            case this.At(p) <    16: return '4'
            case this.At(p) <    32: return '5'
            case this.At(p) <    64: return '6'
            case this.At(p) <   128: return '7'
            case this.At(p) <   256: return '8'
            case this.At(p) <   512: return '9'
            case this.At(p) <  1024: return 'a'
            case this.At(p) <  2048: return 'b'
            case this.At(p) <  4096: return 'c'
            case this.At(p) <  8192: return 'd'
            case this.At(p) < 16384: return 'e'
            case this.At(p) < 32768: return 'f'
            }
            return '+'
        case square.HasWater():
            return '%'
        }
        return '?'
    })
}
