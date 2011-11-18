package main

import "bytes"

type Direction byte

const (
    NORTH = Direction(1)
    EAST  = Direction(2)
    SOUTH = Direction(4)
    WEST  = Direction(8)
    STAY  = Direction(16)

    DIRECTIONS = 5
)

func (this Direction) Left() Direction {
    switch this {
    case NORTH:
        return WEST
    case EAST:
        return NORTH
    case SOUTH:
        return EAST
    case WEST:
        return SOUTH
    }
    return this
}

func (this Direction) Right() Direction {
    switch this {
    case NORTH:
        return EAST
    case EAST:
        return SOUTH
    case SOUTH:
        return WEST
    case WEST:
        return NORTH
    }
    return this
}

func (this Direction) Backward() Direction {
    switch this {
    case NORTH:
        return SOUTH
    case EAST:
        return WEST
    case SOUTH:
        return NORTH
    case WEST:
        return EAST
    }
    return this
}

func (this Direction) IsMultiple() bool {
    switch this {
    case 0:
    case NORTH:
    case EAST:
    case SOUTH:
    case WEST:
    case STAY:
    default:
        return true
    }
    return false
}

func (this Direction) Includes(dir Direction) bool {
    return this & dir > 0
}

func (this Direction) Minus(dir Direction) Direction {
    return this & ^dir
}

func (this Direction) Char() byte {
    switch this {
    case 0:
        return '-'
    case NORTH:
        return 'N'
    case EAST:
        return 'E'
    case SOUTH:
        return 'S'
    case WEST:
        return 'W'
    case STAY:
        return 'X'
    }
    return '+'
}

func (this Direction) String() string {
    b := new(bytes.Buffer)

    ForEachDirection(func(dir Direction) {
        if this.Includes(dir) {
            b.WriteByte(dir.Char())
        }
    })

    return b.String()
}

func ForEachDirection(f func(Direction)) {
    f(NORTH)
    f(EAST)
    f(SOUTH)
    f(WEST)
    f(STAY)
}
