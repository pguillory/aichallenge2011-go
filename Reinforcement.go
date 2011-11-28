// TODO
// rename to "Reinforcement"

package main

type Reinforcement struct {
    time int64
    turn int
    terrain *Terrain
    army *Army
    distanceToTrouble *TravelDistance
    values [MAX_ROWS][MAX_COLS]bool
}

func NewReinforcement(terrain *Terrain, army *Army, distanceToTrouble *TravelDistance) *Reinforcement {
    this := new(Reinforcement)
    this.terrain = terrain
    this.army = army
    this.distanceToTrouble = distanceToTrouble

    this.Calculate()
    return this
}

func (this *Reinforcement) ShouldBeInfectedAt(p Point) bool {
    switch {
    case this.distanceToTrouble.At(p) > 50 && !this.terrain.At(p).HasFriendlyHill():
        return true
    }
    return false
}

func (this *Reinforcement) ShouldBeCuredAt(p Point) bool {
    switch {
    case this.distanceToTrouble.At(p) < 40:
        return true
    }
    return false
}

func (this *Reinforcement) NeighborInfected(p Point) bool {
    result := false
    ForEachNeighbor(p, func(p2 Point) {
        if this.InfectedAt(p2) && !this.ShouldBeCuredAt(p2) {
            result = true
        }
    })
    return result
}

func (this *Reinforcement) Calculate() {
    if this.turn == turn {
        return
    }
    startTime := now()

    var newValues [MAX_ROWS][MAX_COLS]bool

    ForEachPoint(func(p Point) {
        switch {
        case this.terrain.At(p).HasFriendlyAnt() == false:
            newValues[p.row][p.col] = false
        case this.InfectedAt(p):
            if this.ShouldBeCuredAt(p) {
                newValues[p.row][p.col] = false
            } else {
                newValues[p.row][p.col] = true
            }
        case this.ShouldBeInfectedAt(p):
            newValues[p.row][p.col] = true
        default:
            newValues[p.row][p.col] = this.NeighborInfected(p)
        }
    })

    this.values = newValues

    this.time = now() - startTime
    this.turn = turn
}

func (this *Reinforcement) InfectedAt(p Point) bool {
    return this.values[p.row][p.col]
}

func (this *Reinforcement) String() string {
    return GridToString(func(p Point) byte {
        switch {
        case this.terrain.At(p).HasWater():
            return '%'
        case this.terrain.At(p).HasFriendlyAnt():
            if this.InfectedAt(p) {
                return 'A'
            } else {
                return 'a'
            }
        case this.terrain.At(p).HasEnemyAnt():
            return 'b'
        }
        return '.'
    })
}
