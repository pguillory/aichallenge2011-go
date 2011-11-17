package main

import "testing"

func TestCalculate(t *testing.T) {
/*
    terrain := NewTerrain(
    ".......................................\n" +
    ".......................................\n" +
    ".......................................\n" +
    ".......................................\n" +
    "................%%%%%..................\n" +
    "..................1....................\n" +
    ".......................................\n" +
    ".......................................\n" +
    ".......................................\n" +
    ".......................................")
*/

    rows = 200
    cols = 200
    terrain := new(Terrain)
    terrain.SeeHill(Point{180, 50}, Player(1))
    terrain.SeeWater(Point{179, 48})
    terrain.SeeWater(Point{179, 49})
    terrain.SeeWater(Point{179, 50})
    terrain.SeeWater(Point{179, 51})
    terrain.SeeWater(Point{179, 52})
    ForEachPoint(func(p Point) {
        if !terrain.At(p).IsVisible() {
            terrain.SeeLand(p)
        }
    })

    mystery := NewMystery(terrain)
    scent := NewScent(terrain, mystery)
    start := now()
    for i := 0; i < 25; i++ {
        scent.Calculate()
    }
    runtime := now() - start
    if runtime > 150 {
        t.Errorf("runtime=%v ms\n", runtime)
    }
}
