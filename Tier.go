package main

//import "fmt"

type Tier struct {
    time int64
    turn int
    terrain *Terrain
    mystery *Mystery
    potentialEnemy *PotentialEnemy
    distanceAt [MAX_ROWS][MAX_COLS]Distance
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

func (this *Tier) Calculate() {
    queue := new(PointQueue)

    ForEachPoint(func(p Point) {
        this.distanceAt[p.row][p.col] = func() Distance {
            switch {
            case this.terrain.At(p).HasEnemyHill():
                return 0
            // TODO: prioritize enemy ants near friendly hills
            //case potentialEnemy.At(p) && distanceToFriendlyHill.At(p) <= 10:
            //    return distanceToFriendlyHill.At(p)
            case this.mystery.At(p) >= 50:
                return 10
            case this.potentialEnemy.At(p):
                return 15
            case this.mystery.At(p) >= 2:
                return 30
            case this.mystery.At(p) >= 1:
                return 31
            }
            return MAX_TRAVEL_DISTANCE
        }()

        if this.distanceAt[p.row][p.col] < MAX_TRAVEL_DISTANCE {
            queue.Push(p)
        }
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
                    case this.terrain.At(p2).HasFriendlyAnt():
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

    //antQueue.ForEach(func(p Point) {
    //    tier := this.tierAt[p.row][p.col]
    //
    //    ForEachNeighbor(p, func(p2 Point) {
    //        switch {
    //        case this.tierAt[p2.row][p2.col] >= tier:
    //        case this.terrain.At(p2).HasAnt():
    //            this.tierAt[p2.row][p2.col] = tier
    //            antQueue.Push(p2)
    //        }
    //    })
    //})
}

func (this *Tier) ValueOf(move Move) (result Distance) {
    destination := move.Destination()
    result += this.distanceAt[move.from.row][move.from.col]
    result -= this.distanceAt[destination.row][destination.col]
    return
}

func (this *Tier) At(p Point) int {
    return this.tierAt[p.row][p.col]
}

func (this *Tier) String() string {
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
