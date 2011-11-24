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
    forageScent := NewScent(terrain, mystery)
    battleScent := NewScent(terrain, mystery)
    army := NewArmy(terrain)
    command := NewCommand(terrain, forageScent, battleScent, army)
    for i := 0; i < 20; i++ {
        forageScent.Calculate()
        battleScent.CalculateBattle()
    }
    command.Calculate()
    if command.At(Point{1, 1}) != SOUTH {
        t.Error(forageScent)
        t.Error(battleScent)
        t.Error(command)
        t.Errorf("%v %v", forageScent.At(Point{1, 1}), forageScent.At(Point{2, 1}))
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
    holyGround := NewHolyGround(terrain)
    mystery := NewMystery(terrain)
    forageScent := NewForageScent(terrain, holyGround, mystery)
    battleScent := NewBattleScent(terrain, holyGround, mystery)
    army := NewArmy(terrain)
    predictions := NewPredictions(terrain)
    command := NewCommand(terrain, forageScent, battleScent, army, predictions)
    forageScent.Calculate()
    battleScent.Calculate()
    command.Reset()
    before := command.At(Point{2, 0})
    command.PruneOutfocusedMoves()
    after := command.At(Point{2, 0})
    if before.Minus(after) != EAST {
        t.Errorf("%v %v", before, after)
    }
}
