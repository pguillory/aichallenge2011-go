package main

import "testing"

func TestTravelDistance(t *testing.T) {
    terrain := NewTerrain(
    "%%%%%%%%%\n" +
    "%a%...%b%\n" +
    "%.%.%.%.%\n" +
    "%...%...%\n" +
    "%%%%%%%%%")

    mystery := NewMystery(terrain)
    potentialEnemy := NewPotentialEnemy(terrain)
    distanceToTrouble := DistanceToTrouble(terrain, mystery, potentialEnemy)

    if distanceToTrouble.At(Point{1, 1}) <= distanceToTrouble.At(Point{2, 1}) {
        t.Error(terrain)
        t.Error(distanceToTrouble)
    }
}

func TestAssignForagers(t *testing.T) {
    terrain := NewTerrain(
    "%%%%%%%%%%%%%%%%%\n" +
    "%a.....*%a......%\n" +
    "%.......%.......%\n" +
    "%....a..%....a..%\n" +
    "%%%%%%%%%%%%%%%%%")

    expected :=
    ".................\n" +
    ".................\n" +
    ".................\n" +
    ".....x...........\n" +
    "................."

    foragers := AssignForagers(terrain)
    distanceToFood := DistanceToFood(terrain)

    if foragers.String() != expected {
        t.Error(terrain)
        t.Error(foragers)
        t.Error(distanceToFood)
    }
}
