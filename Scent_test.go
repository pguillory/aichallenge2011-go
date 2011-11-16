package main

import "testing"

func TestIterate(t *testing.T) {
/*
    m := MapFromString(
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
    m := new(Map)
    m.SeeHill(Point{180, 50}, Player(1))
    m.SeeWater(Point{179, 48})
    m.SeeWater(Point{179, 49})
    m.SeeWater(Point{179, 50})
    m.SeeWater(Point{179, 51})
    m.SeeWater(Point{179, 52})
    ForEachPoint(func(p Point) {
        if !m.At(p).IsVisible() {
            m.SeeLand(p)
        }
    })

    s := NewScent(m)
    start := now()
    for i := 0; i < 25; i++ {
        s = s.Iterate()
    }
    runtime := now() - start
    if runtime > 150 {
        t.Errorf("runtime=%v ms\n", runtime)
    }
}
