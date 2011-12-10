/*
TODO: track ants waiting to spawn
TODO: bug, ant on hill next to food
TODO: re-order unclaimed food by distance to ant after each food claimed
*/

package main

/*
//import "fmt"

func ForageMoves(terrain *Terrain) *MoveSet {
    land := new(PointSet)
    hills := new(PointSet)
    unclaimedFood := new(PointSet)
    currentAnts := new(PointSet)
    availableAnts := new(PointSet)
    futureAnts := new(PointSet)

    ForEachPoint(func(p Point) {
        if terrain.At(p).HasLand() {
            land.Include(p)
        }
        if terrain.At(p).HasHill() {
            hills.Include(p)
        }
        if terrain.At(p).HasFood() {
            unclaimedFood.Include(p)
        }
        if terrain.At(p).HasFriendlyAnt() {
            currentAnts.Include(p)
            availableAnts.Include(p)
            futureAnts.Include(p)
        }
    })

    var spawns [100]bool
    var arrivals [MAX_ROWS][MAX_COLS]Distance
    var destinations [MAX_ROWS][MAX_COLS]Point

    TimeUntilNextAntAvailableAt := func(p Point) Distance {
        switch {
        case futureAnts.Includes(p):
            return arrivals[p.row][p.col]
        case hills.Includes(p):
            for i := Distance(0); i < 100; i++ {
                if spawns[i] {
                    return i
                }
            }
        }
        return MAX_TRAVEL_DISTANCE
    }

    for i := 0; unclaimedFood.Cardinality() > 0 && i < 20; i++ {
        //fmt.Printf("\n\nIteration %v\n", i)

        foodByDistance := make([]Point, 0, 100)

        maxDistanceToFood := Distance(20)

        NewTravelDistance(func(p Point) Distance {
            return TimeUntilNextAntAvailableAt(p)
        }, func(p Point, distance Distance, dir Direction) bool {
            return land.Includes(p) && distance < maxDistanceToFood
        }, func(p Point, distance Distance, dir Direction) bool {
            if distance <= maxDistanceToFood {
                if unclaimedFood.Includes(p) {
                    unclaimedFood.Exclude(p)
                    foodByDistance = append(foodByDistance, p)
                    if maxDistanceToFood > distance {
                        maxDistanceToFood = distance
                    }
                }
                return land.Includes(p)
            }
            return false
        }, maxDistanceToFood)

        //fmt.Printf("Food at distance %v: %v\n", maxDistanceToFood, foodByDistance)


        for _, food := range foodByDistance {
            //fmt.Printf("\nFood at %v\n", food)

            bestTime := MAX_TRAVEL_DISTANCE
            var closest Point

            NewTravelDistance(func(p Point) Distance {
                if p.Equals(food) {
                    return 0
                }
                return MAX_TRAVEL_DISTANCE
            }, func(p Point, distance Distance, dir Direction) bool {
                return land.Includes(p) && !hills.Includes(p)
            }, func(p Point, distance Distance, dir Direction) bool {
                thisTime := distance + TimeUntilNextAntAvailableAt(p)

                if thisTime < MAX_TRAVEL_DISTANCE && !p.Equals(closest) {
                    //fmt.Printf("Ant at %v can be there in %v turns\n", p, thisTime)
                }

                if thisTime < bestTime {
                    bestTime = thisTime
                    closest = p
                //} else if thisTime == bestTime && p.Equals(closest.from) {
                //    closest.dir |= dir.Backward()
                }

                return land.Includes(p) && distance < bestTime
            }, maxDistanceToFood)

            switch {
            case futureAnts.Includes(closest):
                //fmt.Printf("Taking ant from %v\n", closest)
                futureAnts.Exclude(closest)
                destinations[closest.row][closest.col] = food
                destinations[food.row][food.col] = food
            case hills.Includes(closest):
                for i := Distance(0); i < 100; i++ {
                    if spawns[i] {
                        //fmt.Printf("Taking ant that will spawn in %v turns\n", i)
                        spawns[i] = false
                        break
                    }
                }
            }

            futureAnts.Include(food)
            arrivals[food.row][food.col] = bestTime

            for i := bestTime; i < 100; i++ {
                if !spawns[i] {
                    spawns[i] = true
                    break
                }
            }

            if availableAnts.Includes(closest) {
                availableAnts.Exclude(closest)
            }
        }
    }

    moves := new(MoveSet)
    currentAnts.ForEach(func(p Point) {
        if !availableAnts.Includes(p) {
            destination := destinations[p.row][p.col]
            nextDestination := destinations[destination.row][destination.col]

            //fmt.Printf("Moving ant at %v to %v\n", p, destination)

            dir := WhichWay(p, destination, nextDestination, land)
            if terrain.At(p.Neighbor(dir)).HasFood() {
                if terrain.At(p).HasFriendlyHill() {
                    //fmt.Printf("On hill with food adjacent to %v, going %v or %v\n", dir, dir.Right(), dir.Left())
                    dir = STAY | dir.Right() | dir.Left()
                } else {
                    dir = STAY
                }
            }
            moves.Include(Move{p, dir})
        }
    })

    return moves
}

// spawned on hill next to food, should we move off or stay?
//
//.. .. .. a.
//A. A. A. 0.
//*. .. a. ..
//.. .. .. a.
//
//.. .. .. ..
//A. 0a 0. A.
//*. *. .a ..
//.. .. .. .a
*/
