package main

import "fmt"
import "sort"

type Goal struct {
    name string
    capacity float64
    area *PointSet
    distances *DistanceGrid
}

func (this *Goal) String() string {
    return this.name
}



type AntGoalTuple struct {
    ant Point
    distance Distance
    goal *Goal
}

type AntGoalMatrix struct {
    tuples []AntGoalTuple
}

func NewAntGoalMatrix() *AntGoalMatrix {
    this := new(AntGoalMatrix)
    this.tuples = make([]AntGoalTuple, 0)
    return this
}

func (this *AntGoalMatrix) Add(ant Point, distance Distance, goal *Goal) {
    this.tuples = append(this.tuples, AntGoalTuple{ant, distance, goal})
}

func (this *AntGoalMatrix) ForEachByDistance(f func(Point, *Goal)) {
    sort.Sort(this)

    for _, tuple := range this.tuples {
        f(tuple.ant, tuple.goal)
    }
}

func (this *AntGoalMatrix) Len() int {
    return len(this.tuples)
}

func (this *AntGoalMatrix) Less(i, j int) bool {
    return this.tuples[i].distance < this.tuples[j].distance
}

func (this *AntGoalMatrix) Swap(i, j int) {
    this.tuples[i], this.tuples[j] = this.tuples[j], this.tuples[i]
}


type GoalSet struct {
    engine1, engine2 *SearchEngine
    claimed, friendlyAnts *PointSet
    friendlyAntsAvailable int
    goals []*Goal
    matrix []AntGoalMatrix
}

func NewGoalSet(terrain *Terrain) *GoalSet {
    this := new(GoalSet)
    this.engine1 = NewSearchEngine(terrain)
    this.engine2 = NewSearchEngine(terrain)
    this.claimed = new(PointSet)
    this.friendlyAnts = Where(func(p Point) bool {
        return terrain.At(p).HasFriendlyAnt()
    })
    this.friendlyAntsAvailable = this.friendlyAnts.Cardinality()
    this.goals = make([]*Goal, 0, 100)
    return this
}

func (this *GoalSet) Add(name string, priority Distance, capacityRate float64, maxArea int, matches func(Point) bool) {
    //fmt.Printf("Checking %v for %v\n", origin, name)

    matchedGoalCount := 0

    this.engine1.SearchFrom(this.friendlyAnts, func(origin Point, distanceToFriendlyAnt Distance) (spread bool) {
        spread = true

        if !this.claimed.Includes(origin) && matches(origin) {
            this.claimed.Include(origin)
            matchedGoalCount += 1
            //fmt.Printf("Goal: %v around %v... ", name, origin)

            goal := new(Goal)
            goal.name = fmt.Sprintf("%v-%v", origin, name)
            goal.capacity = capacityRate

            goal.area = new(PointSet)
            goal.area.Include(origin)
            areaSize := 1

            for areaSize < maxArea {
                //fmt.Println("iterating")

                this.engine2.SearchFrom(goal.area, func(p Point, distanceToOrigin Distance) (spread bool) {
                    if areaSize < maxArea && !this.claimed.Includes(p) && matches(p) {
                        spread = true
                        this.claimed.Include(p)
                        goal.area.Include(p)
                        goal.capacity += capacityRate
                        areaSize += 1
                    }

                    return
                })

                if areaSize < maxArea {
                    //fmt.Println("iterating2")
                    found := false

                    this.engine2.SearchFrom(goal.area, func(p Point, distanceToOrigin Distance) (spread bool) {
                        if areaSize < maxArea && distanceToOrigin <= 3 {
                            spread = true

                            if !this.claimed.Includes(p) && matches(p) {
                                this.claimed.Include(p)
                                goal.area.Include(p)
                                goal.capacity += capacityRate
                                areaSize += 1
                                found = true
                            }
                        }

                        return
                    })

                    if !found {
                        break
                    }
                }
            }

            goal.distances = this.engine2.DistanceTo(goal.area)
            //goal.capacity = int(area * capacityRate)

            if goal.capacity >= 1 {
                this.goals = append(this.goals, goal)
                fmt.Printf("Found goal %v (x %v)\n", goal, areaSize)
            } else {
                fmt.Printf("Ignoring goal %v, capacity %v\n", goal, goal.capacity)
            }
        }

        return
    })

    fmt.Printf("Matched %v %v\n", matchedGoalCount, name)
}

func (this *GoalSet) Calculate() *DistanceGridGrid {
    matrix := NewAntGoalMatrix()

    this.friendlyAnts.ForEach(func(ant Point) {
        for _, goal := range this.goals {
            matrix.Add(ant, goal.distances.At(ant), goal)
        }
    })

    result := new(DistanceGridGrid)

    matrix.ForEachByDistance(func(ant Point, goal *Goal) {
        if this.friendlyAnts.Includes(ant) && goal.capacity >= 1 {
            this.friendlyAnts.Exclude(ant)
            goal.capacity -= 1
            result.SetAt(ant, goal.distances)
        }
    })

    return result
}

func (this *GoalSet) String() string {
    var cohorts [MAX_ROWS][MAX_COLS]int

    for i, goal := range this.goals {
        goal.area.ForEach(func(p Point) {
            if cohorts[p.row][p.col] > 0 {
                cohorts[p.row][p.col] = -1
            } else {
                cohorts[p.row][p.col] = i + 1
            }
        })
    }

    return GridToString(func(p Point) byte {
        v := cohorts[p.row][p.col]
        if v == -1 {
            return '+'
        }
        if v == 0 {
            return '.'
        }
        v %= 36
        if v < 10 {
            return '0' + byte(v)
        }
        v -= 10
        return 'a' + byte(v)
    })
}
