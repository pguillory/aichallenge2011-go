package main

type Move struct {
    from Point
    dir Direction
}

type MoveSet struct {
    dirs [MAX_ROWS][MAX_COLS]Direction
}

func (this *MoveSet) At(p Point) Direction {
    return this.dirs[p.row][p.col]
}

func (this *MoveSet) IncludeAllFrom(p Point) {
    this.dirs[p.row][p.col] = NORTH | EAST | SOUTH | WEST | STAY
}

func (this *MoveSet) Includes(move Move) bool {
    return this.dirs[move.from.row][move.from.col] & move.dir > 0
}

func (this *MoveSet) Exclude(move Move) {
    if this.Includes(move) {
        this.dirs[move.from.row][move.from.col] &= ^move.dir

        ForEachDirection(func(dir Direction) {
            if this.At(move.from) == dir {
                this.Select(Move{move.from, dir})
            }
        })
    }
}

func (this *MoveSet) Select(move Move) {
    this.dirs[move.from.row][move.from.col] = 0
    this.ExcludeMovesTo(move.from.Neighbor(move.dir))
    this.dirs[move.from.row][move.from.col] = move.dir
}

func (this *MoveSet) ExcludeMovesTo(p Point) {
    this.Exclude(Move{p, STAY})
    this.Exclude(Move{p.Neighbor(NORTH), SOUTH})
    this.Exclude(Move{p.Neighbor(EAST), WEST})
    this.Exclude(Move{p.Neighbor(SOUTH), NORTH})
    this.Exclude(Move{p.Neighbor(WEST), EAST})
}

func (this *MoveSet) ForEach(f func(Move)) {
    ForEachPoint(func(p Point) {
        if this.At(p) > 0 {
            ForEachDirection(func(dir Direction) {
                move := Move{p, dir}
                if this.Includes(move) {
                    f(move)
                }
            })
        }
    })
}

func (this *MoveSet) Cardinality() int {
    return MAX_ROWS * MAX_COLS
}

func (this *MoveSet) Order(valueFunc func(move Move) float32) *OrderedMoveList {
    list := NewOrderedMoveList(this.Cardinality())

    this.ForEach(func(move Move) {
        list.Add(move, valueFunc(move))
    })

    return list
}
