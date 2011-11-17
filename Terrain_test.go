package main

import "testing"

func TestFromString(t *testing.T) {
    s := "?.%*\naB2?"
    terrain := NewTerrain(s)
    s2 := terrain.String()
    if s != s2 {
        t.Errorf("%#v != %#v", s, s2)
    }
}

func TestUpdate(t *testing.T) {
    rows = 20
    cols = 40
    terrain := new(Terrain)
    u := new(Terrain)
    u.SeeAnt(Point{2, 0}, Player(0))
    u.SeeAnt(Point{2, 5}, Player(1))
    u.SeeAnt(Point{2, 6}, Player(1))
    u.SeeHill(Point{2, 6}, Player(1))
    u.SeeHill(Point{3, 6}, Player(1))
    terrain.Update(u)
    if !terrain.At(Point{2, 1}).IsVisible() {
        t.Errorf("visiblility not updated")
    }
    if terrain.At(Point{2, 10}).IsVisible() {
        t.Errorf("too much visiblility")
    }
    if !terrain.At(Point{2, 5}).HasEnemyAnt() {
        t.Errorf("enemy ant")
    }
}
