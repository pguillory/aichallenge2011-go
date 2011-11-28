package main

import "bytes"
import "fmt"
import "math"

const (
    BLACK   = 0
    RED     = 1
    GREEN   = 2
    YELLOW  = 3
    BLUE    = 4
    MAGENTA = 5
    CYAN    = 6
    WHITE   = 7

    BRIGHT = 60
    BOLD = 1
    UNDERLINE = 4
)

type ColorChar struct {
    symbol string
    foreground, background int
}

func ForEachPointAndLine(pointFunc func(p Point), lineFunc func()) {
    max_row := 0

    ForEachPoint(func(p Point) {
        for max_row < p.row {
            max_row += 1
            lineFunc()
        }
        pointFunc(p)
    })
}

func ForEachPointAndLineRadius2(origin Point, radius2 int, insideFunc func(p Point), outsideFunc func(), lineFunc func()) {
    radius := int(math.Floor(math.Sqrt(float64(radius2))))

    max_row := normalizeRow(origin.row - radius)

    //fmt.Println(radius2, radius, max_row)

    ForEachPointWithinManhattanDistance(origin, radius, func(p Point) {
        //fmt.Println(p.row)
        if max_row > p.row {
            max_row = p.row
            lineFunc()
        }
        for max_row < p.row {
            max_row += 1
            //fmt.Println(" => ", max_row)
            lineFunc()
        }
        if p.Distance2(origin) <= radius2 {
            insideFunc(p)
        } else {
            outsideFunc()
        }
    })
}

func GridToString(f func(p Point) byte) string {
    b := new(bytes.Buffer)

    ForEachPointAndLine(func(p Point) {
        b.WriteByte(f(p))
    }, func() {
        b.WriteByte('\n')
    })

    return b.String()
}

func GridToColorString(viewportSize Point, f func(p Point) ColorChar) string {
    b := new(bytes.Buffer)
    var last ColorChar

    max_row := 0

    ForEachPoint(func(p Point) {
        if p.row > viewportSize.row || p.col > viewportSize.col {
            return
        }

        if max_row < p.row {
            max_row = p.row
            b.WriteByte('\n')
        }

        cc := f(p)
        cc.foreground += 30
        cc.background += 40

        switch {
        //case last.style != cc.style:
        //    b.WriteString(fmt.Sprintf("%c[%v;%v;%vm", 27, cc.style, cc.foreground, cc.background))
        case last.foreground != cc.foreground && last.background != cc.background:
            b.WriteString(fmt.Sprintf("%c[%v;%vm", 27, cc.foreground, cc.background))
        case last.foreground != cc.foreground:
            b.WriteString(fmt.Sprintf("%c[%vm", 27, cc.foreground))
        case last.background != cc.background:
            b.WriteString(fmt.Sprintf("%c[%vm", 27, cc.background))
        }
        last = cc
        //b.WriteByte(cc.c)
        //b.WriteString(fmt.Sprintf("%c%c%c", 0xe2, 0x89, 0xa1))
        b.WriteString(cc.symbol)
    })

    b.WriteString(fmt.Sprintf("%c[0m", 27))

    return b.String()
}

func GridToCsv(f func(p Point) string) string {
    b := new(bytes.Buffer)

    ForEachPointAndLine(func(p Point) {
        if p.col > 0 {
            b.WriteByte(',')
        }
        b.WriteString(f(p))
    }, func() {
        b.WriteByte('\n')
    })

    return b.String()
}

func StringFromRadius2(p Point, radius2 int, f func(p Point) byte) string {
    b := new(bytes.Buffer)

    ForEachPointAndLineRadius2(p, radius2, func(p Point) {
        b.WriteByte(f(p))
    }, func() {
        b.WriteByte(' ')
    }, func() {
        b.WriteByte('\n')
    })

    return b.String()
}
