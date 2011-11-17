package main

import "bytes"

type Map struct {
    squares [MAX_ROWS][MAX_COLS]Square
}

func (this *Map) At(p Point) Square {
    return this.squares[p.row][p.col]
}

func (this *Map) SetAt(p Point, s Square) {
    this.squares[p.row][p.col] = s
}

func (this *Map) SeeWater(p Point) {
    this.SetAt(p, this.At(p).PlusVisible().PlusWater())
}

func (this *Map) SeeLand(p Point) {
    this.SetAt(p, this.At(p).PlusVisible().PlusLand())
}

func (this *Map) SeeFood(p Point) {
    this.SetAt(p, this.At(p).PlusVisible().PlusLand().PlusFood())
}

func (this *Map) SeeAnt(p Point, owner Player) {
    this.SetAt(p, this.At(p).PlusVisible().PlusLand().PlusAnt(owner))
}

func (this *Map) SeeHill(p Point, owner Player) {
    this.SetAt(p, this.At(p).PlusVisible().PlusLand().PlusHill(owner))
}

func (this *Map) Update(m *Map) {
    ForEachPoint(func(p Point) {
        s := this.At(p).MinusVisible().MinusAnt()
        s2 := m.At(p)
        if s2.HasAnt() {
            s = s.PlusAnt(s2.owner)
        }
        this.SetAt(p, s)
    })

    ForEachPoint(func(p Point) {
        if this.At(p).HasFriendlyAnt() {
            ForEachPointWithinRadius2(p, viewradius2, func(p2 Point) {
                this.SetAt(p2, this.At(p2).PlusVisible())
            })
        }
    })

    ForEachPoint(func(p Point) {
        s := this.At(p)
        if s.IsVisible() {
            s2 := m.At(p)

            if s2.HasWater() {
                s = s.PlusWater()
            } else if !s.HasWater() && !s.HasLand() {
                s = s.PlusLand()
            }

            if s2.HasFood() {
                s = s.PlusFood()
            } else if s.HasFood() {
                s = s.MinusFood()
            }

            if s2.HasHill() {
                s = s.PlusHill(s2.owner)
            } else if s.HasHill() {
                s = s.MinusHill()
            }
            this.SetAt(p, s)
        }
    })
}

func (this *Map) String() string {
    b := new(bytes.Buffer)
    max_row := 0

    ForEachPoint(func(p Point) {
        for max_row < p.row {
            max_row += 1
            b.WriteByte('\n')
        }

        s := this.At(p)
        switch {
        case s.HasLand():
            switch {
            case s.HasFood():
                b.WriteByte('*')
            case s.HasAnt() && s.HasHill():
                b.WriteByte('A' + byte(s.owner))
            case s.HasAnt():
                b.WriteByte('a' + byte(s.owner))
            case s.HasHill():
                b.WriteByte('0' + byte(s.owner))
            default:
                b.WriteByte('.')
            }
        case s.HasWater():
            b.WriteByte('%')
        default:
            b.WriteByte('?')
        }
    })

    return b.String()
}

func MapFromString(input string) *Map {
    m := new(Map)

    rows = 0
    cols = 0
    p := Point{0, 0}

    for _, c := range input {
        switch {
        case c == '.':
            m.SeeLand(p)
        case c == '%':
            m.SeeWater(p)
        case c == '*':
            m.SeeFood(p)
        case c >= 'a' && c <= 'z':
            owner := Player(c - 'a')
            m.SeeAnt(p, owner)
        case c >= '0' && c <= '9':
            owner := Player(c - '0')
            m.SeeHill(p, owner)
        case c >= 'A' && c <= 'Z':
            owner := Player(c - 'A')
            m.SeeAnt(p, owner)
            m.SeeHill(p, owner)
        case c == '\n':
            p.row += 1
            p.col = 0
            continue
        }

        if rows <= p.row {
            rows = p.row + 1
        }
        if cols <= p.col {
            cols = p.col + 1
        }
        p.col += 1
    }

    return m
}
