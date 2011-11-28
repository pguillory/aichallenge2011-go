package main

//import "fmt"

type AdjacentWater struct {
    time int64
    turn int
    terrain *Terrain
    value [MAX_ROWS][MAX_COLS]int
}

func NewAdjacentWater(terrain *Terrain) *AdjacentWater {
    this := new(AdjacentWater)
    this.terrain = terrain

    this.Calculate()
    return this
}

func (this *AdjacentWater) Calculate() {
    if this.turn == turn {
        return
    }
    startTime := now()

    ForEachPoint(func(p Point) {
        count := 0
        ForEachNeighbor(p, func(p2 Point) {
            if this.terrain.At(p2).HasWater() {
                count++
            }
        })
        this.value[p.row][p.col] = count
    })

    this.time = now() - startTime
    this.turn = turn
}

func (this *AdjacentWater) At(p Point) int {
    return this.value[p.row][p.col]
}

func (this *AdjacentWater) String() string {
    return GridToString(func(p Point) byte {
        v := this.At(p)
        return '0' + byte(v)
    })
}
