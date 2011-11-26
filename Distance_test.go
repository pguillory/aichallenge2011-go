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
	"+t+nml+f+\n" +
	"+s+o+k+g+\n" +
	"+rqp+jih+\n" +
	"+++++++++"
	//"+++++++++\n" +
	//"+e+876+0+\n" +
	//"+d+9+5+1+\n" +
	//"+cba+432+\n" +
	//"+++++++++"

    mystery := NewMystery(terrain)
    potentialEnemy := NewPotentialEnemy(terrain)
    distanceToTrouble := DistanceToTrouble(terrain, mystery, potentialEnemy)

    if distanceToTrouble.String() != expected {
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
