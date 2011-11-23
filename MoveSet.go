package main

type MoveSet struct {
    dirs [MAX_ROWS][MAX_COLS]Direction
}

func (this *MoveSet) At(p Point) Direction {
    return this.dirs[p.row][p.col]
}

func (this *MoveSet) Include(move Move) {
    this.dirs[move.from.row][move.from.col] |= move.dir
}

func (this *MoveSet) IncludeAllFrom(p Point) {
    this.dirs[p.row][p.col] = NORTH | EAST | SOUTH | WEST | STAY
}

func (this *MoveSet) Includes(move Move) bool {
    return this.dirs[move.from.row][move.from.col] & move.dir > 0
}

func (this *MoveSet) ExcludeAllFrom(p Point) {
    this.dirs[p.row][p.col] = 0
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

func (this *MoveSet) ForEachMoveTo(p Point, f func(Move)) {
    this.ForEach(func(move Move) {
        if move.Destination().Equals(p) {
            f(move)
        }
    })
}

func (this *MoveSet) FocusOn(p Point) (result byte) {
    moves := *this
    //ForEachPointWithinRadius2(p, attackradius2, func(p2 Point) {
    //    this.ForEachMoveTo(p2, func(move Move) {
    //        moves.Include(move)
    //    })
    //})

    var dirs [5]Direction

    for _, move := range spiralPattern {
        p2 := p.Plus(move.from)
        dirs[0] = STAY
        dirs[1] = move.dir
        dirs[2] = move.dir.Right()
        dirs[3] = move.dir.Backward()
        dirs[4] = move.dir.Left()
        for _, dir := range dirs {
            p3 := p2.Neighbor(dir)
            move2 := Move{p3, dir.Backward()}
            if moves.Includes(move2) {
                moves.Select(move2)
                result += 1
                break
            }
        }
    }

    return
}

func (this *MoveSet) String() string {
    return GridToString(func(p Point) byte {
        dir := this.At(p)

        switch {
        case dir.IsSingle():
            return dir.Char()
        case dir.IsMultiple():
            return '+'
        }

        return '.'
    })
}

func (this *MoveSet) Cardinality() int {
    count := 0
    this.ForEach(func(move Move) {
        count += 1
    })
    return count
}

func (this *MoveSet) OrderedList(valueFunc func(move Move) float32) *OrderedMoveList {
    list := NewOrderedMoveList(this.Cardinality())

    this.ForEach(func(move Move) {
        list.Add(move, valueFunc(move))
    })

    return list
}

func (this *MoveSet) ExceptFrom(exceptions *PointSet) *MoveSet {
    result := new(MoveSet)
    *result = *this
    exceptions.ForEach(func(p Point) {
        result.ExcludeAllFrom(p)
    })
    return result
}

func (this *MoveSet) Destinations() *PointSet {
    result := new(PointSet)

    this.ForEach(func(move Move) {
        result.Include(move.Destination())
    })

    return result
}
