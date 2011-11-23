package main

//import "fmt"

var spiralPattern = [21]Move{
    {Point{ 0,  0}, NORTH},
    {Point{-1,  0}, NORTH},
    {Point{ 0,  1}, EAST},
    {Point{ 1,  0}, SOUTH},
    {Point{ 0, -1}, WEST},
    {Point{-2,  0}, NORTH},
    {Point{ 0,  2}, EAST},
    {Point{ 2,  0}, SOUTH},
    {Point{ 0, -2}, WEST},
    {Point{-1,  1}, NORTH},
    {Point{ 1,  1}, EAST},
    {Point{ 1, -1}, SOUTH},
    {Point{-1, -1}, WEST},
    {Point{-2,  1}, NORTH},
    {Point{-1,  2}, EAST},
    {Point{ 1,  2}, EAST},
    {Point{ 2,  1}, SOUTH},
    {Point{ 2, -1}, SOUTH},
    {Point{ 1, -2}, WEST},
    {Point{-1, -2}, WEST},
    {Point{-2, -1}, NORTH},
}

type Command struct {
    terrain *Terrain
    workerScent, battleScent *Scent
    army *Army
    predictions *Predictions
    moves, enemyMoves *MoveSet
    //dirs [MAX_ROWS][MAX_COLS]Direction
    //len int
}

func NewCommand(terrain *Terrain, workerScent, battleScent *Scent, army *Army, predictions *Predictions) *Command {
    this := new(Command)
    this.terrain = terrain
    this.workerScent = workerScent
    this.battleScent = battleScent
    this.army = army
    this.predictions = predictions
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

/*
func FriendlyDestinations() *PointSet {
    result := new(PointSet)

    this.moves.ForEach(func(move Move) {
        result.Include(move.Destination())
    })

    return result
    friendlyDestinations := this.moves.Destinations()
}
*/

func (this *Command) PruneOutfocusedMoves() {
    //log := NewTurnLog("PruneOutfocusedMoves", "log")

    //timer := NewTimer()

    berzerkers := this.army.Berzerkers()

    //timer.Start("enemyDestinations")
    enemyDestinations := this.enemyMoves.Destinations()
    //enemyDestinations.ForEach(func(p Point) {
    //    this.moves.ExcludeMovesTo(p)
    //})
    //timer.Stop()
    //log.WriteString("enemyDestinations\n")
    //log.WriteString(fmt.Sprintf("%v\n\n", enemyDestinations))

    //timer.Start("friendlyDestinations")
    friendlyDestinations := this.moves.ExceptFrom(berzerkers).Destinations()
    //timer.Stop()
    //log.WriteString("friendlyDestinations\n")
    //log.WriteString(fmt.Sprintf("%v\n\n", friendlyDestinations))

    //timer.Start("friendlyFocus")
    friendlyFocus := OpposingFocus(friendlyDestinations, this.enemyMoves)
    //timer.Stop()
    //log.WriteString("friendlyFocus\n")
    //log.WriteString(fmt.Sprintf("%v\n\n", friendlyFocus))

    for i := 0; i < 20; i++ {
        //log.WriteString(fmt.Sprintf("iteration %v\n", i))

        //timer.Start("enemyFocus")
        enemyFocus := OpposingFocus(enemyDestinations, this.moves)
        //timer.Stop()
        //log.WriteString("enemyFocus\n")
        //log.WriteString(fmt.Sprintf("%v\n\n", enemyFocus))

        //timer.Start("maxFriendlyFocus")
        maxFriendlyFocus := MaxFocus(enemyDestinations, enemyFocus)
        //timer.Stop()
        //log.WriteString("maxFriendlyFocus\n")
        //log.WriteString(fmt.Sprintf("%v\n\n", maxFriendlyFocus))

        changed := false

        //timer.Start("excluding moves")
        friendlyDestinations.ForEach(func(p Point) {
            if friendlyFocus.At(p) >= maxFriendlyFocus.At(p) {
                //log.WriteString(fmt.Sprintf("ExcludeMovesTo(%v)\n", p))
                this.moves.ExcludeMovesTo(p)
                changed = true
            }
        })
        //timer.Stop()

        if changed {
            //timer.Start("friendlyDestinations")
            friendlyDestinations = this.moves.ExceptFrom(berzerkers).Destinations()
            //timer.Stop()
            //log.WriteString("friendlyDestinations\n")
            //log.WriteString(fmt.Sprintf("%v\n\n", friendlyDestinations))
        } else {
            break
        }
    }

    //log.WriteString(fmt.Sprintf("\ntimer %v\n", timer))
}

/*
func (this *Command) DoomedAntsTakeHeart() {
    ForEachPoint(func(p Point) {
        if this.terrain.At(p).HasFriendlyAnt() && this.moves.At(p) == 0 {
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
*/

func (this *Command) PickBestMovesByScent() {
    //log := NewTurnLog("PickBestMovesByScent", "log")

    list := this.moves.OrderedList(func(move Move) float32 {
        p2 := move.from.Neighbor(move.dir)

        if this.army.IsSoldierAt(move.from) {
            return this.battleScent.At(p2) - this.battleScent.At(move.from)
        }

        return this.workerScent.At(p2) - this.workerScent.At(move.from)
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
    //log := NewTurnLog("command", "log")
    //log.WriteString(fmt.Sprintf("start\n"))

    //timer := NewTimer()

    //timer.Start("Reset")
    this.Reset()
    //timer.Stop()
    //log.WriteString(fmt.Sprintf("Reset: %v ms\n", timer.times["Reset"]))

    //timer.Start("PruneOutfocusedMoves")
    this.PruneOutfocusedMoves()
    //timer.Stop()
    //log.WriteString(fmt.Sprintf("PruneOutfocusedMoves: %v ms\n", timer.times["PruneOutfocusedMoves"]))

    //timer.Start("DoomedAntsTakeHeart")
    //this.DoomedAntsTakeHeart()
    //timer.Stop()

    //timer.Start("PickBestMovesByScent")
    this.PickBestMovesByScent()
    //timer.Stop()
    //log.WriteString(fmt.Sprintf("PickBestMovesByScent: %v ms\n", timer.times["PickBestMovesByScent"]))

    //timer.Start("SaveCrushedAnts")
    this.SaveCrushedAnts()
    //timer.Stop()
    //log.WriteString(fmt.Sprintf("SaveCrushedAnts: %v ms\n", timer.times["SaveCrushedAnts"]))

    //log.WriteString(fmt.Sprintf("turn %v, timer %v\n", turn, timer))
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
        return this.moves.At(p).Char()
    })
}
