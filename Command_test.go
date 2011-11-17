package main

import "testing"

func TestMoves(t *testing.T) {
    terrain := NewTerrain(
    "a%...%*%\n" +
    ".%.%.%.%\n" +
    "...%...%\n" +
    "%%%%%%%%")
    mystery := NewMystery(terrain)
    workerScent := NewScent(terrain, mystery)
    soldierScent := NewScent(terrain, mystery)
    army := NewArmy(terrain)
    command := NewCommand(terrain, workerScent, soldierScent, army)
    for i := 0; i < 20; i++ {
        workerScent.Calculate()
        soldierScent.CalculateSoldier()
    }
    army.Calculate()
    command.Calculate()
    if command.At(Point{0, 0}) != SOUTH {
        t.Error(workerScent)
        t.Error(soldierScent)
        t.Error(command)
        t.Errorf("%v %v", workerScent.At(Point{0, 0}), workerScent.At(Point{1, 0}))
    }
}
