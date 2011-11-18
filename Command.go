package main

import "fmt"

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
    workerScent, soldierScent *Scent
    army *Army
    moves, enemyMoves *MoveSet
    //dirs [MAX_ROWS][MAX_COLS]Direction
    //len int
}

func NewCommand(terrain *Terrain, workerScent, soldierScent *Scent, army *Army) *Command {
    this := new(Command)
    this.terrain = terrain
    this.workerScent = workerScent
    this.soldierScent = soldierScent
    this.army = army
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
            this.enemyMoves.IncludeAllFrom(p)
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
    log := NewTurnLog("PruneOutfocusedMoves", "log")
    log.WriteString(fmt.Sprintf("turn %v\n", turn))

    timer := NewTimer()

    timer.Start("enemyDestinations")
    enemyDestinations := this.enemyMoves.Destinations()
    timer.Stop()
    //fmt.Println("enemyDestinations")
    //fmt.Println(enemyDestinations)

    timer.Start("friendlyDestinations")
    friendlyDestinations := this.moves.Destinations()
    timer.Stop()
    //fmt.Println("friendlyDestinations")
    //fmt.Println(friendlyDestinations)

    timer.Start("friendlyFocus")
    friendlyFocus := OpposingFocus(friendlyDestinations, this.enemyMoves)
    timer.Stop()
    //fmt.Println("friendlyFocus")
    //fmt.Println(friendlyFocus)

    for i := 0; i < 10; i++ {
        log.WriteString(fmt.Sprintf("iteration %v\n", i))

        timer.Start("enemyFocus")
        enemyFocus := OpposingFocus(enemyDestinations, this.moves)
        timer.Stop()
        //fmt.Println("enemyFocus")
        //fmt.Println(enemyFocus)

        timer.Start("maxFriendlyFocus")
        maxFriendlyFocus := MaxFocus(enemyDestinations, enemyFocus)
        timer.Stop()
        //fmt.Println("maxFriendlyFocus")
        //fmt.Println(maxFriendlyFocus)

        changed := false

        timer.Start("excluding moves")
        friendlyDestinations.ForEach(func(p Point) {
            if friendlyFocus.At(p) >= maxFriendlyFocus.At(p) {
                this.moves.ExcludeMovesTo(p)
                changed = true
            }
        })
        timer.Stop()

        if changed {
            timer.Start("friendlyDestinations (redux)")
            friendlyDestinations = this.moves.Destinations()
            timer.Stop()
            //fmt.Println("friendlyFocus")
            //fmt.Println(friendlyFocus)
        } else {
            break
        }
    }

    log.WriteString(fmt.Sprintf("timer %v\n", timer))
}

func (this *Command) PickBestMovesByScent() {
    list := this.moves.OrderedList(func(move Move) float32 {
        p2 := move.from.Neighbor(move.dir)

        if this.army.IsSoldierAt(move.from) {
            return this.soldierScent.At(p2) - this.soldierScent.At(move.from)
        }

        return this.workerScent.At(p2) - this.workerScent.At(move.from)
    })

    list.ForBestWorst(func(move Move) bool {
        return this.moves.Includes(move)
    }, func(move Move) {
        this.moves.Select(move)
    }, func(move Move) {
        this.moves.Exclude(move)
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
    log := NewTurnLog("command", "log")

    timer := NewTimer()

    timer.Start("Reset")
    this.Reset()
    timer.Stop()

    timer.Start("PruneOutfocusedMoves")
    this.PruneOutfocusedMoves()
    timer.Stop()

    timer.Start("PickBestMovesByScent")
    this.PickBestMovesByScent()
    timer.Stop()

    timer.Start("SaveCrushedAnts")
    this.SaveCrushedAnts()
    timer.Stop()

    log.WriteString(fmt.Sprintf("turn %v, timer %v\n", turn, timer))
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
