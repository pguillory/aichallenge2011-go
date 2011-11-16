/*
package main

import "bytes"

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
    var s Square

    for _, c := range input {
        switch {
        case c == '.':
            m.SetAt(p, s.PlusVisible().PlusLand())
        case c == '%':
            m.SetAt(p, s.PlusVisible().PlusWater())
        case c == '*':
            m.SetAt(p, s.PlusVisible().PlusLand().PlusFood())
        case c >= 'a' && c <= 'z':
            owner := Player(c - 'a')
            m.SetAt(p, s.PlusVisible().PlusLand().PlusAnt(owner))
        case c >= '0' && c <= '9':
            owner := Player(c - '0')
            m.SetAt(p, s.PlusVisible().PlusLand().PlusHill(owner))
        case c >= 'A' && c <= 'Z':
            owner := Player(c - 'A')
            m.SetAt(p, s.PlusVisible().PlusLand().PlusAnt(owner).PlusHill(owner))
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
*/
