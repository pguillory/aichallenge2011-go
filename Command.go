package main

type Command struct {
    terrain *Terrain
    workerScent, soldierScent *Scent
    army *Army
    moves *MoveSet
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

    ForEachPoint(func(p Point) {
        if this.terrain.At(p).HasFriendlyAnt() {
            this.moves.IncludeAllFrom(p)
        }
    })

    ForEachPoint(func(p Point) {
        s := this.terrain.At(p)
        if s.HasWater() || s.HasFood() {
            this.moves.ExcludeMovesTo(p)
        }
    })
}

/*
func Focus() {
    ForEachPoint(func(p Point) {
        if this.terrain.At(p).HasFriendlyAnt() {
            ForEachPoint(func(p Point) {
        }
    }
}
*/

func (this *Command) PickBestMovesByScent() {
    list := this.moves.Order(func(move Move) float32 {
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

func (this *Command) Calculate() {
    this.Reset()
    //friendlyFocus = NewFocus(friendlyAnts, enemyAnts)
    this.PickBestMovesByScent()
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
