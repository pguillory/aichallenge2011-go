package main

//import "fmt"

/*
type DistanceGrid [MAX_ROWS][MAX_COLS]Distance

func (this *DistanceGrid) At(p Point) Distance {
    return this[p.row][p.col]
}

func (this *DistanceGrid) SetAt(p Point, value Distance) {
    this[p.row][p.col] = value
}

func (this *DistanceGrid) String() string {
    return GridToString(func(p Point) byte {
        v := this.At(p)
        if v == MAX_TRAVEL_DISTANCE {
            return '!'
        }
        if v < 10 {
            return '0' + byte(v)
        }
        v -= 10
        v /= 10
        if v < 26 {
            return 'a' + byte(v)
        }
        return '+'
    })
}

type Tier struct {
    time int64
    turn int
    terrain *Terrain
    mystery *Mystery
    potentialEnemy *PotentialEnemy
    distanceAt DistanceGrid
    tierAt [MAX_ROWS][MAX_COLS]int
}

func NewTier(terrain *Terrain, mystery *Mystery, potentialEnemy *PotentialEnemy) *Tier {
    this := new(Tier)
    this.terrain = terrain
    this.mystery = mystery
    this.potentialEnemy = potentialEnemy

    this.Calculate()
    return this
}

func (this *Tier) DetermineThreats() (*PointSet, *DistanceGrid) {
    queue := new(PointQueue)
    threats := new(PointSet)
    distanceToFriendlyHill := new(DistanceGrid)
    defenders := new(PointSet)
    enemyAntQueue := new(PointQueue)

    ForEachPoint(func(p Point) {
        distanceToFriendlyHill[p.row][p.col] = func() Distance {
            switch {
            case this.terrain.At(p).HasFriendlyHill():
                queue.Push(p)
                return 0
            }
            return MAX_TRAVEL_DISTANCE
        }()
    })

    queue.ForEach(func(p Point) {
        distance := distanceToFriendlyHill[p.row][p.col]

        ForEachNeighbor(p, func(p2 Point) {
            switch {
            case p2.Equals(p):
            case distanceToFriendlyHill.At(p2) <= distance + 1:
            case this.terrain.At(p2).HasWater():
            default:
                distanceToFriendlyHill.SetAt(p2, distance + 1)
                queue.Push(p2)

                switch {
                case this.terrain.At(p2).HasFriendlyAnt():
                    defenders.Include(p2)
                case this.terrain.At(p2).HasEnemyAnt():
                    enemyAntQueue.Push(p2)
                }
            }
        })
    })

    enemyAntQueue.ForEachReverse(func(attacker Point) {
        distance := distanceToFriendlyHill[attacker.row][attacker.col]
        alive := true

        for p := attacker; alive && distance > 0; distance-- {
            ForEachNeighbor(p, func(next Point) {
                if alive {
                    if distanceToFriendlyHill.At(next) < distance {
                        p = next
                        ForEachPointWithinRadius2(next, attackradius2, func(defender Point) {
                            if alive && defenders.Includes(defender) {
                                //fmt.Printf("attacker %v killed %v\n", attacker, defender)
                                alive = false
                                defenders.Exclude(defender)
                            }
                        })
                    }
                }
            })
        }

        if alive {
            //fmt.Printf("attacker %v made it to hill %v\n", attacker, p)
            threats.Include(attacker)
        }
    })

    return threats, distanceToFriendlyHill
}

func (this *Tier) Calculate() {
    threats, distanceToFriendlyHill := this.DetermineThreats()
    //fmt.Println(threats)

    queue := new(PointQueue)
    ForEachPoint(func(p Point) {
        this.distanceAt[p.row][p.col] = func() Distance {
            switch {
            case this.terrain.At(p).HasEnemyHill():
                return 0
            // TODO: prioritize enemy ants near friendly hills
            //case potentialEnemy.At(p) && distanceToFriendlyHill.At(p) <= 10:
            //    return distanceToFriendlyHill.At(p)
            case threats.Includes(p):
                return distanceToFriendlyHill[p.row][p.col]
            //case this.mystery.At(p) >= 50:
            //    return 10
            ////case this.potentialEnemy.At(p):
            ////    return 15
            //case this.mystery.At(p) >= 2:
            //    return 30
            //case this.mystery.At(p) >= 1:
            //    return 31

            case this.mystery.At(p) >= 100:
                return 10
            //case this.potentialEnemy.At(p):
            //   return 15
            case this.mystery.At(p) > 0:
                if this.mystery.At(p) < 0 || this.mystery.At(p) > 100 {
                    panic("mystery value out of range")
                }
                return 45 - Distance(this.mystery.At(p) / 4)
            }
            return MAX_TRAVEL_DISTANCE
        }()

        if this.distanceAt[p.row][p.col] < MAX_TRAVEL_DISTANCE {
            queue.Push(p)
        }

        this.tierAt[p.row][p.col] = 0
    })
    
    //fmt.Printf("%v in queue\n", queue.Size())

    maxDistance := Distance(0)
    antQueue := new(PointQueue)

    for tier := 1; !queue.Empty(); tier++ {
        //fmt.Println("Tier", tier)

        maxDistance += 50
        nextQueue := new(PointQueue)

        //touched := new(PointSet)

        queue.ForEach(func(p Point) {
            //touched.Include(p)
            //fmt.Println(touched)

            distance := this.distanceAt[p.row][p.col]

            //fmt.Println(p, distance)

            ForEachNeighbor(p, func(p2 Point) {
                switch {
                case p2.Equals(p):
                case this.distanceAt[p2.row][p2.col] <= distance + 1:
                    //fmt.Printf("Already touched %v at distance %v\n", p2, this.distanceAt[p2.row][p2.col])
                case 0 < this.tierAt[p2.row][p2.col] && this.tierAt[p2.row][p2.col] <= this.tierAt[p.row][p.col]:
                    //fmt.Printf("Same tier %v\n", p2)
                case this.terrain.At(p2).HasWater():
                    //fmt.Printf("Can't cross %v\n", p2)
                default:
                    this.tierAt[p2.row][p2.col] = tier
                    this.distanceAt[p2.row][p2.col] = distance + 1

                    switch {
                    case distance + 1 >= maxDistance:
                        //fmt.Printf("Reached max distance %v at %v\n", distance, p2)
                        nextQueue.Push(p2)
                    case this.terrain.At(p2).HasFriendlyAnt() && distance < 25:
                        //fmt.Printf("Ant at %v\n", p2)
                        nextQueue.Push(p2)
                        antQueue.Push(p2)
                    default:
                        //fmt.Printf("Spreading to %v\n", p2)
                        queue.Push(p2)
                    }
                }
            })
        })

        queue = nextQueue
        //fmt.Println(this)
    }

    antQueue.ForEach(func(p Point) {
       tier := this.tierAt[p.row][p.col]
    
       ForEachPointWithinManhattanDistance(p, 1, func(p2 Point) {
           switch {
           case this.tierAt[p2.row][p2.col] >= tier:
           case this.terrain.At(p2).HasAnt():
               this.tierAt[p2.row][p2.col] = tier
               antQueue.Push(p2)
           }
       })
    })
}

func (this *Tier) DistanceAt(p Point) Distance {
    return this.distanceAt[p.row][p.col]
}

func (this *Tier) At(p Point) int {
    return this.tierAt[p.row][p.col]
}

func (this *Tier) DistanceString() string {
    return GridToString(func(p Point) byte {
        v := this.DistanceAt(p) / 2
        switch {
        case v < 10:
            return '0' + byte(v)
        case v < 36:
            return 'a' + byte(v - 10)
        case v == MAX_TRAVEL_DISTANCE:
            return '!'
        }
        return '+'
    })
}

func (this *Tier) String() string {
    return GridToString(func(p Point) byte {
        v := this.At(p)
        switch {
        case v < 10:
            return '0' + byte(v)
        case v < 36:
            return 'a' + byte(v - 10)
        }
        return '+'
    })
}
*/
