package main

import "fmt"

type DistanceGrid [MAX_ROWS][MAX_COLS]Distance

func (this *DistanceGrid) At(p Point) Distance {
    return this[p.row][p.col]
}

func (this *DistanceGrid) SetAt(p Point, value Distance) {
    this[p.row][p.col] = value
}

func (this *DistanceGrid) String() string {
    return GridToString(func(p Point) byte {
        v := this.At(p)
        if v == MAX_TRAVEL_DISTANCE {
            return '!'
        }
        if v < 10 {
            return '0' + byte(v)
        }
        v -= 10
        v /= 10
        if v < 26 {
            return 'a' + byte(v)
        }
        return '+'
    })
}


type DistanceGridGrid [MAX_ROWS][MAX_COLS]*DistanceGrid

func (this *DistanceGridGrid) At(p Point) *DistanceGrid {
    return this[p.row][p.col]
}

func (this *DistanceGridGrid) SetAt(p Point, value *DistanceGrid) {
    this[p.row][p.col] = value
}


type SearchNode struct {
    point Point
    neighbors [4]*SearchNode
    hasFriendlyAnt bool
    distance *Distance
}

func (this *SearchNode) String() string {
    return fmt.Sprintf("Node")
}

type SearchEngine struct {
    nodes [MAX_ROWS][MAX_COLS]SearchNode
}

func NewSearchEngine(terrain *Terrain) *SearchEngine {
    this := new(SearchEngine)
    ForEachPoint(func(p Point) {
        node := &this.nodes[p.row][p.col]
        node.point = p

        if !terrain.At(p).HasWater() {
            var p2 Point
            p2 = p.Neighbor(NORTH); this.nodes[p2.row][p2.col].neighbors[0] = node
            p2 = p.Neighbor(EAST);  this.nodes[p2.row][p2.col].neighbors[1] = node
            p2 = p.Neighbor(SOUTH); this.nodes[p2.row][p2.col].neighbors[2] = node
            p2 = p.Neighbor(WEST);  this.nodes[p2.row][p2.col].neighbors[3] = node
        }
        node.hasFriendlyAnt = terrain.At(p).HasFriendlyAnt()
    })
    return this
}

const SEARCH_BUFFER_SIZE = 500

func (this *SearchEngine) SearchFrom(destinations *PointSet, f func(Point, Distance) bool) {
    distances := new(DistanceGrid)

    var queue [SEARCH_BUFFER_SIZE]*SearchNode
    var tail int

    var p Point
    for p.row = 0; p.row < rows; p.row++ {
        nodeRow := &this.nodes[p.row]
        distanceRow := &distances[p.row]
        destinationRow := &destinations[p.row]
        for p.col = 0; p.col < cols; p.col++ {
            node := &nodeRow[p.col]
            node.distance = &distanceRow[p.col]
            if destinationRow[p.col] {
                queue[tail] = node
                tail++
            } else {
                *node.distance = MAX_TRAVEL_DISTANCE
            }
        }
    }

    var head int
    for head != tail {
        node := queue[head]
        head += 1
        head %= SEARCH_BUFFER_SIZE
        distance2 := *node.distance + 1
        for _, node2 := range node.neighbors {
            if node2 != nil && *node2.distance > distance2 {
                *node2.distance = distance2
                if !node2.hasFriendlyAnt && f(node2.point, distance2) {
                    queue[tail] = node2
                    tail += 1
                    tail %= SEARCH_BUFFER_SIZE
                }
            }
        }
    }
}

func (this *SearchEngine) DistanceTo(destinations *PointSet) *DistanceGrid {
    distances := new(DistanceGrid)

    var queue [SEARCH_BUFFER_SIZE]*SearchNode
    var tail int

    var p Point
    for p.row = 0; p.row < rows; p.row++ {
        nodeRow := &this.nodes[p.row]
        distanceRow := &distances[p.row]
        destinationRow := &destinations[p.row]
        for p.col = 0; p.col < cols; p.col++ {
            node := &nodeRow[p.col]
            node.distance = &distanceRow[p.col]
            if destinationRow[p.col] {
                queue[tail] = node
                tail++
            } else {
                *node.distance = MAX_TRAVEL_DISTANCE
            }
        }
    }

    var head int
    for head != tail {
        node := queue[head]
        head += 1
        head %= SEARCH_BUFFER_SIZE
        distance2 := *node.distance + 1
        for _, node2 := range node.neighbors {
            if node2 != nil && *node2.distance > distance2 {
                *node2.distance = distance2
                if !node2.hasFriendlyAnt {
                    queue[tail] = node2
                    tail += 1
                    tail %= SEARCH_BUFFER_SIZE
                }
            }
        }
    }

    return distances
}
