package main

import "testing"

func TestNormalize(t *testing.T) {
    rows = 10
    cols = 10

    var cases = map[int]int {
        -13: 7,
        -2: 8,
        -1: 9,
        0: 0,
        1: 1,
        2: 2,
        9: 9,
        10: 0,
        11: 1,
        12: 2,
    }

    for k, v := range cases {
        if (Point{0, k}.Normalize().col != v) {
            t.Errorf("%v != %v\n", k, v)
        }
    }
}

func TestDistance2(t *testing.T) {
    rows = 10
    cols = 10

    p := Point{1, 1}

    type TestCase struct {
        p Point
        distance2 int
    }
    cases := []TestCase{
        {Point{0, 0}, 2},
        {Point{1, 1}, 0},
        {Point{1, 2}, 1},
        {Point{2, 2}, 2},
        {Point{1, 3}, 4},
        {Point{9, 9}, 8},
        {Point{7, 8}, 25},
    }
    for _, v := range cases {
        if (v.p.Distance2(p) != v.distance2) {
            t.Errorf("%+v %v\n", v, v.p.Distance2(p))
        }
    }
}

func TestNeighbor(t *testing.T) {
    rows = 10
    cols = 10

    type TestCase struct {
        from Point
        dir Direction
        to Point
    }
    cases := []TestCase{
        {Point{0, 0}, NORTH, Point{9, 0}},
        {Point{0, 0}, EAST, Point{0, 1}},
        {Point{0, 0}, SOUTH, Point{1, 0}},
        {Point{0, 0}, WEST, Point{0, 9}},
        {Point{1, 1}, EAST, Point{1, 2}},
    }
    for _, v := range cases {
        if !v.from.Neighbor(v.dir).Equals(v.to) {
            t.Errorf("%+v %v", v, v.from.Neighbor(v.dir))
        }
    }
}

func TestForEachPoint(t *testing.T) {
    rows = 11
    cols = 13

    count := 0
    ForEachPoint(func(p Point) {
        count += 1
    })
    if count != rows * cols {
        t.Fail()
    }
}

func TestForEachPointWithinManhattanDistance(t *testing.T) {
    rows = 11
    cols = 13

    count := 0
    ForEachPointWithinManhattanDistance(Point{5, 5}, 1, func(p Point) {
        count += 1
    })
    if count != 9 {
        t.Errorf("count: %v", count)
    }
}

func TestForEachPointWithinRadius2(t *testing.T) {
    rows = 11
    cols = 13

    count := 0
    ForEachPointWithinRadius2(Point{5, 5}, 5, func(p Point) {
        count += 1
    })
    if count != 21 {
        t.Errorf("count: %v", count)
    }
}

func TestForEachNeighbor(t *testing.T) {
    rows = 11
    cols = 13

    count := 0
    ForEachNeighbor(Point{0, 0}, func(p Point) {
        count += 1
    })
    if count != 5 {
        t.Errorf("count: %v", count)
    }
}
