package main

//import "fmt"

type Distance uint16

const MAX_TRAVEL_DISTANCE = Distance(MAX_ROWS * MAX_COLS)

type TravelDistance struct {
    time int64
    turn int
    value [MAX_ROWS][MAX_COLS]Distance
    initialValue func(Point) Distance
    isEnterable, isExitable func(Point, Distance, Direction) bool
    max Distance
}

func DistanceToFriendlyHill(terrain *Terrain) *TravelDistance {
    return NewTravelDistance(func(p Point) Distance {
        if terrain.At(p).HasFriendlyHill() {
            return 0
        }
        return MAX_TRAVEL_DISTANCE
    }, func(p Point, distance Distance, dir Direction) bool {
        return terrain.At(p).HasLand()
    }, func(p Point, distance Distance, dir Direction) bool {
        return terrain.At(p).HasLand()
    }, 20)
}

func DistanceToFood(terrain *Terrain) *TravelDistance {
    distance := NewTravelDistance(func(p Point) Distance {
        square := terrain.At(p)
        switch {
        case terrain.At(p).HasEnemyHill():
            return 0
        case square.HasFood():
            return 2
        }
        return MAX_TRAVEL_DISTANCE
    }, func(p Point, distance Distance, dir Direction) bool {
        square := terrain.At(p)
        return square.HasLand()
    }, func(p Point, distance Distance, dir Direction) bool {
        square := terrain.At(p)
        return square.HasLand()
    }, 22)

    return distance
}

//func DistanceToFewerFriendliesThan(max byte, terrain *Terrain) *TravelDistance {
//    return NewTravelDistance(func(p Point) Distance {
//        if terrain.VisibleFriendliesAt(p) < max {
//            return 0
//        }
//        return MAX_TRAVEL_DISTANCE
//    }, func(p Point) bool {
//        square := terrain.At(p)
//        return !square.HasWater() && !square.HasFriendlyAnt()
//    }, func(p Point) bool {
//        square := terrain.At(p)
//        return !square.HasWater() && !square.HasFriendlyHill()
//    }, 50)
//}

func DistanceToTrouble(terrain *Terrain, mystery *Mystery, potentialEnemy *PotentialEnemy) *TravelDistance {
    //var distanceToFewerFriendliesThan []*TravelDistance
    //distanceToFewerFriendliesThan[1].At(p)

    distance := NewTravelDistance(func(p Point) Distance {
        switch {
        // TODO: prioritize enemy ants near friendly hills
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
        }
        return MAX_TRAVEL_DISTANCE
    }, func(p Point, distance Distance, dir Direction) bool {
        square := terrain.At(p)
        return !square.HasWater() && !square.HasFriendlyAnt()
    }, func(p Point, distance Distance, dir Direction) bool {
        square := terrain.At(p)
        return !square.HasWater() && !square.HasFriendlyHill()
    }, MAX_TRAVEL_DISTANCE)

    return distance
}

func DistanceToDoom(terrain *Terrain, mystery *Mystery, potentialEnemy *PotentialEnemy) *TravelDistance {
    distance := NewTravelDistance(func(p Point) Distance {
        switch {
        case terrain.At(p).HasEnemyHill():
            return 0
        case potentialEnemy.At(p):
            return 10
        }
        return MAX_TRAVEL_DISTANCE
    }, func(p Point, distance Distance, dir Direction) bool {
        square := terrain.At(p)
        // TODO && !army.IsBerzerker2()
        return !square.HasWater()
    }, func(p Point, distance Distance, dir Direction) bool {
        square := terrain.At(p)
        return !square.HasWater()
    }, MAX_TRAVEL_DISTANCE)

    return distance
}

func DistanceToSoldier(terrain *Terrain, army *Army) *TravelDistance {
    // TODO
    //not tested!

    distance := NewTravelDistance(func(p Point) Distance {
        switch {
        case terrain.At(p).HasFriendlyAnt() && army.IsSoldierAt(p):
            return 0
        }
        return MAX_TRAVEL_DISTANCE
    }, func(p Point, distance Distance, dir Direction) bool {
        square := terrain.At(p)
        return !square.HasWater() && !square.HasAnt()
    }, func(p Point, distance Distance, dir Direction) bool {
        square := terrain.At(p)
        return !square.HasWater()
    }, 10)

    return distance
}

/*
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
*/

func WhichWay(origin, destination, nextDestination Point, passablePoints *PointSet) (result Direction) {
    checked := new(PointSet)
    queue := new(PointQueue)

    checked.Include(destination)
    queue.Push(destination)

    found := false

    queue.ForEach(func(p Point) {
        ForEachDirection(func(dir Direction) {
            p2 := p.Neighbor(dir)

            switch {
            case found:
            case checked.Includes(p2):
            case origin.Equals(p2):
                result = dir.Backward()
                found = true

                queue.Clear()
            case passablePoints.Includes(p2):
                checked.Include(p2)
                queue.Push(p2)
            }
        })
    })

    return
}

func NewTravelDistance(initialValue func(Point) Distance, isEnterable, isExitable func(Point, Distance, Direction) bool, max Distance) *TravelDistance {
    this := new(TravelDistance)
    this.initialValue = initialValue
    this.isExitable = isExitable
    this.isEnterable = isEnterable
    this.max = max

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
        if value < this.max {
            queue.Push(p)
        }
    })

    //fmt.Printf("%v in queue\n", queue.Size())

    queue.ForEach(func(p Point) {
        distance := this.value[p.row][p.col] + 1
        //fmt.Printf("got %v: %v\n", p, v)

        ForEachDirection(func(dir Direction) {
            p2 := p.Neighbor(dir)
            distance2 := this.value[p2.row][p2.col]
            if distance2 > distance && distance <= this.max {
                if this.isExitable(p2, distance, dir) {
                    this.value[p2.row][p2.col] = distance
                    if this.isEnterable(p2, distance, dir) {
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
