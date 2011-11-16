package main

type Direction byte

const (
    NORTH = Direction(1)
    EAST  = Direction(2)
    SOUTH = Direction(4)
    WEST  = Direction(8)
    STAY  = Direction(16)
)

func (this Direction) Left() Direction {
    switch (this) {
    case NORTH:
        return WEST
    case SOUTH:
        return EAST
    case EAST:
        return NORTH
    case WEST:
        return SOUTH
    }
    return this
}

func (this Direction) Right() Direction {
    switch (this) {
    case NORTH:
        return EAST
    case SOUTH:
        return WEST
    case EAST:
        return SOUTH
    case WEST:
        return NORTH
    }
    return this
}

func (this Direction) Backward() Direction {
    switch (this) {
    case NORTH:
        return SOUTH
    case SOUTH:
        return NORTH
    case EAST:
        return WEST
    case WEST:
        return EAST
    }
    return this
}

func (this Direction) Char() byte {
    switch (this) {
    case 0:
        return '-'
    case NORTH:
        return 'N'
    case SOUTH:
        return 'S'
    case EAST:
        return 'E'
    case WEST:
        return 'W'
    case STAY:
        return 'X'
    }
    return '+'
}

func (this Direction) String() string {
    switch (this) {
    case 0:
        return "-"
    case NORTH:
        return "N"
    case SOUTH:
        return "S"
    case EAST:
        return "E"
    case WEST:
        return "W"
    case STAY:
        return "X"
    }
    return "+"
}

func ForEachDirection(f func(Direction)) {
    var dir Direction
    for dir = 1; dir <= STAY; dir *= 2 {
        f(dir)
    }
}
