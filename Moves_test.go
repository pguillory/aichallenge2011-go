package main

import "testing"

func TestMoves(t *testing.T) {
    m := NewMap(
    "a%...%*%\n" +
    ".%.%.%.%\n" +
    ".%.%.%.%\n" +
    ".......%\n" +
    "%%%%%%%%")
    mystery := NewMystery(m)
    workerScent := NewScent(m, mystery)
    soldierScent := NewScent(m, mystery)
    army := NewArmy(m)
    moves := NewMoves(m, workerScent, soldierScent, army)
    for i := 0; i < 20; i++ {
        workerScent.Iterate()
        soldierScent.IterateSoldier()
    }
    army.Iterate()
    moves.Calculate()
    if moves.At(Point{0, 0}) != SOUTH {
        t.Error(workerScent)
        t.Error(soldierScent)
        t.Error(moves)
        t.Errorf("%v %v", workerScent.At(Point{0, 0}), workerScent.At(Point{1, 0}))
    }
}
