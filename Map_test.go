package main

import "testing"

func TestFromString(t *testing.T) {
    s := "?.%*\naB2?"
    m := NewMap(s)
    s2 := m.String()
    if s != s2 {
        t.Errorf("%#v != %#v", s, s2)
    }
}

func TestUpdate(t *testing.T) {
    rows = 20
    cols = 40
    m := new(Map)
    u := new(Map)
    u.SeeAnt(Point{2, 0}, Player(0))
    u.SeeAnt(Point{2, 5}, Player(1))
    u.SeeAnt(Point{2, 6}, Player(1))
    u.SeeHill(Point{2, 6}, Player(1))
    u.SeeHill(Point{3, 6}, Player(1))
    m.Update(u)
    if !m.At(Point{2, 1}).IsVisible() {
        t.Errorf("visiblility not updated")
    }
    if m.At(Point{2, 10}).IsVisible() {
        t.Errorf("too much visiblility")
    }
    if !m.At(Point{2, 5}).HasEnemyAnt() {
        t.Errorf("enemy ant")
    }
}
