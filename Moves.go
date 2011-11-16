package main

import "bytes"

type Moves struct {
    dirs [MAX_ROWS][MAX_COLS]Direction
}

func (this *Moves) At(p Point) Direction {
    return this.dirs[p.row][p.col]
}

func (this *Moves) Include(p Point, dir Direction) {
    this.dirs[p.row][p.col] |= dir
}

func (this *Moves) Exclude(p Point, dir Direction) {
    this.dirs[p.row][p.col] &= ^dir

    ForEachDirection(func(dir Direction) {
        if this.dirs[p.row][p.col] == dir {
            this.Select(p, dir)
        }
    })
}

func (this *Moves) Select(p Point, dir Direction) {
    this.ExcludeMovesTo(p.Neighbor(dir))
    this.dirs[p.row][p.col] = dir
}

func (this *Moves) ExcludeMovesTo(p Point) {
    ForEachDirection(func(dir Direction) {
        this.Exclude(p.Neighbor(dir), dir.Backward())
    })
}

func (this *Moves) PickBestMove(scent *Scent) (best_p Point, best_dir Direction, worst_p Point, worst_dir Direction, count int) {
    var best_value, worst_value float64

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
                if this.At(p) & dir == dir {
                    value := scent.At(p.Neighbor(dir)) - scent.At(p)

                    if count == 0 || best_value < value {
                        best_value = value
                        best_p = p
                        best_dir = dir
                    }

                    if count == 0 || worst_value > value {
                        worst_value = value
                        worst_p = p
                        worst_dir = dir
                    }

                    count += 1
                }
            })
        }
    })

    return
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

func NewMoves(m *Map, scent *Scent) *Moves {
    this := new(Moves)

    ForEachPoint(func(p Point) {
        if m.At(p).HasFriendlyAnt() {
            ForEachDirection(func(dir Direction) {
                this.Include(p, dir)
            })
        }
    })

    ForEachPoint(func(p Point) {
        if m.At(p).HasWater() || m.At(p).HasFood() {
            this.ExcludeMovesTo(p)
        }
    })

    for {
        best_p, best_dir, worst_p, worst_dir, count := this.PickBestMove(scent)
        if count == 0 {
            break
        }
        this.Select(best_p, best_dir)
        this.Exclude(worst_p, worst_dir)
    }

    return this
}
