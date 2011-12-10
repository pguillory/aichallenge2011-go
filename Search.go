package main

import "fmt"

type Search struct {
    time int64
    turn int
    terrain *Terrain
    mystery *Mystery
    potentialEnemy *PotentialEnemy
    values *DistanceGridGrid
}

func NewSearch(terrain *Terrain, mystery *Mystery, potentialEnemy *PotentialEnemy) *Search {
    this := new(Search)
    this.terrain = terrain
    this.mystery = mystery
    this.potentialEnemy = potentialEnemy

    this.Calculate()
    return this
}

/*
func SearchFromFriendlyAnts(terrain *Terrain, f func(Point)) {
    queue := new(PointQueue)
    covered := new(PointSet)

    ForEachPoint(func(p Point) {
        square := terrain.At(p)
        switch {
        case square.HasFriendlyAnt():
            queue.Push(p)
            covered.Include(p)
        case terrain.At(p).HasWater():
            covered.Include(p)
        }
    })

    queue.ForEach(func(p Point) {
        ForEachNeighbor(p, func(p2 Point) {
            if !covered.Includes(p2) {
                covered.Include(p2)
                queue.Push(p2)
                f(p2)
            }
        })
    })
}
*/

func (this *Search) Calculate() {
    if this.turn == turn {
        return
    }
    startTime := now()

    goals := NewGoalSet(this.terrain)

    goals.Add("enemy hill", 0, 100, 1, func(p Point) bool {
        return this.terrain.At(p).HasEnemyHill()
    })
    goals.Add("unexplored", 10, 0.01, 100, func(p Point) bool {
       return this.mystery.UnexploredAt(p)
    })
    goals.Add("food", 10, 1, 1, func(p Point) bool {
        return this.terrain.At(p).HasFood()
    })
    goals.Add("enemy", 15, 5, 1000, func(p Point) bool {
        return this.terrain.At(p).HasEnemyAnt()
    })
    goals.Add("potential enemy", 15, 0.01, 100, func(p Point) bool {
       return this.potentialEnemy.At(p)
    })

    fmt.Println("turn", turn)
    fmt.Println(goals)

    this.values = goals.Calculate()

/*
    potentialDestinations := Everywhere()
    engine := NewSearchEngine(this.terrain)

    friendlyAnts := Where(func(p Point) bool {
        return terrain.At(p).HasFriendlyAnt()
    })

    spread := func(origin Point, matches func(Point) bool) (area *PointSet, cardinality int) {
        if !matches(origin) {
            return
        }

        area = new(PointSet)
        cardinality = 0

        queue := new(PointQueue)
        queue.Push(origin)
        queue.ForEach(func(p Point) {
            ForEachPointWithinManhattanDistance(p, 1, func(p2 Point) {
                if potentialDestinations.Includes(p2) && matches(p2) {
                    potentialDestinations.Exclude(p2)
                    area.Include(p2)
                    cardinality += 1
                    queue.Push(p2)
                }
            })
        })

        return
    }

    NewSearchEngine(this.terrain).uSearch(friendlyAnts, 0, func(origin) {
        for _, objective := range objectives {
            if potentialDestinations.Includes(origin) {
                area, cardinality := spread(origin, objective.matches)
                
            }
        }
    })
*/


    this.time = now() - startTime
    this.turn = turn
}

func (this *Search) DistancesFor(p Point) *DistanceGrid {
    result := this.values.At(p)
    if result == nil {
        panic(fmt.Sprintf("No goal for ant at %v", p))
    }
    return result
}
