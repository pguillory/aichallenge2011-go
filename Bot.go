package main

import "fmt"
import "os"

type Bot struct {
    terrain *Terrain
    update *Terrain
    holyGround *HolyGround
    mystery *Mystery
    forageScent, battleScent *Scent
    army *Army
    predictions *Predictions
    command *Command
    hud *os.File
}

func (this *Bot) Ready() {
    VerifySituationSize()

    this.terrain = new(Terrain)
    this.holyGround = NewHolyGround(this.terrain)
    this.mystery = NewMystery(this.terrain)
    this.forageScent = NewForageScent(this.terrain, this.holyGround, this.mystery)
    this.battleScent = NewBattleScent(this.terrain, this.holyGround, this.mystery)
    this.army = NewArmy(this.terrain)
    this.predictions = NewPredictions(this.terrain)
    this.command = NewCommand(this.terrain, this.forageScent, this.battleScent, this.army, this.predictions)

    this.hud = NewLog("hud", "txt")
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
    timeLog := NewTurnLog("times", "txt")

    timer.Start("map")
    this.terrain.Update(this.update)
    timer.Stop()
    timeLog.WriteString(fmt.Sprintf("map: %v\n", timer.times["map"]))

    timer.Start("holyGround")
    this.holyGround.Calculate()
    timer.Stop()
    timeLog.WriteString(fmt.Sprintf("holyGround: %v\n", timer.times["holyGround"]))

    timer.Start("mystery")
    this.mystery.Calculate()
    timer.Stop()
    timeLog.WriteString(fmt.Sprintf("mystery: %v\n", timer.times["mystery"]))

    timer.Start("forageScent")
    this.forageScent.Calculate()
    timer.Stop()
    timeLog.WriteString(fmt.Sprintf("forageScent: %v\n", timer.times["forageScent"]))

    timer.Start("battleScent")
    this.battleScent.Calculate()
    timer.Stop()
    timeLog.WriteString(fmt.Sprintf("battleScent: %v\n", timer.times["battleScent"]))

    timer.Start("army")
    this.army.Calculate()
    timer.Stop()
    timeLog.WriteString(fmt.Sprintf("army: %v\n", timer.times["army"]))

    timer.Start("predictions")
    this.predictions.Calculate()
    timer.Stop()
    timeLog.WriteString(fmt.Sprintf("predictions: %v\n", timer.times["predictions"]))

    timer.Start("command")
    this.command.Calculate()
    timer.Stop()
    timeLog.WriteString(fmt.Sprintf("command: %v\n", timer.times["command"]))

    timer.Start("issueOrder")
    this.command.ForEach(func(move Move) {
        issueOrder(move.from.row, move.from.col, move.dir.Char())
    })
    timer.Stop()
    timeLog.WriteString(fmt.Sprintf("issueOrder: %v\n", timer.times["issueOrder"]))

    // TODO: do this in a goroutine
    this.hud.WriteString(fmt.Sprintf("%v\n", this.ColorString()))
    this.hud.WriteString(fmt.Sprintf("turn %v, times %v\n", turn, timer.String()))
    //NewTurnLog("map", "txt").WriteString(this.terrain.String())
    //NewTurnLog("mystery", "txt").WriteString(this.mystery.String())
    //NewTurnLog("forageScent", "txt").WriteString(this.forageScent.String())
    //NewTurnLog("forageScent", "csv").WriteString(this.forageScent.Csv())
    //NewTurnLog("battleScent", "txt").WriteString(this.battleScent.String())
    //NewTurnLog("battleScent", "csv").WriteString(this.battleScent.Csv())
    //NewTurnLog("army", "txt").WriteString(this.army.String())
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
                    return ColorChar{'a' + byte(s.owner), BLACK, HIGH_INTENSITY + RED, style}
                }
            case s.HasAnt():
                if s.IsFriendly() {
                    return ColorChar{'a' + byte(s.owner), GREEN, background, style}
                } else {
                    return ColorChar{'a' + byte(s.owner), RED, background, style}
                }
            case s.HasHill():
                if s.IsFriendly() {
                    return ColorChar{' ', BLACK, HIGH_INTENSITY + GREEN, style}
                } else {
                    return ColorChar{' ', BLACK, HIGH_INTENSITY + RED, style}
                }
            }
            return ColorChar{'.', HIGH_INTENSITY + BLACK, background, style}
        case s.HasWater():
            return ColorChar{'%', BLUE, background, style}
        }
        return ColorChar{'?', WHITE, background, style}
    })
}
