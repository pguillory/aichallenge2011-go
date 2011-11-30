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


func TestDistanceToSoldier(t *testing.T) {
    terrain := NewTerrain(
        "%%%%%%%%%%%%%%%%\n" +
        "%aaaa....b.b...%\n" +
        "%aaaaa...bbb...%\n" +
        "%aaaaaa..b.....%\n" +
        "%aaaaaa..b.....%\n" +
        "%aaaaaa........%\n" +
        "%%%%%%%%%%%%%%%%\n" +
        "%aaaa....b.....%\n" +
        "%%%%%%%%%%%%%%%%")

    expected :=
	    "++++++++++++++++\n" +
	    "+000012345++a+++\n" +
	    "+000001234789a++\n" +
	    "+0000001236789a+\n" +
	    "+00000012356789+\n" +
	    "+00000012345678+\n" +
	    "++++++++++++++++\n" +
	    "++++++++++++++++\n" +
	    "++++++++++++++++"

    army := NewArmy(terrain)
    distanceToBerzerker := DistanceToSoldier(terrain, army)

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
