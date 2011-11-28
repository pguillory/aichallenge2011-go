package main

type Repulsion struct {
    time int64
    turn int
    terrain *Terrain
    north, south, east, west [MAX_ROWS][MAX_COLS]byte
}

func NewRepulsion(terrain *Terrain) *Repulsion {
    this := new(Repulsion)
    this.terrain = terrain

    this.Calculate()
    return this
}

func (this *Repulsion) Calculate() {
    ForEachPoint(func(p Point) {
        this.north[p.row][p.col] = 0
        this.east[p.row][p.col] = 0
        this.south[p.row][p.col] = 0
        this.west[p.row][p.col] = 0
    })

    ForEachPoint(func(p Point) {
        if this.terrain.At(p).HasFriendlyAnt() {
            ForEachPointWithinRadius2(p, viewradius2, func(p2 Point) {
                dRow := (p.row - p2.row + rows + rows / 2) % rows - rows / 2
                switch {
                case dRow < 0:
                    this.north[p2.row][p2.col] += 1
                case dRow > 0:
                    this.south[p2.row][p2.col] += 1
                }

                dCol := (p.col - p2.col + cols + cols / 2) % cols - cols / 2
                switch {
                case dCol < 0:
                    this.west[p2.row][p2.col] += 1
                case dCol > 0:
                    this.east[p2.row][p2.col] += 1
                }
            })
        }
    })
}

func (this *Repulsion) To(move Move) int {
    switch move.dir {
    case NORTH:
        return int(this.north[move.from.row][move.from.col])
    case EAST:
        return int(this.east[move.from.row][move.from.col])
    case SOUTH:
        return int(this.south[move.from.row][move.from.col])
    case WEST:
        return int(this.west[move.from.row][move.from.col])
    }
    return 0
}

func (this *Repulsion) StringFor(dir Direction) string {
    return GridToString(func(p Point) byte {
        return '0' + byte(this.To(Move{p, dir}))
    })
}

func (this *Repulsion) String() string {
    return this.StringFor(NORTH)
}
