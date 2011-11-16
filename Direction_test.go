package main

import "testing"

func TestDirection(t *testing.T) {
    dir := NORTH
    if (dir.Backward() != SOUTH) {
        t.Fail()
    }
    if (dir.Char() != 'N') {
        t.Fail()
    }
}

func TestForEachDirection(t *testing.T) {
    seen := make(map[string]int)
    ForEachDirection(func(dir Direction) {
        seen[dir.String()] = 1
    })
    if (len(seen) != 5 || seen["N"] != 1 || seen["E"] != 1 || seen["S"] != 1 || seen["W"] != 1 || seen["X"] != 1) {
        t.Errorf("%v", seen)
    }
}
