package main

type Focus [MAX_ROWS][MAX_COLS]byte

func OpposingFocus(destinations *PointSet, opposingMoves *MoveSet) *Focus {
    this := new(Focus)

    destinations.ForEach(func(p Point) {
        this[p.row][p.col] = opposingMoves.FocusOn(p)
    })

    return this
}

func MaxFocus(destinations *PointSet, focus *Focus) *Focus {
    this := new(Focus)

    ForEachPoint(func(p Point) {
        this[p.row][p.col] = 255
    })

    destinations.ForEach(func(p Point) {
        v := focus[p.row][p.col]
        //if v > 0
        ForEachPointWithinRadius2(p, attackradius2, func(p2 Point) {
            if this[p2.row][p2.col] > v {
                this[p2.row][p2.col] = v
            }
        })
    })

    return this
}

func (this *Focus) At(p Point) byte {
    return this[p.row][p.col]
}

func (this *Focus) String() string {
    return GridToString(func(p Point) byte {
        value := this.At(p)
        switch {
        case value == 0:
            return '.'
        case value < 10:
            return '0' + byte(value)
        case value < 36:
            return 'a' + byte(value - 36)
        //case value == 255:
        //    return '.'
        }
        return '+'
    })
}
