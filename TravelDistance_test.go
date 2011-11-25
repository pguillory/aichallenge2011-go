package main

import "testing"

func TestTravelDistance(t *testing.T) {
    terrain := NewTerrain(
    "%%%%%%%%%\n" +
    "%a%...%b%\n" +
    "%.%.%.%.%\n" +
    "%...%...%\n" +
    "%%%%%%%%%")

    expected :=
	"+++++++++\n" +
	"+e+876+0+\n" +
	"+d+9+5+1+\n" +
	"+cba+432+\n" +
	"+++++++++"

    distanceToEnemy := DistanceToEnemy(terrain)

    if distanceToEnemy.String() != expected {
        t.Error(distanceToEnemy)
    }
}
