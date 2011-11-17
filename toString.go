package main

import "bytes"
import "fmt"

const (
    BLACK   = 0
    RED     = 1
    GREEN   = 2
    YELLOW  = 3
    BLUE    = 4
    MAGENTA = 5
    CYAN    = 6
    WHITE   = 7

    HIGH_INTENSITY = 60
    BOLD = 1
    UNDERLINE = 4
)

type ColorChar struct {
    c byte
    foreground, background, style int
}

func ForEachPointAndLine(f1 func(p Point), f2 func()) {
    max_row := 0

    ForEachPoint(func(p Point) {
        for max_row < p.row {
            max_row += 1
            f2()
        }
        f1(p)
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

func GridToColorString(f func(p Point) ColorChar) string {
    b := new(bytes.Buffer)
    var last ColorChar

    ForEachPointAndLine(func(p Point) {
        if p.row < rows - 50 || p.col > 170 {
            return
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
        b.WriteByte(cc.c)
    }, func() {
        b.WriteByte('\n')
    })

    b.WriteString(fmt.Sprintf("%c[0m", 27))

    return b.String()
}
