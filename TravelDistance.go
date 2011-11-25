package main

//import "fmt"

const MAX_TRAVEL_DISTANCE = MAX_ROWS * MAX_COLS

type TravelDistance struct {
    time int64
    turn int
    terrain *Terrain
    value [MAX_ROWS][MAX_COLS]int
    isDestination func(p Point) bool
}

func DistanceToFriendlyHill(terrain *Terrain) *TravelDistance {
    return NewTravelDistance(terrain, func(p Point) bool {
        return terrain.At(p).HasFriendlyHill()
    })
}

func DistanceToEnemy(terrain *Terrain) *TravelDistance {
    return NewTravelDistance(terrain, func(p Point) bool {
        square := terrain.At(p)
        return square.HasEnemyAnt() || square.HasEnemyHill() || !square.EverSeen()
    })
}

func NewTravelDistance(terrain *Terrain, isDestination func(p Point) bool) *TravelDistance {
    this := new(TravelDistance)
    this.terrain = terrain
    this.isDestination = isDestination

    this.Calculate()
    return this
}

func (this *TravelDistance) Calculate() {
    if this.turn == turn {
        return
    }
    startTime := now()

    queue := new(PointQueue)

    ForEachPoint(func(p Point) {
        if this.isDestination(p) {
            this.value[p.row][p.col] = 0
            queue.Push(p)
        } else {
            this.value[p.row][p.col] = MAX_TRAVEL_DISTANCE
        }
    })

    queue.ForEach(func(p Point) {
        v := this.value[p.row][p.col]
        ForEachNeighbor(p, func(p2 Point) {
            if this.value[p2.row][p2.col] > v + 1 && this.terrain.At(p2).HasLand() {
                this.value[p2.row][p2.col] = v + 1
                queue.Push(p2)
            }
        })
    })

    this.time = now() - startTime
    this.turn = turn
}

func (this *TravelDistance) At(p Point) int {
    return this.value[p.row][p.col]
}

func (this *TravelDistance) String() string {
    return GridToString(func(p Point) byte {
        v := this.At(p)
        switch {
        case v < 10:
            return '0' + byte(v)
        case v < 36:
            return 'a' + byte(v) - 10
        }
        return '+'
    })
}
