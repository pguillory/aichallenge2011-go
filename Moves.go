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
    scent *Scent
    dirs [MAX_ROWS][MAX_COLS]Direction
    //moves []Move
    len int
}

func NewMoves(m *Map, scent *Scent) *Moves {
    this := new(Moves)
    this.m = m
    this.scent = scent
    return this
}

func (this *Moves) At(p Point) Direction {
    return this.dirs[p.row][p.col]
}

/*
func (this *Moves) Include(p Point, dir Direction) {
    this.dirs[p.row][p.col] |= dir
}
*/

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
    ForEachDirection(func(dir Direction) {
        this.Exclude(p.Neighbor(dir), dir.Backward())
    })
}

func (this *Moves) PickBestMovesByScent() {
    var best_value, worst_value float32
    var best_p, worst_p Point
    var best_dir, worst_dir Direction

    for i := 0; i < rows * cols; i++ {
        found := false

        ForEachPoint(func(p Point) {
            switch this.At(p) {
            case 0:
            case NORTH:
            case EAST:
            case SOUTH:
            case WEST:
            case STAY:
            default:
                ForEachDirection(func(dir Direction) {
                    if this.At(p) & dir > 0 {
                        value := this.scent.At(p.Neighbor(dir)) - this.scent.At(p)

                        if !found || best_value < value {
                            best_value = value
                            best_p = p
                            best_dir = dir
                        }

                        if !found || worst_value > value {
                            worst_value = value
                            worst_p = p
                            worst_dir = dir
                        }

                        found = true
                    }
                })
            }
        })

        if found {
            this.Select(best_p, best_dir)
            this.Exclude(worst_p, worst_dir)
        } else {
            break
        }
    }
}

func (this *Moves) PickBestMovesByScent2() {
    moves := make(MoveSet, 0)

    ForEachPoint(func(p Point) {
        if this.At(p).IsMultiple() {
            ForEachDirection(func(dir Direction) {
                if this.At(p).Includes(dir) {
                    moves = append(moves, Move{p, dir, this.scent.At(p.Neighbor(dir)) - this.scent.At(p)})
                }
            })
        }
    })

    sort.Sort(moves)
    NewLog("moves", "log").TurnFile().WriteString(fmt.Sprintf("%v", moves))

    i := 0
    j := len(moves) - 1
    for i < j {
        for this.At(moves[j].p) & moves[j].dir == 0 {
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

        for this.At(moves[i].p) & moves[i].dir == 0 {
            i++
            if i > j {
                return
            }
        }
        this.Exclude(moves[i].p, moves[i].dir)
        i++

/*
        worst := moves[0]
        for this.At(worst.p) & worst.dir == 0 {
            moves = moves[1:]
            if len(moves) == 0 {
                break
            }
            worst = moves[0]
        }
        this.Exclude(worst.p, worst.dir)
        moves = moves[1:]
*/
    }
}

func (this *Moves) Calculate() {
    ForEachPoint(func(p Point) {
        this.dirs[p.row][p.col] = 0
        if this.m.At(p).HasFriendlyAnt() {
            ForEachDirection(func(dir Direction) {
                //this.Include(p, dir)
                this.dirs[p.row][p.col] |= dir
            })
        }
    })

    ForEachPoint(func(p Point) {
        if this.m.At(p).HasWater() || this.m.At(p).HasFood() {
            this.ExcludeMovesTo(p)
        }
    })

    this.PickBestMovesByScent2()
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
