package main

const (
    MAX_ROWS = 200
    MAX_COLS = 200
)

var turn = 0
var loadtime = 3000
var turntime = 1000
var rows = MAX_ROWS
var cols = MAX_COLS
var turns = 500
var viewradius2 = 55
var attackradius2 = 5
var spawnradius2 = 1

func normalizeRow(row int) int {
    for row < 0 {
        row += rows
    }
    return row % rows
}

func normalizeCol(col int) int {
    for col < 0 {
        col += cols
    }
    return col % cols
}
