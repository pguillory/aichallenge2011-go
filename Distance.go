package main

//import "fmt"

type Distance uint16

const MAX_TRAVEL_DISTANCE = MAX_ROWS * MAX_COLS

type TravelDistance struct {
    time int64
    turn int
    value [MAX_ROWS][MAX_COLS]Distance
    initialValue func(p Point) Distance
    isEnterable, isExitable func(p Point) bool
}

func DistanceToFriendlyHill(terrain *Terrain) *TravelDistance {
    return NewTravelDistance(func(p Point) Distance {
        if terrain.At(p).HasFriendlyHill() {
            return 0
        }
        return MAX_TRAVEL_DISTANCE
    }, func(p Point) bool {
        return terrain.At(p).HasLand()
    }, func(p Point) bool {
        return terrain.At(p).HasLand()
    })
}

func DistanceToFood(terrain *Terrain) *TravelDistance {
    distance := NewTravelDistance(func(p Point) Distance {
        square := terrain.At(p)
        switch {
        case square.HasFood():
            return 0
        }
        return MAX_TRAVEL_DISTANCE
    }, func(p Point) bool {
        square := terrain.At(p)
        return !square.HasWater()
    }, func(p Point) bool {
        square := terrain.At(p)
        return !square.HasWater()
    })

    return distance
}

func DistanceToTrouble(terrain *Terrain, mystery *Mystery, potentialEnemy *PotentialEnemy, scrum *Scrum) *TravelDistance {
    distance := NewTravelDistance(func(p Point) Distance {
        switch {
        case terrain.At(p).HasEnemyHill():
            return 0
        case mystery.At(p) >= 50:
            return 10
        case potentialEnemy.At(p):
            return 15
        case mystery.At(p) >= 2:
            return 30
        case mystery.At(p) >= 1:
            return 31
        case scrum.At(p):
            return 200
        }
        return MAX_TRAVEL_DISTANCE
    }, func(p Point) bool {
        square := terrain.At(p)
        return !square.HasWater() && !square.HasFriendlyAnt()
    }, func(p Point) bool {
        square := terrain.At(p)
        return !square.HasWater()
    })

    return distance
}

func DistanceToDoom(terrain *Terrain, mystery *Mystery, potentialEnemy *PotentialEnemy, scrum *Scrum) *TravelDistance {
    distance := NewTravelDistance(func(p Point) Distance {
        switch {
        case terrain.At(p).HasEnemyHill():
            return 0
        case mystery.At(p) >= 50:
            return 10
        case potentialEnemy.At(p):
            return 15
        case scrum.At(p):
            return 200
        }
        return MAX_TRAVEL_DISTANCE
    }, func(p Point) bool {
        square := terrain.At(p)
        return !square.HasWater()
    }, func(p Point) bool {
        square := terrain.At(p)
        return !square.HasWater()
    })

    return distance
}

func AssignForagers(terrain *Terrain) *PointSet {
    land := new(PointSet)
    food := new(PointSet)
    availableAnts := new(PointSet)
    foragers := new(PointSet)

    ForEachPoint(func(p Point) {
        if terrain.At(p).HasLand() {
            land.Include(p)
        }
        if terrain.At(p).HasFood() {
            food.Include(p)
        }
        if terrain.At(p).HasFriendlyAnt() {
            availableAnts.Include(p)
        }
    })

    food.ForEach(func(p Point) {
        move, found := FindNearestMoveTo(p, 20, land, availableAnts)
        if found {
            availableAnts.Exclude(move.from)
            foragers.Include(move.from)
        }
    })

    return foragers
}

func FindNearestMoveTo(destination Point, within Distance, passablePoints, fromPoints *PointSet) (result Move, found bool) {
    var values [MAX_ROWS][MAX_COLS]Distance
    var queue PointQueue

    values[destination.row][destination.col] = MAX_TRAVEL_DISTANCE

    queue.Push(destination)

    queue.ForEach(func(p Point) {
        value := values[p.row][p.col]
        if MAX_TRAVEL_DISTANCE - value <= within {
            ForEachDirection(func(dir Direction) {
                p2 := p.Neighbor(dir)

                if values[p2.row][p2.col] < value - 1 {
                    switch {
                    case fromPoints.Includes(p2):
                        result = Move{p2, dir.Backward()}
                        found = true
                        queue.Clear()
                    case passablePoints.Includes(p2):
                        values[p2.row][p2.col] = value - 1
                        queue.Push(p2)
                    }
                }
            })
        }
    })

    return
}

func NewTravelDistance(initialValue func(p Point) Distance, isEnterable, isExitable func(p Point) bool) *TravelDistance {
    this := new(TravelDistance)
    this.initialValue = initialValue
    this.isExitable = isExitable
    this.isEnterable = isEnterable

    this.Calculate()
    return this
}

func (this *TravelDistance) Calculate() {
    if this.turn == turn {
        return
    }
    startTime := now()

    //queue := new(PointQueue)
    var queue PointQueue

    //fmt.Printf("%v in queue\n", queue.Size())

    ForEachPoint(func(p Point) {
        value := this.initialValue(p)
        this.value[p.row][p.col] = value
        if value < MAX_TRAVEL_DISTANCE {
            queue.Push(p)
        }
    })

    //fmt.Printf("%v in queue\n", queue.Size())

    queue.ForEach(func(p Point) {
        v := this.value[p.row][p.col]
        //fmt.Printf("got %v: %v\n", p, v)

        ForEachNeighbor(p, func(p2 Point) {
            if this.value[p2.row][p2.col] > v + 1 {
                if this.isExitable(p2) {
                    //fmt.Printf("spreading to %v\n", p2)
                    this.value[p2.row][p2.col] = v + 1
                    if this.isEnterable(p2) {
                        queue.Push(p2)
                    }
                } else {
                    //fmt.Printf("%v not exitable\n", p2)
                }
            } else {
                //fmt.Printf("%v too low (%v)\n", p2, this.value[p2.row][p2.col])
            }
        })
    })

    this.time = now() - startTime
    this.turn = turn
}

func (this *TravelDistance) At(p Point) int {
    return int(this.value[p.row][p.col])
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
