package main

type Mystery struct {
    m *Map
    value [MAX_ROWS][MAX_COLS]float32
}

func NewMystery(m *Map) *Mystery {
    this := new(Mystery)
    this.m = m
    ForEachPoint(func(p Point) {
        if !m.At(p).IsVisible() {
            this.value[p.row][p.col] = 0.25
        }
    })
    return this
}

func (this *Mystery) At(p Point) float32 {
    return this.value[p.row][p.col]
}

func (this *Mystery) Iterate() {
    ForEachPoint(func(p Point) {
        if this.m.At(p).IsVisible() {
            this.value[p.row][p.col] = 0
        } else {
            this.value[p.row][p.col] += 0.005
            if this.value[p.row][p.col] > 1.0 {
                this.value[p.row][p.col] = 1.0
            }
        }
    })
}
