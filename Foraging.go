/*
TODO:
- track ants waiting to spawn
*/

package main

//import "fmt"

func ForageMoves(terrain *Terrain) *MoveSet {
    foods := make([]Point, 0, 100)
    currentAnts := new(PointSet)
    availableAnts := new(PointSet)
    futureAnts := new(PointSet)
    hills := new(PointSet)
    land := new(PointSet)

    /*distanceToFriendlyAnt :=*/ NewTravelDistance(func(p Point) Distance {
        if terrain.At(p).HasLand() {
            land.Include(p)
        }
        if terrain.At(p).HasHill() {
            hills.Include(p)
        }
        if terrain.At(p).HasFriendlyAnt() {
            currentAnts.Include(p)
            availableAnts.Include(p)
            futureAnts.Include(p)
            //fmt.Printf("ant at %v\n", p)
            return 0
        }
        return MAX_TRAVEL_DISTANCE
    }, func(p Point, distance Distance, dir Direction) bool {
        return terrain.At(p).HasLand()
    }, func(p Point, distance Distance, dir Direction) bool {
        if terrain.At(p).HasFood() {
            foods = append(foods, p)
        }
        return terrain.At(p).HasLand()
    }, 20)

    //fmt.Println(distanceToFriendlyAnt)

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

    MoveAnt := func(p Point, destination Point) {
        switch {
        case futureAnts.Includes(p):
            //fmt.Printf("Taking ant from %v\n", p)
            futureAnts.Exclude(p)
            destinations[p.row][p.col] = destination
            destinations[destination.row][destination.col] = destination
        case hills.Includes(p):
            for i := Distance(0); i < 100; i++ {
                if spawns[i] {
                    //fmt.Printf("Taking ant that will spawn in %v turns\n", i)
                    spawns[i] = false
                    break
                }
            }
        }
    }

    AntWillArriveAt := func(p Point, distance Distance) {
        futureAnts.Include(p)
        arrivals[p.row][p.col] = distance
    }

    AntWillSpawn := func(distance Distance) {
        for i := distance; i < 100; i++ {
            if !spawns[i] {
                spawns[i] = true
                break
            }
        }
    }

    for _, food := range foods {
        //fmt.Printf("\nFood at %v\n", food)

        bestTime := MAX_TRAVEL_DISTANCE
        var closest Point

        /*distanceToThisFood :=*/ NewTravelDistance(func(p Point) Distance {
            if p.Equals(food) {
                return 0
            }
            return MAX_TRAVEL_DISTANCE
        }, func(p Point, distance Distance, dir Direction) bool {
            return terrain.At(p).HasLand() && !terrain.At(p).HasFriendlyHill()
        }, func(p Point, distance Distance, dir Direction) bool {
            thisTime := distance + TimeUntilNextAntAvailableAt(p)

            //if thisTime < MAX_TRAVEL_DISTANCE && !p.Equals(closest.from) {
            //    //fmt.Printf("Ant at %v can be there in %v turns from %v\n", p, thisTime, distance)
            //}

            if thisTime < bestTime {
                //fmt.Printf("%v beats %v\n", thisTime, bestTime)
                bestTime = thisTime
                closest = p
            //} else if thisTime == bestTime && p.Equals(closest.from) {
            //    closest.dir |= dir.Backward()
            }

            return terrain.At(p).HasLand() && distance < bestTime
        }, 20)
        //fmt.Println(distanceToThisFood)

        MoveAnt(closest, food)
        AntWillArriveAt(food, bestTime)
        AntWillSpawn(bestTime)

        if availableAnts.Includes(closest) {
            availableAnts.Exclude(closest)
        }
    }

    //fmt.Println("moves")
    moves := new(MoveSet)
    currentAnts.ForEach(func(p Point) {
        if !availableAnts.Includes(p) {
            destination := destinations[p.row][p.col]
            nextDestination := destinations[destination.row][destination.col]

            //fmt.Printf("moving ant at %v to %v (then %v)\n", p, destination, nextDestination)

            dir := WhichWay(p, destination, nextDestination, land)
            moves.Include(Move{p, dir})
        }
    })

    //fmt.Println(land)

    return moves
}
