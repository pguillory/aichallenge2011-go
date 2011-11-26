package main

type PotentialEnemy struct {
    time int64
    turn int
    terrain *Terrain
    values [MAX_ROWS][MAX_COLS]bool
}

func NewPotentialEnemy(terrain *Terrain) *PotentialEnemy {
    this := new(PotentialEnemy)
    this.terrain = terrain

    ForEachPoint(func(p Point) {
        square := this.terrain.At(p)
        if !square.IsVisible() || square.HasEnemyAnt() {
            this.values[p.row][p.col] = true
        }
    })

    return this
}

func (this *PotentialEnemy) Calculate() {
    if this.turn == turn {
        return
    }
    startTime := now()

    var newValues [MAX_ROWS][MAX_COLS]bool

    ForEachPoint(func(p Point) {
        square := this.terrain.At(p)
        if square.IsVisible() {
            newValues[p.row][p.col] = square.HasEnemyAnt()
        } else {
            switch {
            case this.At(p):
                newValues[p.row][p.col] = true
            case square.HasWater():
                newValues[p.row][p.col] = false
            default:
                newValues[p.row][p.col] = false
                ForEachNeighbor(p, func(p2 Point) {
                    if this.At(p2) {
                        newValues[p.row][p.col] = true
                    }
                })
            }
        }
    })

    this.values = newValues

    this.time = now() - startTime
    this.turn = turn
}

func (this *PotentialEnemy) At(p Point) bool {
    return this.values[p.row][p.col]
}

func (this *PotentialEnemy) String() string {
    return GridToString(func(p Point) byte {
        switch {
        case this.terrain.At(p).HasWater():
            return '%'
        case this.terrain.At(p).HasFriendlyAnt():
            return 'a'
        case this.At(p):
            return 'b'
        }
        return '.'
    })
}
