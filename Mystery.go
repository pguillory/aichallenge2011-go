package main

const STARTING_MYSTERY = 0.25

type Mystery struct {
    time int64
    turn int
    terrain *Terrain
    value [MAX_ROWS][MAX_COLS]float32
}

func NewMystery(terrain *Terrain) *Mystery {
    this := new(Mystery)
    this.terrain = terrain
    ForEachPoint(func(p Point) {
        if !terrain.At(p).IsVisible() {
            this.value[p.row][p.col] = STARTING_MYSTERY
        }
    })

    this.Calculate()
    return this
}

func (this *Mystery) At(p Point) float32 {
    return this.value[p.row][p.col]
}

func (this *Mystery) Calculate() {
    if this.turn == turn {
        return
    }
    startTime := now()

    ForEachPoint(func(p Point) {
        square := this.terrain.At(p)
        if square.IsVisible() || square.HasWater() {
            this.value[p.row][p.col] = 0
        } else {
            this.value[p.row][p.col] += 0.001
            if this.value[p.row][p.col] > 1.0 {
                this.value[p.row][p.col] = 1.0
            }
        }
    })

    this.time = now() - startTime
    this.turn = turn
}

func (this *Mystery) String() string {
    return GridToString(func(p Point) byte {
        v := int(this.At(p) * 10)
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
