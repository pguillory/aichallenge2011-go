package main

import "testing"

func TestForageMoves(t *testing.T) {
    terrain := NewTerrain(
        "%%%%%%%%%%%%%%%%%%%%%%%\n" +
        "%.....................%\n" +
        "%.*...*...............%\n" +
        "%.....................%\n" +
        "%.....A...............%\n" +
        "%....a................%\n" +
        "%.................*...%\n" +
        "%.....................%\n" +
        "%.....................%\n" +
        "%%%%%%%%%%%%%%%%%%%%%%%")
    expected :=
        ".......................\n" +
        ".......................\n" +
        ".......................\n" +
        ".......................\n" +
        "......N................\n" +
        ".....E.................\n" +
        ".......................\n" +
        ".......................\n" +
        ".......................\n" +
        "......................."

    foraging := NewForaging(terrain)

    if foraging.moves.String() != expected {
        t.Error(foraging.moves)
    }
}

func TestForageUsingFutureSpawns(t *testing.T) {
    terrain := NewTerrain(
        "%%%%%%%%%%%%%%%%%%%%%%%\n" +
        "%.....................%\n" +
        "%.*...*...............%\n" +
        "%.....................%\n" +
        "%.....A...............%\n" +
        "%...a.................%\n" +
        "%.................*...%\n" +
        "%.....................%\n" +
        "%.....................%\n" +
        "%%%%%%%%%%%%%%%%%%%%%%%")
    expected :=
        ".......................\n" +
        ".......................\n" +
        ".......................\n" +
        ".......................\n" +
        "......N................\n" +
        "....N..................\n" +
        ".......................\n" +
        ".......................\n" +
        ".......................\n" +
        "......................."

    foraging := NewForaging(terrain)

    if foraging.moves.String() != expected {
        t.Error(foraging.moves)
    }
}

func TestForageAroundObstacles(t *testing.T) {
    terrain := NewTerrain(
        "%%%%%%%%%%%%%%%%%%%%%%%%%\n" +
        "%.....%.................%\n" +
        "%.*%..%.a...............%\n" +
        "%%%%..%a%...............%\n" +
        "%.....%%................%\n" +
        "%.......................%\n" +
        "%.......................%\n" +
        "%.......................%\n" +
        "%.......................%\n" +
        "%%%%%%%%%%%%%%%%%%%%%%%%%")
    expected :=
        ".........................\n" +
        ".........................\n" +
        "........E................\n" +
        ".........................\n" +
        ".........................\n" +
        ".........................\n" +
        ".........................\n" +
        ".........................\n" +
        ".........................\n" +
        "........................."

    foraging := NewForaging(terrain)

    if foraging.moves.String() != expected {
        t.Error(foraging.moves)
    }
}

func TestSpawnedNextToFood(t *testing.T) {
    terrain := NewTerrain(
        "%%%%%%%%%%%%%%%%%%%%%%%%%\n" +
        "%...........*...........%\n" +
        "%.......................%\n" +
        "%.......................%\n" +
        "%...........A...........%\n" +
        "%...........*...........%\n" +
        "%.......................%\n" +
        "%.......................%\n" +
        "%.......................%\n" +
        "%%%%%%%%%%%%%%%%%%%%%%%%%")
    expected := EAST | WEST | STAY

    foraging := NewForaging(terrain)

    if foraging.moves.At(Point{4, 12}) != expected {
        t.Error(foraging.moves)
    }
}
