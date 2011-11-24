package main

import "fmt"

type ScentChannel struct {
    pathways [DIRECTIONS]*float32
    multiplier float32
    additive float32
}

func ConfigureForageScentChannel(this *Scent, p Point, channel *ScentChannel) {
    channel.multiplier *= 0.95

    s := this.terrain.At(p)

    switch {
    case s.HasFood():
        channel.additive = 50.0
    case s.HasHill():
        if s.IsEnemy() {
            channel.additive = 500.0
        } else {
            channel.multiplier = 0.0
        }
    case s.HasAnt():
        if s.IsEnemy() {
            switch {
            case this.holyGround.At(p) < 20:
                channel.additive = 500.0 * (20.0 - float32(this.holyGround.At(p))) / 20.0
            case this.terrain.VisibleFriendliesAt(p) >= 2:
                channel.additive = 5.0
            }
        } else {
            channel.multiplier *= 0.1
        }
    case s.HasWater():
        channel.multiplier = 0.0
    default:
        channel.additive = this.mystery.At(p) * 5.0
    }
}

func ConfigureBattleScentChannel(this *Scent, p Point, channel *ScentChannel) {
    channel.multiplier *= 0.99

    s := this.terrain.At(p)

    switch {
    case s.HasHill():
        if s.IsEnemy() {
            channel.additive = 500.0
        } else {
            channel.multiplier = 0.0
        }
    case s.HasAnt():
        if s.IsEnemy() {
            switch {
            case this.holyGround.At(p) < 20:
                channel.additive = 500.0 * (20.0 - float32(this.holyGround.At(p))) / 20.0
            default:
                channel.additive = 5.0
            }
        } else {
            channel.multiplier *= 0.8
        }
    case s.HasWater():
        channel.multiplier = 0.0
    default:
        channel.additive = this.mystery.At(p) * 1.0
    }
}

type Scent struct {
    terrain *Terrain
    holyGround *HolyGround
    mystery *Mystery
    value [MAX_ROWS][MAX_COLS]float32
    channels [MAX_ROWS][MAX_COLS]ScentChannel
    configureChannel func(*Scent, Point, *ScentChannel)
}

func NewForageScent(terrain *Terrain, holyGround *HolyGround, mystery *Mystery) *Scent {
    this := new(Scent)
    this.terrain = terrain
    this.holyGround = holyGround
    this.mystery = mystery
    this.configureChannel = ConfigureForageScentChannel
    return this
}

func NewBattleScent(terrain *Terrain, holyGround *HolyGround, mystery *Mystery) *Scent {
    this := new(Scent)
    this.terrain = terrain
    this.holyGround = holyGround
    this.mystery = mystery
    this.configureChannel = ConfigureBattleScentChannel
    return this
}

func (this *Scent) At(p Point) float32 {
    return this.value[p.row][p.col]
}

func (this *Scent) BuildChannels() {
    ForEachPoint(func(p Point) {
        channel := &this.channels[p.row][p.col]
        channel.pathways[0] = &this.value[(p.row - 1 + rows) % rows][(p.col           )       ]
        channel.pathways[1] = &this.value[(p.row           )       ][(p.col - 1 + cols) % cols]
        channel.pathways[2] = &this.value[(p.row           )       ][(p.col           )       ]
        channel.pathways[3] = &this.value[(p.row           )       ][(p.col + 1       ) % cols]
        channel.pathways[4] = &this.value[(p.row + 1       ) % rows][(p.col           )       ]

        channel.multiplier = 1.0 / (5.0 - float32(this.terrain.waterNeighbors[p.row][p.col]))
        channel.additive = 0

        this.configureChannel(this, p, channel)
    })
}

/*
func (this *Scent) EmanateForage() {
    ForEachPoint(func(p Point) {
        channel := &this.channels[p.row][p.col]

        channel.multiplier *= 0.95

        s := this.terrain.At(p)

        switch {
        case s.HasFood():
            channel.additive = 50.0
        case s.HasHill():
            if s.IsEnemy() {
                channel.additive = 500.0
            } else {
                channel.multiplier = 0.0
            }
        case s.HasAnt():
            if s.IsEnemy() {
                switch {
                case this.holyGround.At(p) < 20:
                    channel.additive = 500.0 * (20.0 - float32(this.holyGround.At(p))) / 20.0
                case this.terrain.VisibleFriendliesAt(p) >= 2:
                    channel.additive = 5.0
                }
            } else {
                channel.multiplier *= 0.1
            }
        case s.HasWater():
            channel.multiplier = 0.0
        default:
            channel.additive = this.mystery.At(p) * 5.0
        }
    })
}

func (this *Scent) EmanateBattle() {
    ForEachPoint(func(p Point) {
        channel := &this.channels[p.row][p.col]

        channel.multiplier *= 0.99

        s := this.terrain.At(p)

        switch {
        case s.HasHill():
            if s.IsEnemy() {
                channel.additive = 500.0
            } else {
                channel.multiplier = 0.0
            }
        case s.HasAnt():
            if s.IsEnemy() {
                switch {
                case this.holyGround.At(p) < 20:
                    channel.additive = 500.0 * (20.0 - float32(this.holyGround.At(p))) / 20.0
                default:
                    channel.additive = 5.0
                }
            } else {
                channel.multiplier *= 0.8
            }
        case s.HasWater():
            channel.multiplier = 0.0
        default:
            channel.additive = this.mystery.At(p) * 1.0
        }
    })
}
*/

func (this *Scent) Spread() {
    var newValue [MAX_ROWS][MAX_COLS]float32

    var p Point
    for p.row = 0; p.row < rows; p.row++ {
        for p.col = 0; p.col < cols; p.col++ {
            channel := &this.channels[p.row][p.col]

            v := *channel.pathways[0]
            v += *channel.pathways[1]
            v += *channel.pathways[2]
            v += *channel.pathways[3]
            v += *channel.pathways[4]

            newValue[p.row][p.col] = v * channel.multiplier + channel.additive * 10000000000000000000000000000.0
        }
    }

    this.value = newValue
}

func (this *Scent) Calculate() {
    this.BuildChannels()

    //this.EmanateForage()

    for i := 0; i < 100; i++ {
        this.Spread()
    }
}

/*
func (this *Scent) CalculateBattle() {
    this.BuildChannels()

    this.EmanateBattle()

    for i := 0; i < 100; i++ {
        this.Spread()
    }
}
*/

func (this *Scent) String() string {
    return GridToString(func(p Point) byte {
        square := this.terrain.At(p)
        switch {
        case square.HasFood():
            return '*'
        case square.HasLand():
            switch {
            case this.At(p) <     0: return '-'
            case this.At(p) <     1: return '0'
            case this.At(p) <     2: return '1'
            case this.At(p) <     4: return '2'
            case this.At(p) <     8: return '3'
            case this.At(p) <    16: return '4'
            case this.At(p) <    32: return '5'
            case this.At(p) <    64: return '6'
            case this.At(p) <   128: return '7'
            case this.At(p) <   256: return '8'
            case this.At(p) <   512: return '9'
            case this.At(p) <  1024: return 'a'
            case this.At(p) <  2048: return 'b'
            case this.At(p) <  4096: return 'c'
            case this.At(p) <  8192: return 'd'
            case this.At(p) < 16384: return 'e'
            case this.At(p) < 32768: return 'f'
            }
            return '+'
        case square.HasWater():
            return '_'
        }
        return ' '
    })
}

func (this *Scent) Csv() string {
    return GridToCsv(func(p Point) string {
        return fmt.Sprintf("%v", this.At(p))
    })
}
