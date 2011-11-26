package main

import "testing"

func TestPotentialEnemy(t *testing.T) {
    terrain := NewTerrain(
    "%%%%%%%%%\n" +
    "%.%...%b%\n" +
    "%.%.%.%.%\n" +
    "%...%...%\n" +
    "%%%%%%%%%")

    expected :=
	"%%%%%%%%%\n" +
	"%.%...%b%\n" +
	"%.%.%b%b%\n" +
	"%...%bbb%\n" +
	"%%%%%%%%%"

    potentialEnemy := NewPotentialEnemy(terrain)
    terrain.Update(new(Terrain))

    for turn = 1; turn <= 5; turn++ {
        potentialEnemy.Calculate()
        //t.Error(potentialEnemy)
    }

    if potentialEnemy.String() != expected {
        t.Error(potentialEnemy)
    }
}
