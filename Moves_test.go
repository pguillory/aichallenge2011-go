package main

import "testing"

func TestMoves(t *testing.T) {
    m := MapFromString(
    "a%...%*%\n" +
    ".%.%.%.%\n" +
    ".%.%.%.%\n" +
    ".......%\n" +
    "%%%%%%%%")
    mystery := NewMystery(m)
    scent := NewScent(m, mystery)
    moves := NewMoves(m, scent)
    for i := 0; i < 20; i++ {
        scent.Iterate()
    }
    moves.Calculate()
    if moves.At(Point{0, 0}) != SOUTH {
        t.Errorf("%v", scent)
        t.Errorf("%v", moves)
        t.Errorf("%v %v", scent.At(Point{0, 0}), scent.At(Point{1, 0}))
    }
}
