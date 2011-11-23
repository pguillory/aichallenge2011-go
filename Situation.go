package main

const (
    SITUATION_WATER    = byte(1)
    SITUATION_FOOD     = byte(2)
    SITUATION_ANT      = byte(4)
    SITUATION_HILL     = byte(8)
    SITUATION_ENEMY    = byte(16)

    SITUATION_BITS_STORED = 5
    SITUATION_RADIUS2 = 19
    SITUATION_POINTS_WITHIN_RADIUS2 = 61
    SITUATION_SIZE = (SITUATION_POINTS_WITHIN_RADIUS2 + SITUATION_BITS_STORED - 1) * SITUATION_BITS_STORED / 32
)

type Situation [SITUATION_SIZE]uint32

func VerifySituationSize() {
    i := 0
    ForEachPointWithinRadius2(Point{0, 0}, SITUATION_RADIUS2, func(p2 Point) {
        i++
    })
    if i != SITUATION_POINTS_WITHIN_RADIUS2 {
        panic("VerifySituationSize() failed")
    }
}

func NewSituation(terrain *Terrain, p Point) (this *Situation, friendlyNearby bool) {
    this = new(Situation)

    player := terrain.At(p).owner
    i := uint(0)

    ForEachPointWithinRadius2(p, SITUATION_RADIUS2, func(p2 Point) {
        square := terrain.At(p2.Normalize())

        switch {
        case square.HasWater():
            this[(i + uint(0)) / 32] |= 1 << ((i + uint(0)) % 32)
        case square.HasFood():
            this[(i + uint(1)) / 32] |= 1 << ((i + uint(1)) % 32)
        default:
            if square.HasAnt() {
                this[(i + uint(2)) / 32] |= 1 << ((i + uint(2)) % 32)
                if square.owner != player {
                    this[(i + uint(4)) / 32] |= 1 << ((i + uint(4)) % 32)
                }
                if square.owner == 0 {
                    friendlyNearby = true
                }
            }
            if square.HasHill() {
                this[(i + uint(3)) / 32] |= 1 << ((i + uint(3)) % 32)
                if square.owner != player {
                    this[(i + uint(4)) / 32] |= 1 << ((i + uint(4)) % 32)
                }
            }
        }

        i += 5
    })

    return
}

func (this *Situation) Matches(situation *Situation) bool {
    for i := 0; i < SITUATION_SIZE; i++ {
        if this[i] != situation[i] {
            return false
        }
    }
    return true
}

func (this *Situation) String() string {
    i := uint(0)

    return StringFromRadius2(Point{0, 0}, SITUATION_RADIUS2, func(Point) byte {
        hasWater := this[(i + uint(0)) / 32] & (1 << ((i + uint(0)) % 32)) > 0
        hasFood  := this[(i + uint(1)) / 32] & (1 << ((i + uint(1)) % 32)) > 0
        hasAnt   := this[(i + uint(2)) / 32] & (1 << ((i + uint(2)) % 32)) > 0
        hasHill  := this[(i + uint(3)) / 32] & (1 << ((i + uint(3)) % 32)) > 0
        isEnemy  := this[(i + uint(4)) / 32] & (1 << ((i + uint(4)) % 32)) > 0
        i += 5

        switch {
        case hasWater:
            return '%'
        case hasFood:
            return '*'
        case hasAnt && hasHill:
            if isEnemy {
                return 'B'
            } else {
                return 'A'
            }
        case hasAnt:
            if isEnemy {
                return 'b'
            } else {
                return 'a'
            }
        case hasHill:
            if isEnemy {
                return '1'
            } else {
                return '0'
            }
        }

        return '.'
    })
}
