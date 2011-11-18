package main

import "fmt"
import "os"

type Bot struct {
    terrain *Terrain
    update *Terrain
    mystery *Mystery
    workerScent, soldierScent *Scent
    army *Army
    command *Command
    hud *os.File
}

func (this *Bot) Ready() {
    this.terrain = new(Terrain)
    this.mystery = NewMystery(this.terrain)
    this.workerScent = NewScent(this.terrain, this.mystery)
    this.soldierScent = NewScent(this.terrain, this.mystery)
    this.army = NewArmy(this.terrain)
    this.command = NewCommand(this.terrain, this.workerScent, this.soldierScent, this.army)
    if debugMode {
        this.hud = NewLog("hud", "log")
    }
}

func (this *Bot) Turn() {
    this.update = new(Terrain)
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

func (this *Bot) Go(issueOrder func(int, int, byte)) {
    timer := NewTimer()

    timer.Start("map")
    this.terrain.Update(this.update)
    timer.Stop()

    timer.Start("mystery")
    this.mystery.Calculate()
    timer.Stop()

    timer.Start("scent")
    for i := 0; i < 25; i++ {
        this.workerScent.Calculate()
        this.soldierScent.CalculateSoldier()
    }
    timer.Stop()

    timer.Start("army")
    this.army.Calculate()
    timer.Stop()

    timer.Start("command")
    this.command.Calculate()
    timer.Stop()

    timer.Start("issueOrder")
    this.command.ForEach(func(move Move) {
        issueOrder(move.from.row, move.from.col, move.dir.Char())
    })
    timer.Stop()

    if debugMode {
        this.hud.WriteString(fmt.Sprintf("%v\n", this.ColorString()))
        this.hud.WriteString(fmt.Sprintf("turn %v, times %v\n", turn, timer.String()))
        NewTurnLog("map", "log").WriteString(this.terrain.String())
        NewTurnLog("army", "log").WriteString(this.army.String())
    }
}

func (this *Bot) ColorString() string {
    return GridToColorString(func(p Point) ColorChar {
        s := this.terrain.At(p)

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
                    return ColorChar{'a' + byte(s.owner), BLACK, HIGH_INTENSITY + GREEN, style}
                } else {
                    return ColorChar{'a' + byte(s.owner), RED, background, style}
                }
            case s.HasAnt():
                if s.IsFriendly() {
                    return ColorChar{'a' + byte(s.owner), GREEN, background, style}
                } else {
                    return ColorChar{'a' + byte(s.owner), RED, background, style}
                }
            case s.HasHill():
                if s.IsFriendly() {
                    return ColorChar{' ' + byte(s.owner), BLACK, HIGH_INTENSITY + GREEN, style}
                } else {
                    return ColorChar{' ' + byte(s.owner), RED, background, style}
                }
            }
            return ColorChar{'.', HIGH_INTENSITY + BLACK, background, style}
        case s.HasWater():
            return ColorChar{'%', BLUE, background, style}
        }
        return ColorChar{'?', WHITE, background, style}
    })
}
