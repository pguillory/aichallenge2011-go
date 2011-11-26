package main

const STARTING_MYSTERY = 100
const MAX_MYSTERY = 100

type Mystery struct {
    time int64
    turn int
    terrain *Terrain
    values [MAX_ROWS][MAX_COLS]byte
}

func NewMystery(terrain *Terrain) *Mystery {
    this := new(Mystery)
    this.terrain = terrain

    ForEachPoint(func(p Point) {
        if !this.terrain.At(p).IsVisible() {
            this.values[p.row][p.col] = STARTING_MYSTERY
        }
    })

    return this
}

func (this *Mystery) Calculate() {
    if this.turn == turn {
        return
    }
    startTime := now()

    ForEachPoint(func(p Point) {
        if this.terrain.At(p).IsVisible() {
            this.values[p.row][p.col] = 0
        } else {
            if this.values[p.row][p.col] < MAX_MYSTERY {
                this.values[p.row][p.col] += 1
            }
        }
    })

    this.time = now() - startTime
    this.turn = turn
}

func (this *Mystery) At(p Point) byte {
    return this.values[p.row][p.col]
}

func (this *Mystery) String() string {
    return GridToString(func(p Point) byte {
        v := this.At(p) / 10
        switch {
        case v < 0:
            return '-'
        case v < 10:
            return '0' + byte(v)
        case v == 10:
            return 'a'
        }
        return '+'
    })
}
