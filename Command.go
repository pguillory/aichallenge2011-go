package main

//import "fmt"

type Command struct {
    time int64
    turn int
    terrain *Terrain
    forageScent *Scent
    army *Army
    predictions *Predictions
    moves, enemyMoves *MoveSet
    enemyDestinations *PointSet
    friendlyFocus, maxFriendlyFocus *Focus
    //dirs [MAX_ROWS][MAX_COLS]Direction
    //len int
}

func NewCommand(terrain *Terrain, forageScent *Scent, army *Army, predictions *Predictions) *Command {
    this := new(Command)
    this.terrain = terrain
    this.forageScent = forageScent
    //this.battleScent = battleScent
    this.army = army
    this.predictions = predictions

    this.Calculate()
    return this
}

func (this *Command) At(p Point) Direction {
    return this.moves.At(p)
}

func (this *Command) Reset() {
    this.moves = new(MoveSet)
    this.enemyMoves = new(MoveSet)

    ForEachPoint(func(p Point) {
        s := this.terrain.At(p)
        if s.HasFriendlyAnt() {
            this.moves.IncludeAllFrom(p)
        } else if s.HasEnemyAnt() {
            dir := this.predictions.At(p)
            this.enemyMoves.Include(Move{p, dir})
            //this.enemyMoves.IncludeAllFrom(p)
        }
    })

    ForEachPoint(func(p Point) {
        s := this.terrain.At(p)
        if s.HasWater() || s.HasFood() {
            this.moves.ExcludeMovesTo(p)
            this.enemyMoves.ExcludeMovesTo(p)
        }
    })
}

func (this *Command) PruneOutfocusedMoves() {
    //log := NewTurnLog("PruneOutfocusedMoves", "txt")

    timer := NewTimer()

    timer.Start("berzerkers")
    berzerkers := this.army.Berzerkers()
    timer.Stop()
    //log.WriteString(fmt.Sprintf("berzerkers: %v ms\n", timer.times["berzerkers"]))
    //log.WriteString(fmt.Sprintf("%v\n\n", berzerkers))

    timer.Start("enemyDestinations")
    this.enemyDestinations = this.enemyMoves.Destinations()
    this.enemyDestinations.ForEach(func(p Point) {
       this.moves.ExcludeMovesTo(p)
    })
    timer.Stop()
    //log.WriteString(fmt.Sprintf("enemyDestinations: %v ms\n", timer.times["enemyDestinations"]))
    //log.WriteString(fmt.Sprintf("%v\n\n", this.enemyDestinations))

    timer.Start("enemyVisibility")
    enemyVisibility := this.enemyDestinations.Visibility()
    timer.Stop()
    //log.WriteString(fmt.Sprintf("enemyVisibility: %v ms\n", timer.times["enemyVisibility"]))
    //log.WriteString(fmt.Sprintf("%v\n\n", enemyVisibility))

    timer.Start("friendlyDestinations")
    friendlyDestinations := this.moves.ExceptFrom(berzerkers).Destinations().Intersection(enemyVisibility)
    timer.Stop()
    //log.WriteString(fmt.Sprintf("friendlyDestinations: %v ms\n", timer.times["friendlyDestinations"]))
    //log.WriteString(fmt.Sprintf("%v\n\n", friendlyDestinations))

    timer.Start("friendlyFocus")
    this.friendlyFocus = OpposingFocus(friendlyDestinations, this.enemyMoves)
    timer.Stop()
    //log.WriteString(fmt.Sprintf("friendlyFocus: %v ms\n", timer.times["friendlyFocus"]))
    //log.WriteString(fmt.Sprintf("%v\n\n", this.friendlyFocus))

    for i := 0; i < 10; i++ {
        //log.WriteString(fmt.Sprintf("\n\n*** iteration %v ***\n\n", i))

        timer.Start("enemyFocus")
        enemyFocus := OpposingFocus(this.enemyDestinations, this.moves)
        timer.Stop()
        //log.WriteString(fmt.Sprintf("enemyFocus: %v ms\n", timer.times["enemyFocus"]))
        //log.WriteString(fmt.Sprintf("%v\n\n", enemyFocus))

        timer.Start("maxFriendlyFocus")
        this.maxFriendlyFocus = MaxFocus(this.enemyDestinations, enemyFocus)
        timer.Stop()
        //log.WriteString(fmt.Sprintf("maxFriendlyFocus: %v ms\n", timer.times["maxFriendlyFocus"]))
        //log.WriteString(fmt.Sprintf("%v\n\n", this.maxFriendlyFocus))

        changed := false

        timer.Start("excluding moves")
        friendlyDestinations.ForEach(func(p Point) {
            if this.friendlyFocus.At(p) >= this.maxFriendlyFocus.At(p) {
                //log.WriteString(fmt.Sprintf("ExcludeMovesTo(%v)\n", p))
                this.moves.ExcludeMovesTo(p)
                changed = true
            }
        })
        timer.Stop()

        if changed {
            //timer.Start("friendlyDestinations")
            //friendlyDestinations = this.moves.ExceptFrom(berzerkers).Destinations().Intersection(enemyVisibility)
            //timer.Stop()
            //log.WriteString(fmt.Sprintf("friendlyDestinations: %v ms\n", timer.times["friendlyDestinations"]))
            //log.WriteString(fmt.Sprintf("%v\n\n", friendlyDestinations))
        } else {
            break
        }
    }

    //log.WriteString(fmt.Sprintf("\ntimer %v\n", timer))
}

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

func (this *Command) ValueAt(p Point) float32 {
    if this.friendlyFocus.At(p) >= this.maxFriendlyFocus.At(p) {
        return this.forageScent.At(p) - float32(this.friendlyFocus.At(p)) * 1e30
    }
    return this.forageScent.At(p) + float32(this.friendlyFocus.At(p)) * 1e30
}

/*
func (this *Command) ArmyValueAt(p Point) float32 {
    if this.friendlyFocus.At(p) >= this.maxFriendlyFocus.At(p) {
        return this.battleScent.At(p) - float32(this.friendlyFocus.At(p)) * 1e30
    }
    return this.battleScent.At(p) + float32(this.friendlyFocus.At(p)) * 1e30
}
*/

func (this *Command) PickBestMovesByScent() {
    //log := NewTurnLog("PickBestMovesByScent", "txt")

    list := this.moves.OrderedList(func(move Move) float32 {
        //if this.forageScent.At(move.from.Neighbor(NORTH)) < 1e-40 && 
        //    this.forageScent.At(move.from.Neighbor(EAST)) < 1e-40 &&
        //    this.forageScent.At(move.from.Neighbor(SOUTH)) < 1e-40 &&
        //    this.forageScent.At(move.from.Neighbor(WEST)) < 1e-40 {

        //if this.army.IsSoldierAt(move.from) {
        //    return this.ArmyValueAt(move.Destination()) - this.ArmyValueAt(move.from)
        //}
        return this.ValueAt(move.Destination()) - this.ValueAt(move.from)
    })

    list.ForBestWorst(func(move Move) bool {
        return this.moves.Includes(move)
    }, func(move Move) {
        //log.WriteString(fmt.Sprintf("Select %v\n", move))
        this.moves.Select(move)
    }, func(move Move) {
        if this.moves.At(move.from).IsMultiple() {
            //log.WriteString(fmt.Sprintf("Exclude %v\n", move))
            this.moves.Exclude(move)
        } else {
            //log.WriteString(fmt.Sprintf("Exclude %v -- skipped, not multiple\n", move))
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

//        grid(moves[0], p) = STAY;
//        p2 = neighbor(p, NORTH); exclude_move(p2, SOUTH); save_crushed_ant(p2);
//        p2 = neighbor(p, EAST);  exclude_move(p2, WEST);  save_crushed_ant(p2);
//        p2 = neighbor(p, SOUTH); exclude_move(p2, NORTH); save_crushed_ant(p2);
//        p2 = neighbor(p, WEST);  exclude_move(p2, EAST);  save_crushed_ant(p2);

func (this *Command) Calculate() {
    this.forageScent.Calculate()
    //this.battleScent.Calculate()
    this.army.Calculate()
    this.predictions.Calculate()

    if this.turn == turn {
        return
    }
    startTime := now()

    //log := NewTurnLog("command", "txt")
    //log.WriteString(fmt.Sprintf("start\n"))

    timer := NewTimer()

    timer.Start("Reset")
    this.Reset()
    timer.Stop()
    //log.WriteString(fmt.Sprintf("Reset: %v ms\n", timer.times["Reset"]))

    timer.Start("PruneOutfocusedMoves")
    this.PruneOutfocusedMoves()
    timer.Stop()
    //log.WriteString(fmt.Sprintf("PruneOutfocusedMoves: %v ms\n", timer.times["PruneOutfocusedMoves"]))

    //timer.Start("DoomedAntsTakeHeart")
    //this.DoomedAntsTakeHeart()
    //timer.Stop()

    timer.Start("PickBestMovesByScent")
    this.PickBestMovesByScent()
    timer.Stop()
    //log.WriteString(fmt.Sprintf("PickBestMovesByScent: %v ms\n", timer.times["PickBestMovesByScent"]))

    timer.Start("SaveCrushedAnts")
    this.SaveCrushedAnts()
    timer.Stop()
    //log.WriteString(fmt.Sprintf("SaveCrushedAnts: %v ms\n", timer.times["SaveCrushedAnts"]))

    //log.WriteString(fmt.Sprintf("turn %v, timer %v\n", turn, timer))

    this.time = now() - startTime
    this.turn = turn
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
