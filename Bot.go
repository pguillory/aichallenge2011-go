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
    scent *Scent
    moves *Moves
}

func (this *Bot) Ready() {
    this.m = new(Map)
    this.mystery = NewMystery(this.m)
    this.scent = NewScent(this.m, this.mystery)
    this.moves = NewMoves(this.m, this.scent)
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
        this.scent.Iterate()
    }
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
        NewLog("time", "log").File().WriteString(fmt.Sprintf("%v\n", timer.String()))
        NewLog("map", "log").TurnFile().WriteString(this.m.String())
    }

    return orders
}
