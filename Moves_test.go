package main

import "testing"

func TestMoves(t *testing.T) {
    m := MapFromString(
    "a%...%*%\n" +
    ".%.%.%.%\n" +
    "...%...%\n" +
    ".......%\n" +
    "%%%%%%%%")
    scent := NewScent(m)
    for i := 0; i < 20; i++ {
        scent = scent.Iterate()
    }
    moves := NewMoves(m, scent)
    if moves.At(Point{0, 0}) != SOUTH {
        t.Errorf("%v", scent)
        t.Errorf("%v", moves)
    }
}
