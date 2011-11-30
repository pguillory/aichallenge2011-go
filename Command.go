/*
TODO: break into two modules: tactical vs scent-based
TODO: rush a hill if someone else is bout to cap it
*/

package main

//import "fmt"

type Command struct {
    time int64
    turn int
    terrain *Terrain
    army *Army
    predictions *Predictions
    distanceToFood, distanceToTrouble, distanceToDoom *TravelDistance
    reinforcement *Reinforcement
    moves, enemyMoves *MoveSet
    enemies *PointSet
    enemyDestinations *PointSet
    friendlyFocus, maxFriendlyFocus, maxFriendlyFocus_STAY *Focus
    //dirs [MAX_ROWS][MAX_COLS]Direction
    //len int
}

func NewCommand(terrain *Terrain, army *Army, predictions *Predictions, distanceToFood, distanceToTrouble, distanceToDoom *TravelDistance, reinforcement *Reinforcement) *Command {
    this := new(Command)
    this.terrain = terrain
    this.army = army
    this.predictions = predictions
    this.distanceToFood = distanceToFood
    this.distanceToTrouble = distanceToTrouble
    this.distanceToDoom = distanceToDoom
    this.reinforcement = reinforcement

    this.Calculate()
    return this
}

func (this *Command) At(p Point) Direction {
    return this.moves.At(p)
}

func (this *Command) Reset() {
    this.moves = new(MoveSet)
    this.enemyMoves = new(MoveSet)
    this.enemies = new(PointSet)

    distanceToSoldier := DistanceToSoldier(this.terrain, this.army)

    ForEachPoint(func(p Point) {
        s := this.terrain.At(p)
        if s.HasFriendlyAnt() {
            this.moves.IncludeAllFrom(p)
        } else if s.HasEnemyAnt() {
            // TODO: should this include berzerkers also?
            if distanceToSoldier.At(p) < 6 {
                this.enemyMoves.IncludeAllFrom(p)
            } else {
                this.enemyMoves.Include(Move{p, this.predictions.At(p)})
            }
            this.enemies.Include(p)
        }
    })

    ForEachPoint(func(p Point) {
        s := this.terrain.At(p)
        switch {
        case s.HasWater() || s.HasFood():
            this.moves.ExcludeMovesTo(p)
            this.enemyMoves.ExcludeMovesTo(p)
        case s.HasFriendlyHill():
            ForEachPointWithinRadius2(p, 5, func(p2 Point) {
                this.enemyMoves.ExcludeAllFrom(p2)
            })
        }
    })
}

func (this *Command) PruneOutfocusedMoves() {
    //log := NewTurnLog("PruneOutfocusedMoves", "txt")

    timer := NewTimer()

    timer.Start("enemyDestinations")
    this.enemyDestinations = this.enemyMoves.Destinations()
    this.enemyDestinations.ForEach(func(p Point) {
       this.moves.ExcludeMovesTo(p)
    })
    timer.Stop()
    //log.WriteString(fmt.Sprintf("enemyDestinations: %v ms\n", timer.times["enemyDestinations"]))

    timer.Start("enemyVisibility")
    enemyVisibility := this.enemyDestinations.Visibility()
    timer.Stop()
    //log.WriteString(fmt.Sprintf("enemyVisibility: %v ms\n", timer.times["enemyVisibility"]))

    timer.Start("friendlyDestinations")
    friendlyDestinations := this.moves.Destinations().Intersection(enemyVisibility)
    timer.Stop()
    //log.WriteString(fmt.Sprintf("friendlyDestinations: %v ms\n", timer.times["friendlyDestinations"]))

    timer.Start("friendlyFocus")
    this.friendlyFocus = OpposingFocus(friendlyDestinations, this.enemyMoves)
    timer.Stop()
    //log.WriteString(fmt.Sprintf("friendlyFocus: %v ms\n", timer.times["friendlyFocus"]))

    timer.Start("enemyFocus")
    enemyFocus := OpposingFocus(this.enemyDestinations, this.moves)
    changedPoints := new(PointSet)
    timer.Stop()
    //log.WriteString(fmt.Sprintf("enemyFocus: %v ms\n", timer.times["enemyFocus"]))

    for i := 0; i < 10; i++ {
        //log.WriteString(fmt.Sprintf("\n\n*** iteration %v ***\n\n", i))

        timer.Start("maxFriendlyFocus")
        this.maxFriendlyFocus = MaxFocus(this.enemyDestinations, enemyFocus)
        this.maxFriendlyFocus_STAY = MaxFocus(this.enemies, enemyFocus)
        timer.Stop()
        //log.WriteString(fmt.Sprintf("maxFriendlyFocus: %v ms\n", timer.times["maxFriendlyFocus"]))

        changed := false

        timer.Start("excluding moves")
        this.moves.ForEach(func(move Move) {
            p := move.Destination()
            switch {
            case this.army.IsScoutAt(move.from):
                if this.friendlyFocus.At(p) >= this.maxFriendlyFocus.At(p) {
                    this.moves.Exclude(move)
                    changed = true
                    ForEachPointWithinRadius2(p, 19, func(p2 Point) {
                        changedPoints.Include(p2)
                    })
                }
            case this.army.IsSoldierAt(move.from):
                if this.friendlyFocus.At(p) > this.maxFriendlyFocus.At(p) {
                    this.moves.Exclude(move)
                    changed = true
                    ForEachPointWithinRadius2(p, 19, func(p2 Point) {
                        changedPoints.Include(p2)
                    })
                }
            case this.army.IsBerzerkerAt(move.from):
               if this.friendlyFocus.At(p) > this.maxFriendlyFocus_STAY.At(p) {
                   this.moves.Exclude(move)
                   changed = true
                   ForEachPointWithinRadius2(p, 19, func(p2 Point) {
                       changedPoints.Include(p2)
                   })
               }
            default:
                panic("What is it then?")
            }
        })
        timer.Stop()

        if changed {
            //timer.Start("friendlyDestinations")
            //friendlyDestinations = this.moves.ExceptFrom(berzerkers).Destinations().Intersection(enemyVisibility)
            //timer.Stop()
            //log.WriteString(fmt.Sprintf("friendlyDestinations: %v ms\n", timer.times["friendlyDestinations"]))
            //log.WriteString(fmt.Sprintf("%v\n\n", friendlyDestinations))

            timer.Start("enemyFocus")
            enemyFocus.UpdateOpposingFocus(this.enemyDestinations.Intersection(changedPoints), this.moves)
            changedPoints = new(PointSet)
            timer.Stop()
            //log.WriteString(fmt.Sprintf("enemyFocus: %v ms\n", timer.times["enemyFocus"]))
        } else {
            break
        }
    }

    //log.WriteString(fmt.Sprintf("\ntimer %v\n", timer))
}

// TODO: beseech allies to "make them pay"
// TODO: don't retreat into dead ends, short ones anyway
func (this *Command) DoomedAntsTakeHeart() {
    //log := NewTurnLog("DoomedAntsTakeHeart", "txt")

    ForEachPoint(func(p Point) {
        if this.terrain.At(p).HasFriendlyAnt() && this.moves.At(p) == 0 {
            //log.WriteString(fmt.Sprintf("Restoring %v\n", p))

            ForEachDirection(func(dir Direction) {
                move := Move{p, dir}
                s := this.terrain.At(move.Destination())
                if !s.HasWater() && !s.HasFood() {
                    this.moves.Include(move)
                }
            })
        }
    })
}

func (this *Command) PickBestMoves() {
    //log := NewTurnLog("PickBestMoves", "txt")

    // TODO
    // return an EvaluatedMoveSet
    // ignores hills!
    //timer.Start("forage")
    forageMoves := ForageMoves(this.terrain)
    //.ForEach(func(move Move) {
    //    if !this.terrain.At(move.Destination()).HasFood() {
    //        this.moves.Select(move)
    //    }
    //})
    //timer.Stop()

    //fmt.Println(forageMoves)

    //foragers := AssignForagers(this.terrain)

    //var distanceToFewerFriendliesThan [10]*TravelDistance
    //for i := byte(1); i < 10; i++ {
    //    distanceToFewerFriendliesThan[i] = DistanceToFewerFriendliesThan(i, this.terrain)
    //}

    list := this.moves.OrderedList(func(move Move) float32 {
        destination := move.Destination()

        var result float32

        switch {
        //case foragers.Includes(move.from):
        //    result += float32(this.distanceToFood.At(move.from))
        //    result -= float32(this.distanceToFood.At(destination))
        case this.reinforcement.At(move.from):
            result += float32(this.distanceToDoom.At(move.from))
            result -= float32(this.distanceToDoom.At(destination))
        //case this.distanceToTrouble.At(move.from) > 25:
        //    friendlies := this.terrain.VisibleFriendliesAt(move.from)
        //    if friendlies > 9 {
        //        friendlies = 9
        //    }
        //    result += float32(distanceToFewerFriendliesThan[friendlies].At(move.from))
        //    result -= float32(distanceToFewerFriendliesThan[friendlies].At(destination))
        default:
            result += float32(this.distanceToTrouble.At(move.from))
            result -= float32(this.distanceToTrouble.At(destination))

            // discourage foragers from following each other
            //if this.terrain.At(destination).HasFriendlyAnt() {
            //    result -= 0.1
            //}
        }

        if forageMoves.Includes(move) {
            //fmt.Printf("%v is a forage move\n", move)
            result += 19.0
        }

        switch {
        case this.army.IsScoutAt(move.from):
            fromFocus := this.friendlyFocus.At(move.from)
            if fromFocus >= this.maxFriendlyFocus.At(move.from) {
               result += float32(fromFocus * 20)
            } else {
               result -= float32(fromFocus * 20)
            }
        
            toFocus := this.friendlyFocus.At(destination)
            if toFocus >= this.maxFriendlyFocus.At(destination) {
               result -= float32(toFocus * 20)
            } else {
               result += float32(toFocus * 20)
            }
        case this.army.IsSoldierAt(move.from):
            fromFocus := this.friendlyFocus.At(move.from)
            if fromFocus > this.maxFriendlyFocus.At(move.from) {
               result += float32(fromFocus * 10)
            } else {
               result -= float32(fromFocus * 10)
            }
        
            toFocus := this.friendlyFocus.At(destination)
            if toFocus > this.maxFriendlyFocus.At(destination) {
               result -= float32(toFocus * 10)
            } else {
               result += float32(toFocus * 10)
            }
        }

        if result > 0 {
            result = (result * result)
        } else {
            result = -(result * result)
        }

        //fmt.Printf("considering move: %v -- %v\n", move, result)

        return result
    })

    list.ForBestWorst(func(move Move) bool {
        return this.moves.Includes(move)
    }, func(move Move) {
        // best move
        this.moves.Select(move)
    }, func(move Move) {
        // worst move
        if this.moves.At(move.from).IsMultiple() {
            this.moves.Exclude(move)
        }
    })
}

func (this *Command) SaveCrushedAnts() {
    ForEachPoint(func(p Point) {
        this.SaveCrushedAntsAt(p)
    })
}

func (this *Command) SaveCrushedAntsAt(p Point) {
    if this.terrain.At(p).HasFriendlyAnt() && this.moves.At(p) == 0 {
        this.moves.Select(Move{p, STAY})
        ForEachNeighbor(p, func(p2 Point) {
            this.SaveCrushedAntsAt(p2)
        })
    }
}

func (this *Command) Calculate() {
    if this.turn == turn {
        return
    }
    startTime := now()

    //log := NewLog("command", "txt")
    //log.WriteString(fmt.Sprintf("start\n"))

    timer := NewTimer()

    timer.Start("reset")
    this.Reset()
    timer.Stop()
    //log.WriteString(fmt.Sprintf("Reset: %v ms\n", timer.times["Reset"]))

    timer.Start("prune")
    this.PruneOutfocusedMoves()
    timer.Stop()
    //log.WriteString(fmt.Sprintf("PruneOutfocusedMoves: %v ms\n", timer.times["PruneOutfocusedMoves"]))

    timer.Start("reinvigorate")
    this.DoomedAntsTakeHeart()
    timer.Stop()

    timer.Start("pick")
    this.PickBestMoves()
    timer.Stop()
    //log.WriteString(fmt.Sprintf("PickBestMoves: %v ms\n", timer.times["PickBestMoves"]))

    timer.Start("save")
    this.SaveCrushedAnts()
    timer.Stop()
    //log.WriteString(fmt.Sprintf("SaveCrushedAnts: %v ms\n", timer.times["SaveCrushedAnts"]))

    this.moves.ReplaceLoops()

    this.time = now() - startTime
    this.turn = turn

    //log.WriteString(fmt.Sprintf("turn %v, total=%3v, %3v reset, %3v prune, %3v reinvigorate, %3v pick, %3v save\n", turn, this.time,
    //    timer.times["reset"], timer.times["prune"], timer.times["reinvigorate"], timer.times["pick"], timer.times["save"]))
}

func (this *Command) ForEach(f func(Move)) {
    this.moves.ForEach(func(move Move) {
        switch move.dir {
        case NORTH, EAST, SOUTH, WEST:
            f(move)
        }
    })
}

func (this *Command) String() string {
    return GridToString(func(p Point) byte {
        square := this.terrain.At(p)
        switch {
        case square.HasFood():
            return '*'
        case square.HasLand():
            return this.moves.At(p).Char()
        case square.HasWater():
            return '%'
        }
        return ' '
    })
}
