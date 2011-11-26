package main

import "fmt"
import "os"

type Bot struct {
    terrain, update *Terrain
    mystery *Mystery
    potentialEnemy *PotentialEnemy
    distanceToTrouble *TravelDistance
    army *Army
    predictions *Predictions
    command *Command
    hud *os.File
}

func (this *Bot) Ready() {
    VerifySituationSize()

    this.terrain = new(Terrain)
    //this.distanceToEnemy = DistanceToEnemy(this.terrain)
    //this.distanceToFriendlyHill = DistanceToFriendlyHill(this.terrain)
    //this.mystery = NewMystery(this.terrain)
    //this.forageScent = NewForageScent(this.terrain, this.distanceToEnemy, this.distanceToFriendlyHill, this.mystery)
    //this.battleScent = NewBattleScent(this.terrain, this.distanceToEnemy, this.distanceToFriendlyHill, this.mystery)
    this.mystery = NewMystery(this.terrain)
    this.potentialEnemy = NewPotentialEnemy(this.terrain)
    this.distanceToTrouble = DistanceToTrouble(this.terrain, this.mystery, this.potentialEnemy)
    this.army = NewArmy(this.terrain)
    this.predictions = NewPredictions(this.terrain)
    this.command = NewCommand(this.terrain, this.army, this.predictions, this.distanceToTrouble)

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

func (this *Bot) Go(issueOrder func(int, int, byte), done func()) {
    this.terrain.Update(this.update)

    this.mystery.Calculate()
    this.potentialEnemy.Calculate()
    this.command.Calculate()

    this.command.ForEach(func(move Move) {
        issueOrder(move.from.row, move.from.col, move.dir.Char())
    })
    done()

    // TODO: do this in a goroutine
    this.hud.WriteString(fmt.Sprintf("%v\n", this.ColorString()))
    //this.hud.WriteString(fmt.Sprintf("turn %v, times: map %v, dH %v, dE %v, myst %v, for %v, army %v, pred %v, comm %v\n", turn, this.terrain.time, this.distanceToFriendlyHill.time, this.distanceToEnemy.time, this.mystery.time, this.forageScent.time, this.army.time, this.predictions.time, this.command.time))
    //NewTurnLog("map", "txt").WriteString(this.terrain.String())
    //NewTurnLog("mystery", "txt").WriteString(this.mystery.String())
    //NewTurnLog("potentialEnemy", "txt").WriteString(this.potentialEnemy.String())
    //NewTurnLog("forageScent", "txt").WriteString(this.forageScent.String())
    //NewTurnLog("forageScent", "csv").WriteString(this.forageScent.Csv())
    //NewTurnLog("battleScent", "txt").WriteString(this.battleScent.String())
    //NewTurnLog("battleScent", "csv").WriteString(this.battleScent.Csv())
    //NewTurnLog("army", "txt").WriteString(this.army.String())
    //NewTurnLog("distanceToTrouble", "txt").WriteString(this.distanceToTrouble.String())
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
