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

func TestDistanceToBerzerker(t *testing.T) {
    terrain := NewTerrain(
        "%%%%%%%%%%%%%%%%%\n" +
        "%aaaaa....b.b...%\n" +
        "%aaaaaa...bbb...%\n" +
        "%aaaaaaa..b.....%\n" +
        "%aaaaaaa..b.....%\n" +
        "%aaaaaaa........%\n" +
        "%%%%%%%%%%%%%%%%%\n" +
        "%aaaaa....b.....%\n" +
        "%%%%%%%%%%%%%%%%%")

    expected :=
	    "+++++++++++++++++\n" +
	    "+0000012345+babc+\n" +
	    "+0000001234789ab+\n" +
	    "+00000001236789a+\n" +
	    "+000000012356789+\n" +
	    "+000000012345678+\n" +
	    "+++++++++++++++++\n" +
	    "+++++++++++++++++\n" +
	    "+++++++++++++++++"

    army := NewArmy(terrain)
    distanceToBerzerker := DistanceToBerzerker(terrain, army)

    if distanceToBerzerker.String() != expected {
        t.Error(army)
        t.Error(distanceToBerzerker)
        t.Error(GridToString(func(p Point) byte {
            if army.IsBerzerkerAt(p) {
                return '#'
            }
            return '.'
        }))
        t.Error(GridToString(func(p Point) byte {
            return byte('0' + army.CountAt(p))
        }))
        t.Error(army.CountAt(Point{3, 14}))
    }
}
