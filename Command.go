/*
TODO:
break into two modules
tactical
scent-based

don't rampage wastefully
eliminate swaps
prune berzerker-1 moves assuming the enemy will STAY
rush a hill if someone else is bout to cap it
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
        switch {
        case s.HasWater() || s.HasFood():
            this.moves.ExcludeMovesTo(p)
            this.enemyMoves.ExcludeMovesTo(p)
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

    //foragers := AssignForagers(this.terrain)

    //var distanceToFewerFriendliesThan [10]*TravelDistance
    //for i := byte(1); i < 10; i++ {
    //    distanceToFewerFriendliesThan[i] = DistanceToFewerFriendliesThan(i, this.terrain)
    //}

    list := this.moves.OrderedList(func(move Move) float32 {
        destination := move.Destination()

        var result float32

        result += float32(this.distanceToTrouble.At(move.from))
        result -= float32(this.distanceToTrouble.At(destination))

        if forageMoves.Includes(move) {
            result += 19.0
        }

        if result > 0 {
            return (result * result)
        }
        return -(result * result)
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

    timer.Start("pick")
    this.PickBestMoves()
    timer.Stop()
    //log.WriteString(fmt.Sprintf("PickBestMoves: %v ms\n", timer.times["PickBestMoves"]))

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
