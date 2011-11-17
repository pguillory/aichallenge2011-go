package main

import "bytes"

type Scent struct {
    m *Map
    mystery *Mystery
    value [MAX_ROWS][MAX_COLS]float32
}

func NewScent(m *Map, mystery *Mystery) *Scent {
    this := new(Scent)
    this.m = m
    this.mystery = mystery
    return this
}

func (this *Scent) At(p Point) float32 {
    return this.value[p.row][p.col]
}

func (this *Scent) Iterate() {
    var newValue [MAX_ROWS][MAX_COLS]float32

    ForEachPoint(func(p Point) {
        var v float32

        s := this.m.At(p)
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

/*
func (this *Scent) Stabilize() {
    for i := 0; i < 100; i++ {
        scent.Iterate()
    }
}
*/

func (this *Scent) String() string {
    b := new(bytes.Buffer)
    max_row := 0

    ForEachPoint(func(p Point) {
        for max_row < p.row {
            max_row += 1
            b.WriteByte('\n')
        }

        square := this.m.At(p)
        switch {
        case square.HasFood():
            b.WriteByte('*')
        case square.HasLand():
            switch {
            case this.At(p) <     0: b.WriteByte('-')
            case this.At(p) <     1: b.WriteByte('0')
            case this.At(p) <     2: b.WriteByte('1')
            case this.At(p) <     4: b.WriteByte('2')
            case this.At(p) <     8: b.WriteByte('3')
            case this.At(p) <    16: b.WriteByte('4')
            case this.At(p) <    32: b.WriteByte('5')
            case this.At(p) <    64: b.WriteByte('6')
            case this.At(p) <   128: b.WriteByte('7')
            case this.At(p) <   256: b.WriteByte('8')
            case this.At(p) <   512: b.WriteByte('9')
            case this.At(p) <  1024: b.WriteByte('a')
            case this.At(p) <  2048: b.WriteByte('b')
            case this.At(p) <  4096: b.WriteByte('c')
            case this.At(p) <  8192: b.WriteByte('d')
            case this.At(p) < 16384: b.WriteByte('e')
            case this.At(p) < 32768: b.WriteByte('f')
            default:
                b.WriteByte('+')
            }
        case square.HasWater():
            b.WriteByte('%')
        default:
            b.WriteByte('?')
        }
    })

    return b.String()
}
