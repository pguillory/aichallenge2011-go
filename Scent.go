package main

import "fmt"

type ScentChannel struct {
    tickValue float32
    northTickValue, eastTickValue, southTickValue, westTickValue *float32
    tockValue float32
    northTockValue, eastTockValue, southTockValue, westTockValue *float32
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
    //value [MAX_ROWS][MAX_COLS]float32
    channels [MAX_ROWS][MAX_COLS]ScentChannel
}

func NewForageScent(terrain *Terrain, distanceToEnemy *TravelDistance, distanceToFriendlyHill *TravelDistance, mystery *Mystery) *Scent {
    return NewScent(terrain, distanceToEnemy, distanceToFriendlyHill, mystery, ConfigureForageScentChannel)
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

func (this *Scent) ChannelAt(p Point) *ScentChannel {
    return &this.channels[p.row][p.col]
}

func (this *Scent) Calculate() {
    this.distanceToEnemy.Calculate()
    this.distanceToFriendlyHill.Calculate()
    this.mystery.Calculate()

    if this.turn == turn {
        return
    }
    startTime := now()

    ForEachPoint(func(p Point) {
        channel := this.ChannelAt(p)

        channel.northTickValue = &this.ChannelAt(p.Neighbor(NORTH)).tickValue
        channel.eastTickValue  = &this.ChannelAt(p.Neighbor(EAST)).tickValue
        channel.southTickValue = &this.ChannelAt(p.Neighbor(SOUTH)).tickValue
        channel.westTickValue  = &this.ChannelAt(p.Neighbor(WEST)).tickValue

        channel.northTockValue = &this.ChannelAt(p.Neighbor(NORTH)).tockValue
        channel.eastTockValue  = &this.ChannelAt(p.Neighbor(EAST)).tockValue
        channel.southTockValue = &this.ChannelAt(p.Neighbor(SOUTH)).tockValue
        channel.westTockValue  = &this.ChannelAt(p.Neighbor(WEST)).tockValue

        channel.multiplier = 1.0 / (5.0 - float32(this.adjacentWater.At(p)))
        channel.additive = 0.0
        channel.tickValue = 0
    })

    ForEachPoint(func(p Point) {
        this.configureChannel(this, p, this.ChannelAt(p))
    })

    for i := 0; i < 200; i++ {
        this.Spread()
    }

    this.time = now() - startTime
    this.turn = turn
}

func ConfigureForageScentChannel(this *Scent, p Point, channel *ScentChannel) {
    //channel.tickValue = 1e30
    channel.multiplier *= 0.9

    s := this.terrain.At(p)

    switch {
    case s.HasFood():
        channel.additive += 1e15
    case s.HasHill():
        if s.IsEnemy() {
            channel.additive += 1e20
        } else {
            channel.multiplier *= 0.0
        }
    case s.HasAnt():
        if s.IsEnemy() {
            if this.distanceToFriendlyHill.At(p) < 20 {
                channel.additive += 1e20 * float32(20 - this.distanceToFriendlyHill.At(p)) / 20.0
            } else if this.terrain.VisibleFriendliesAt(p) >= 2 {
                channel.additive += 1e14
            }
        } else {
            if this.terrain.VisibleEnemiesAt(p) > 0 {
                channel.multiplier *= 1e-3
            } else {
                channel.multiplier *= 0.0
            }
        }
    case s.HasWater():
        channel.multiplier *= 0.0
    default:
        //channel.additive = this.mystery.At(p) * 0.002
        //this.value[p.row][p.col] += this.mystery.At(p) * 0.002
        channel.additive += this.mystery.At(p) * 1e12

        //if MAX_TRAVEL_DISTANCE > this.distanceToEnemy.At(p) {
        //    channel.additive += 1e12 * (float32(MAX_TRAVEL_DISTANCE - this.distanceToEnemy.At(p)) / float32(MAX_TRAVEL_DISTANCE))
        //}
    }
}

func (this *Scent) Spread() {
    var p Point

    for p.row = 0; p.row < rows; p.row++ {
        for p.col = 0; p.col < cols; p.col++ {
            channel := &this.channels[p.row][p.col]

            channel.tockValue = channel.tickValue
            channel.tockValue += *channel.northTickValue
            channel.tockValue += *channel.eastTickValue
            channel.tockValue += *channel.southTickValue
            channel.tockValue += *channel.westTickValue
            channel.tockValue *= channel.multiplier
            channel.tockValue += channel.additive
        }
    }

    for p.row = 0; p.row < rows; p.row++ {
        for p.col = 0; p.col < cols; p.col++ {
            channel := &this.channels[p.row][p.col]

            channel.tickValue = channel.tockValue
            channel.tickValue += *channel.northTockValue
            channel.tickValue += *channel.eastTockValue
            channel.tickValue += *channel.southTockValue
            channel.tickValue += *channel.westTockValue
            channel.tickValue *= channel.multiplier
            //channel.tickValue += channel.additive
        }
    }
}

func (this *Scent) At(p Point) float32 {
    //return this.value[p.row][p.col] + float32(MAX_TRAVEL_DISTANCE - this.distanceToEnemy.At(p)) * 1e-30
    return this.channels[p.row][p.col].tickValue
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
            case this.At(p) <= 1e2: return '1'
            case this.At(p) <= 1e4: return '2'
            case this.At(p) <= 1e6: return '3'
            case this.At(p) <= 1e8: return '4'
            case this.At(p) <= 1e10: return '5'
            case this.At(p) <= 1e12: return '6'
            case this.At(p) <= 1e14: return '7'
            case this.At(p) <= 1e16: return '8'
            case this.At(p) <= 1e18: return '9'
            case this.At(p) <= 1e20: return 'a'
            case this.At(p) <= 1e22: return 'b'
            case this.At(p) <= 1e24: return 'c'
            case this.At(p) <= 1e26: return 'd'
            case this.At(p) <= 1e28: return 'e'
            case this.At(p) <= 1e30: return 'f'
            case this.At(p) <= 1e32: return 'g'
            case this.At(p) <= 1e34: return 'h'
            case this.At(p) <= 1e36: return 'i'
            case this.At(p) <= 1e36: return 'j'
            case this.At(p) <= 1e38: return 'k'
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
