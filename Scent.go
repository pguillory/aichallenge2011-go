package main

import "fmt"

type ScentChannel struct {
    pathways [DIRECTIONS]*float32
    multiplier float32
    additive float32
}

type Scent struct {
    time int64
    turn int
    terrain *Terrain
    distanceToEnemy, distanceToFriendlyHill *TravelDistance
    mystery *Mystery
    configureChannel func(*Scent, Point, *ScentChannel)
    adjacentWater *AdjacentWater
    value [MAX_ROWS][MAX_COLS]float32
    channels [MAX_ROWS][MAX_COLS]ScentChannel
}

func NewForageScent(terrain *Terrain, distanceToEnemy *TravelDistance, distanceToFriendlyHill *TravelDistance, mystery *Mystery) *Scent {
    return NewScent(terrain, distanceToEnemy, distanceToFriendlyHill, mystery, ConfigureForageScentChannel)
}

func NewBattleScent(terrain *Terrain, distanceToEnemy *TravelDistance, distanceToFriendlyHill *TravelDistance, mystery *Mystery) *Scent {
    return NewScent(terrain, distanceToEnemy, distanceToFriendlyHill, mystery, ConfigureBattleScentChannel)
}

func NewScent(terrain *Terrain, distanceToEnemy *TravelDistance, distanceToFriendlyHill *TravelDistance, mystery *Mystery, configureChannel func(*Scent, Point, *ScentChannel)) *Scent {
    this := new(Scent)
    this.terrain = terrain
    this.distanceToEnemy = distanceToEnemy
    this.distanceToFriendlyHill = distanceToFriendlyHill
    this.mystery = mystery
    this.configureChannel = configureChannel
    this.adjacentWater = NewAdjacentWater(terrain)

    this.Calculate()
    return this
}

func (this *Scent) Calculate() {
    this.distanceToEnemy.Calculate()
    this.distanceToFriendlyHill.Calculate()
    this.mystery.Calculate()

    if this.turn == turn {
        return
    }
    startTime := now()

    this.BuildChannels()

    for i := 0; i < 100; i++ {
        this.Spread()
    }

    this.time = now() - startTime
    this.turn = turn
}

func (this *Scent) BuildChannels() {
    ForEachPoint(func(p Point) {
        channel := &this.channels[p.row][p.col]
        channel.pathways[0] = &this.value[(p.row - 1 + rows) % rows][(p.col           )       ]
        channel.pathways[1] = &this.value[(p.row           )       ][(p.col - 1 + cols) % cols]
        channel.pathways[2] = &this.value[(p.row           )       ][(p.col           )       ]
        channel.pathways[3] = &this.value[(p.row           )       ][(p.col + 1       ) % cols]
        channel.pathways[4] = &this.value[(p.row + 1       ) % rows][(p.col           )       ]

        channel.multiplier = 1.0 / (5.0 - float32(this.adjacentWater.At(p)))
        //channel.multiplier = 0.2
        channel.additive = 0

        this.configureChannel(this, p, channel)
    })
}

func ConfigureForageScentChannel(this *Scent, p Point, channel *ScentChannel) {
    this.value[p.row][p.col] *= 1e-1

    //channel.multiplier *= 0.95

    s := this.terrain.At(p)

    switch {
    case s.HasFood():
        channel.additive = 0.1
    case s.HasHill():
        if s.IsEnemy() {
            channel.additive = 1.0
        } else {
            channel.multiplier = 0.0
        }
    case s.HasAnt():
        if s.IsEnemy() {
            switch {
            case this.distanceToFriendlyHill.At(p) < 20:
                channel.additive = 1.0 * float32(20 - this.distanceToFriendlyHill.At(p)) / 20.0
            case this.terrain.VisibleFriendliesAt(p) >= 2:
                channel.additive = 0.01
            }
        } else {
            channel.multiplier = 0.0
        }
    case s.HasWater():
        channel.multiplier = 0.0
    default:
        channel.additive = this.mystery.At(p) * 0.002
    }
}

func ConfigureBattleScentChannel(this *Scent, p Point, channel *ScentChannel) {
    this.value[p.row][p.col] *= 1e-1

    //channel.multiplier *= 0.99

    s := this.terrain.At(p)

    switch {
    case s.HasHill():
        if s.IsEnemy() {
            channel.additive = 1.0
        } else {
            channel.multiplier = 0.0
        }
    case s.HasAnt():
        if s.IsEnemy() {
            switch {
            case this.distanceToFriendlyHill.At(p) < 20:
                channel.additive = 1.0 * float32(20 - this.distanceToFriendlyHill.At(p)) / 20.0
            default:
                channel.additive = 0.01
            }
        } else {
            channel.multiplier *= 1e-2
        }
    case s.HasWater():
        channel.multiplier = 0.0
    default:
        if this.mystery.At(p) > STARTING_MYSTERY {
            channel.additive = this.mystery.At(p) * 0.002
        }
    }
}

func (this *Scent) Spread() {
    var newValue [MAX_ROWS][MAX_COLS]float32

    var p Point
    for p.row = 0; p.row < rows; p.row++ {
        for p.col = 0; p.col < cols; p.col++ {
            channel := &this.channels[p.row][p.col]

            v := *channel.pathways[0] +
                *channel.pathways[1] +
                *channel.pathways[2] +
                *channel.pathways[3] +
                *channel.pathways[4]

            newValue[p.row][p.col] = v * channel.multiplier + channel.additive //* 1e20
        }
    }

    this.value = newValue
}

func (this *Scent) At(p Point) float32 {
    //return this.value[p.row][p.col] + float32(MAX_TRAVEL_DISTANCE - this.distanceToEnemy.At(p)) * 1e-30
    return this.value[p.row][p.col]
}

func (this *Scent) String() string {
    return GridToString(func(p Point) byte {
        square := this.terrain.At(p)
        switch {
        //case square.HasFood():
        //    return '*'
        case square.HasLand():
            switch {
            case this.At(p) < 0.0: return '-'
            case this.At(p) == 0.0: return '0'
            case this.At(p) <= 1e0: return '1'
            case this.At(p) <= 1e1: return '2'
            case this.At(p) <= 1e2: return '3'
            case this.At(p) <= 1e3: return '4'
            case this.At(p) <= 1e4: return '5'
            case this.At(p) <= 1e5: return '6'
            case this.At(p) <= 1e6: return '7'
            case this.At(p) <= 1e7: return '8'
            case this.At(p) <= 1e8: return '9'
            case this.At(p) <= 1e9: return 'a'
            case this.At(p) <= 1e10: return 'b'
            case this.At(p) <= 1e11: return 'c'
            case this.At(p) <= 1e12: return 'd'
            case this.At(p) <= 1e13: return 'e'
            case this.At(p) <= 1e14: return 'f'
            }
            return '+'
        case square.HasWater():
            return '%'
        }
        return ' '
    })
}

func (this *Scent) Csv() string {
    return GridToCsv(func(p Point) string {
        return fmt.Sprintf("%v", this.At(p))
    })
}
