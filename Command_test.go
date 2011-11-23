package main

import "testing"

/*
func TestFollowScentThroughMaze(t *testing.T) {
    terrain := NewTerrain(
    "%%%%%%%%%\n" +
    "%a%...%*%\n" +
    "%.%.%.%.%\n" +
    "%...%...%\n" +
    "%%%%%%%%%")
    mystery := NewMystery(terrain)
    workerScent := NewScent(terrain, mystery)
    battleScent := NewScent(terrain, mystery)
    army := NewArmy(terrain)
    command := NewCommand(terrain, workerScent, battleScent, army)
    for i := 0; i < 20; i++ {
        workerScent.Calculate()
        battleScent.CalculateBattle()
    }
    command.Calculate()
    if command.At(Point{1, 1}) != SOUTH {
        t.Error(workerScent)
        t.Error(battleScent)
        t.Error(command)
        t.Errorf("%v %v", workerScent.At(Point{1, 1}), workerScent.At(Point{2, 1}))
    }
}
*/

func TestMoves(t *testing.T) {
    terrain := NewTerrain(
    "...................................................................%\n" +
    "...................................................................%\n" +
    "a...b..............................................................%\n" +
    "...................................................................%\n" +
    "...................................................................%\n" +
    "%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%")
    mystery := NewMystery(terrain)
    workerScent := NewScent(terrain, mystery)
    battleScent := NewScent(terrain, mystery)
    army := NewArmy(terrain)
    predictions := NewPredictions(terrain)
    command := NewCommand(terrain, workerScent, battleScent, army, predictions)
    workerScent.Calculate()
    battleScent.CalculateBattle()
    command.Reset()
    before := command.At(Point{2, 0})
    command.PruneOutfocusedMoves()
    after := command.At(Point{2, 0})
    if before.Minus(after) != EAST {
        t.Errorf("%v %v", before, after)
    }
}
