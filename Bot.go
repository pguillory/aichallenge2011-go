package main

import "fmt"
import "os"

type Bot struct {
    terrain, update *Terrain
    mystery *Mystery
    potentialEnemy *PotentialEnemy
    scrum *Scrum
    distanceToFood, distanceToTrouble, distanceToDoom *TravelDistance
    army *Army
    predictions *Predictions
    command *Command
    hud *os.File
    hudCenter Point
}

func (this *Bot) Ready() {
    VerifySituationSize()

    this.terrain = new(Terrain)
    this.mystery = NewMystery(this.terrain)
    this.potentialEnemy = NewPotentialEnemy(this.terrain)
    this.army = NewArmy(this.terrain)
    this.predictions = NewPredictions(this.terrain)
    this.scrum = NewScrum()
    this.distanceToFood = DistanceToFood(this.terrain)
    this.distanceToTrouble = DistanceToTrouble(this.terrain, this.mystery, this.potentialEnemy, this.scrum)
    this.distanceToDoom = DistanceToDoom(this.terrain, this.mystery, this.potentialEnemy, this.scrum)
    this.command = NewCommand(this.terrain, this.army, this.predictions, this.scrum, this.distanceToFood, this.distanceToTrouble, this.distanceToDoom)

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
    this.army.Calculate()
    this.predictions.Calculate()
/*
    this.distanceToFood.Calculate()
    this.distanceToTrouble.Calculate()
    this.distanceToDoom.Calculate()
*/
    this.command.Calculate()

    this.command.ForEach(func(move Move) {
        issueOrder(move.from.row, move.from.col, move.dir.Char())
    })
    done()

    // TODO: do this in a goroutine
    this.hud.WriteString(fmt.Sprintf("\n%v\n", this.ColorString()))
    this.hud.WriteString(fmt.Sprintf("turn %v, times: map %v, myst %v, potE %v, army %v, pred %v, dF %v, dT %v, dD %v, comm %v", turn,
    this.terrain.time, this.mystery.time, this.potentialEnemy.time, this.army.time, this.predictions.time, this.distanceToFood.time, this.distanceToTrouble.time, this.distanceToDoom.time, this.command.time))
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
    ForEachPoint(func(p Point) {
        if this.terrain.At(p).HasFriendlyHill() {
            this.hudCenter = p
        }
    })

    topLeftCorner := this.hudCenter.Plus(Point{-25, -85})
    if cols < 170 {
        topLeftCorner.col -= (cols - 170) / 2
    }

    return GridToColorString(func(p1 Point) ColorChar {
        p := p1.Plus(topLeftCorner)
        s := this.terrain.At(p)

        background := BLACK + HIGH_INTENSITY
        if s.IsVisible() {
            background -= HIGH_INTENSITY
        }

        style := 0

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
                    if this.army.IsBerzerkerAt(p) {
                        return ColorChar{'a' + byte(s.owner), MAGENTA, background, style}
                    //} else if this.army.IsSoldierAt(p) {
                    //    return ColorChar{'a' + byte(s.owner), CYAN, background, style}
                    } else {
                        return ColorChar{'a' + byte(s.owner), GREEN, background, style}
                    }
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
            if this.potentialEnemy.At(p) {
                return ColorChar{'.', RED, background, style}
            } else {
                return ColorChar{'.', HIGH_INTENSITY + BLACK, background, style}
            }
        case s.HasWater():
            return ColorChar{'%', BLUE, background, style}
        }
        return ColorChar{'?', WHITE, background, style}
    })
}
