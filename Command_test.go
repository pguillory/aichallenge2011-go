package main

import "testing"

func TestFollowScentThroughMaze(t *testing.T) {
    terrain := NewTerrain(
    "%%%%%%%%%\n" +
    "%a%...%*%\n" +
    "%.%.%.%.%\n" +
    "%...%...%\n" +
    "%%%%%%%%%")

    //distanceToEnemy := DistanceToEnemy(terrain)
    //distanceToFriendlyHill := DistanceToFriendlyHill(terrain)
    //mystery := NewMystery(terrain)
    //forageScent := NewForageScent(terrain, distanceToEnemy, distanceToFriendlyHill, mystery)
    //battleScent := NewBattleScent(terrain, distanceToEnemy, distanceToFriendlyHill, mystery)
    army := NewArmy(terrain)
    predictions := NewPredictions(terrain)
    mystery := NewMystery(terrain)
    potentialEnemy := NewPotentialEnemy(terrain)
    distanceToFood := DistanceToFood(terrain)
    distanceToTrouble := DistanceToTrouble(terrain, mystery, potentialEnemy)
    distanceToDoom := DistanceToTrouble(terrain, mystery, potentialEnemy)
    rageVirus := NewRageVirus(terrain, army, distanceToTrouble)
    command := NewCommand(terrain, army, predictions, distanceToFood, distanceToTrouble, distanceToDoom, rageVirus)

    if command.At(Point{1, 1}) != SOUTH {
        //t.Error(forageScent)
        //t.Error(battleScent)
        t.Error(command)
        //t.Errorf("%v %v", forageScent.At(Point{1, 1}), forageScent.At(Point{2, 1}))
    }
}

func TestMoves(t *testing.T) {
    terrain := NewTerrain(
    "...................................................................%\n" +
    "...................................................................%\n" +
    "a...b..............................................................%\n" +
    "...................................................................%\n" +
    "...................................................................%\n" +
    "%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%")

    //distanceToEnemy := DistanceToEnemy(terrain)
    //distanceToFriendlyHill := DistanceToFriendlyHill(terrain)
    //mystery := NewMystery(terrain)
    //forageScent := NewForageScent(terrain, distanceToEnemy, distanceToFriendlyHill, mystery)
    //battleScent := NewBattleScent(terrain, distanceToEnemy, distanceToFriendlyHill, mystery)
    army := NewArmy(terrain)
    predictions := NewPredictions(terrain)
    mystery := NewMystery(terrain)
    potentialEnemy := NewPotentialEnemy(terrain)
    distanceToFood := DistanceToFood(terrain)
    distanceToTrouble := DistanceToTrouble(terrain, mystery, potentialEnemy)
    distanceToDoom := DistanceToTrouble(terrain, mystery, potentialEnemy)
    rageVirus := NewRageVirus(terrain, army, distanceToTrouble)
    command := NewCommand(terrain, army, predictions, distanceToFood, distanceToTrouble, distanceToDoom, rageVirus)

    command.Reset()
    before := command.At(Point{2, 0})
    command.PruneOutfocusedMoves()
    after := command.At(Point{2, 0})
    if before.Minus(after) != EAST {
        t.Errorf("%v -> %v", before, after)
    }
}
