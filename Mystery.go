package main

type Mystery struct {
    terrain *Terrain
    value [MAX_ROWS][MAX_COLS]float32
}

func NewMystery(terrain *Terrain) *Mystery {
    this := new(Mystery)
    this.terrain = terrain
    ForEachPoint(func(p Point) {
        if !terrain.At(p).IsVisible() {
            this.value[p.row][p.col] = 0.25
        }
    })
    return this
}

func (this *Mystery) At(p Point) float32 {
    return this.value[p.row][p.col]
}

func (this *Mystery) Calculate() {
    ForEachPoint(func(p Point) {
        if this.terrain.At(p).IsVisible() {
            this.value[p.row][p.col] = 0
        } else {
            this.value[p.row][p.col] += 0.001
            if this.value[p.row][p.col] > 1.0 {
                this.value[p.row][p.col] = 1.0
            }
        }
    })
}
