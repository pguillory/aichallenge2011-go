package main

import "fmt"

type Order struct {
    row, col int
    dir byte
}

type Bot struct {
    m *Map
    update *Map
    mystery *Mystery
    workerScent, soldierScent *Scent
    army *Army
    moves *Moves
}

func (this *Bot) Ready() {
    this.m = new(Map)
    this.mystery = NewMystery(this.m)
    this.workerScent = NewScent(this.m, this.mystery)
    this.soldierScent = NewScent(this.m, this.mystery)
    this.army = NewArmy(this.m)
    this.moves = NewMoves(this.m, this.workerScent, this.soldierScent, this.army)
}

func (this *Bot) Turn() {
    this.update = new(Map)
}

func (this *Bot) SeeWater(row, col int) {
    this.update.SeeWater(Point{row, col})
}

func (this *Bot) SeeFood(row, col int) {
    this.update.SeeFood(Point{row, col})
}

func (this *Bot) SeeAnt(row, col, player int) {
    this.update.SeeAnt(Point{row, col}, Player(player))
}

func (this *Bot) SeeHill(row, col, player int) {
    this.update.SeeHill(Point{row, col}, Player(player))
}

func (this *Bot) SeeDeadAnt(row, col, player int) {
}

func (this *Bot) Go() []Order {
    timer := NewTimer()

    timer.Start("map")
    this.m.Update(this.update)
    timer.Stop()

    timer.Start("mystery")
    this.mystery.Iterate()
    timer.Stop()

    timer.Start("scent")
    for i := 0; i < 10; i++ {
        this.workerScent.Iterate()
        this.soldierScent.IterateSoldier()
    }
    timer.Stop()

    timer.Start("army")
    this.army.Iterate()
    timer.Stop()

    timer.Start("moves")
    this.moves.Calculate()
    timer.Stop()

    timer.Start("bot")
    orders := make([]Order, 0)
    ForEachPoint(func(p Point) {
        c := this.moves.At(p).Char()
        if (c == 'N' || c == 'E' || c == 'S' || c == 'W') {
            orders = append(orders, Order{p.row, p.col, c})
        }
    })
    timer.Stop()

    if debugMode {
        hud := NewLog("hud", "log").File()
        hud.WriteString(fmt.Sprintf("%v\nturn %v, times %v\n", this.ColorString(), turn, timer.String()))
        NewLog("map", "log").TurnFile().WriteString(this.m.String())
        NewLog("army", "log").TurnFile().WriteString(this.army.String())
    }

    return orders
}

func (this *Bot) ColorString() string {
    return GridToColorString(func(p Point) ColorChar {
        s := this.m.At(p)

        background := BLACK + HIGH_INTENSITY
        if s.IsVisible() {
            background -= HIGH_INTENSITY
        }

        style := 0
        if this.army.IsSoldierAt(p) {
            style = UNDERLINE
        }

        switch {
        case s.HasLand():
            switch {
            case s.HasFood():
                return ColorChar{'*', YELLOW, background, style}
            case s.HasAnt() && s.HasHill():
                if s.IsFriendly() {
                    return ColorChar{'A' + byte(s.owner), HIGH_INTENSITY + GREEN, background, style}
                } else {
                    return ColorChar{'A' + byte(s.owner), RED, background, style}
                }
            case s.HasAnt():
                if s.IsFriendly() {
                    return ColorChar{'a' + byte(s.owner), GREEN, background, style}
                } else {
                    return ColorChar{'a' + byte(s.owner), RED, background, style}
                }
            case s.HasHill():
                if s.IsFriendly() {
                    return ColorChar{'0' + byte(s.owner), HIGH_INTENSITY + GREEN, background, style}
                } else {
                    return ColorChar{'0' + byte(s.owner), RED, background, style}
                }
            }
            return ColorChar{'.', HIGH_INTENSITY + BLACK, background, style}
        case s.HasWater():
            return ColorChar{'%', BLUE, background, style}
        }
        return ColorChar{'?', WHITE, background, style}
    })
}
