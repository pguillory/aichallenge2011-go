package main

import "testing"

func TestPointSet(t *testing.T) {
    rows = 10
    cols = 10

    points := new(PointSet)
    points.Include(Point{1, 1})

    expected :=
	"..........\n" +
	".x........\n" +
	"..........\n" +
	"..........\n" +
	"..........\n" +
	"..........\n" +
	"..........\n" +
	"..........\n" +
	"..........\n" +
	".........."

    if points.String() != expected {
        t.Error(points)
    }
    if !points.Includes(Point{1, 1}) {
        t.Error("!points.Includes(Point{1, 1})")
    }
    if points.Includes(Point{0, 0}) {
        t.Error("points.Includes(Point{0, 0})")
    }
}
