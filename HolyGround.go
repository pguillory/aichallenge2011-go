package main

//import "fmt"

const HOLYGROUND_MAX = 255

type HolyGround struct {
    terrain *Terrain
    value [MAX_ROWS][MAX_COLS]byte
    friendlyHills int
}

func NewHolyGround(terrain *Terrain) *HolyGround {
    this := new(HolyGround)
    this.terrain = terrain
    return this
}

func (this *HolyGround) At(p Point) byte {
    return HOLYGROUND_MAX - this.value[p.row][p.col]
}

func (this *HolyGround) Calculate() {
    queue := new(PointQueue)

    ForEachPoint(func(p Point) {
        if this.terrain.At(p).HasFriendlyHill() {
            this.value[p.row][p.col] = HOLYGROUND_MAX
            queue.Push(p)
        } else {
            this.value[p.row][p.col] = 0
        }
    })

    queue.ForEach(func(p Point) {
        v := this.value[p.row][p.col]
        if v > HOLYGROUND_MAX - 30 {
            ForEachNeighbor(p, func(p2 Point) {
                if this.value[p2.row][p2.col] < v - 1 && this.terrain.At(p2).HasLand() {
                    this.value[p2.row][p2.col] = v - 1
                    queue.Push(p2)
                }
            })
        }
    })
}

func (this *HolyGround) String() string {
    return GridToString(func(p Point) byte {
        v := this.At(p)
        switch {
        case v < 10:
            return '0' + v
        case v < 36:
            return 'a' + v - 10
        }
        return '+'
    })
}
