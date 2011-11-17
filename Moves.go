package main

import "sort"
import "fmt"
import "bytes"

type Move struct {
    p Point
    dir Direction
    value float32
}

func (this Move) String() string {
    return fmt.Sprintf("\n{%v:%v %c %v}", this.p.row, this.p.col, this.dir.Char(), this.value)
}

type MoveSet []Move

// implements sort.Interface
func (this MoveSet) Len() int {
    return len(this)
}

func (this MoveSet) Less(i, j int) bool {
    return this[i].value < this[j].value
}

func (this MoveSet) Swap(i, j int) {
    this[i], this[j] = this[j], this[i]
}

type Moves struct {
    m *Map
    workerScent, soldierScent *Scent
    army *Army
    dirs [MAX_ROWS][MAX_COLS]Direction
    len int
}

func NewMoves(m *Map, workerScent, soldierScent *Scent, army *Army) *Moves {
    this := new(Moves)
    this.m = m
    this.workerScent = workerScent
    this.soldierScent = soldierScent
    this.army = army
    return this
}

func (this *Moves) Reset() {
    ForEachPoint(func(p Point) {
        if this.m.At(p).HasFriendlyAnt() {
            this.dirs[p.row][p.col] = NORTH | EAST | SOUTH | WEST | STAY
        } else {
            this.dirs[p.row][p.col] = 0
        }
    })

    ForEachPoint(func(p Point) {
        if this.m.At(p).HasWater() || this.m.At(p).HasFood() {
            this.ExcludeMovesTo(p)
        }
    })
}

func (this *Moves) At(p Point) Direction {
    return this.dirs[p.row][p.col]
}

func (this *Moves) Exclude(p Point, dir Direction) {
    if this.dirs[p.row][p.col] & dir > 0 {
        this.dirs[p.row][p.col] &= ^dir

        ForEachDirection(func(dir Direction) {
            if this.dirs[p.row][p.col] == dir {
                this.Select(p, dir)
            }
        })
    }
}

func (this *Moves) Select(p Point, dir Direction) {
    this.dirs[p.row][p.col] = 0
    this.ExcludeMovesTo(p.Neighbor(dir))
    this.dirs[p.row][p.col] = dir
}

func (this *Moves) ExcludeMovesTo(p Point) {
    this.Exclude(p, STAY)
    this.Exclude(p.Neighbor(NORTH), SOUTH)
    this.Exclude(p.Neighbor(EAST), WEST)
    this.Exclude(p.Neighbor(SOUTH), NORTH)
    this.Exclude(p.Neighbor(WEST), EAST)
}

func (this *Moves) MoveValue(p Point, dir Direction) float32 {
    p2 := p.Neighbor(dir)

    if this.army.IsSoldierAt(p) {
        return this.soldierScent.At(p2) - this.soldierScent.At(p)
    }

    return this.workerScent.At(p2) - this.workerScent.At(p)
}

func (this *Moves) PickBestMovesByScent() {
    moves := make(MoveSet, 0)

    ForEachPoint(func(p Point) {
        if this.At(p).IsMultiple() {
            ForEachDirection(func(dir Direction) {
                if this.At(p).Includes(dir) {
                    moves = append(moves, Move{p, dir, this.MoveValue(p, dir)})
                }
            })
        }
    })

    sort.Sort(moves)
    NewLog("moves", "log").TurnFile().WriteString(fmt.Sprintf("%v", moves))

    i := 0
    j := len(moves) - 1
    for i < j {
        for this.At(moves[j].p).Includes(moves[j].dir) == false {
            j--
            if i > j {
                return
            }
        }
        this.Select(moves[j].p, moves[j].dir)
        j--
        if i > j {
            return
        }

        for this.At(moves[i].p).Includes(moves[i].dir) == false {
            i++
            if i > j {
                return
            }
        }
        this.Exclude(moves[i].p, moves[i].dir)
        i++
    }
}

func (this *Moves) Calculate() {
    this.Reset()

    this.PickBestMovesByScent()
}

func (this *Moves) String() string {
    b := new(bytes.Buffer)
    max_row := 0

    ForEachPoint(func(p Point) {
        for max_row < p.row {
            max_row += 1
            b.WriteByte('\n')
        }

        b.WriteByte(this.At(p).Char())
    })

    return b.String()
}
