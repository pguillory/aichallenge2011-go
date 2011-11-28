package main

import "fmt"
import "os"

type Bot struct {
    terrain, update *Terrain
    mystery *Mystery
    potentialEnemy *PotentialEnemy
    army *Army
    predictions *Predictions
    distanceToFood, distanceToTrouble, distanceToDoom *TravelDistance
    rageVirus *RageVirus
    command *Command
    hud *os.File
    hudCenter Point
}

func (this *Bot) Ready() {
    //PrepareRadius2Tables(attackradius2)
    //PrepareRadius2Tables(viewradius2)
    //PrepareRadius2Tables(SITUATION_RADIUS2)
    VerifySituationSize()

    this.terrain = new(Terrain)
    this.mystery = NewMystery(this.terrain)
    this.potentialEnemy = NewPotentialEnemy(this.terrain)
    this.army = NewArmy(this.terrain)
    this.predictions = NewPredictions(this.terrain)
    this.distanceToFood = DistanceToFood(this.terrain)
    this.distanceToTrouble = DistanceToTrouble(this.terrain, this.mystery, this.potentialEnemy)
    this.distanceToDoom = DistanceToDoom(this.terrain, this.mystery, this.potentialEnemy)
    this.rageVirus = NewRageVirus(this.terrain, this.army, this.distanceToTrouble)
    this.command = NewCommand(this.terrain, this.army, this.predictions, this.distanceToFood, this.distanceToTrouble, this.distanceToDoom, this.rageVirus)

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
    startTime := now()

    this.terrain.Update(this.update)

    this.mystery.Calculate()
    this.potentialEnemy.Calculate()
    this.army.Calculate()
    this.predictions.Calculate()
    this.distanceToFood.Calculate()
    this.distanceToTrouble.Calculate()
    this.distanceToDoom.Calculate()
    this.rageVirus.Calculate()
    this.command.Calculate()

    this.command.ForEach(func(move Move) {
        issueOrder(move.from.row, move.from.col, move.dir.Char())
    })
    done()

    time := now() - startTime

    this.hud.WriteString(fmt.Sprintf("\n%v\n", this.ColorString()))
    this.hud.WriteString(fmt.Sprintf("turn %v, time %v (map %v, myst %v, potE %v, army %v, pred %v, dF %v, dT %v, dD %v, comm %v)",
            turn,
            time,
            this.terrain.time,
            this.mystery.time,
            this.potentialEnemy.time,
            this.army.time,
            this.predictions.time,
            this.distanceToFood.time, this.distanceToTrouble.time, this.distanceToDoom.time,
            this.command.time))

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

func (this *Bot) End() {
    this.hud.WriteString(fmt.Sprintf("\nGame over.\n"))
    //for i, count := range counts {
    //    if count > 0 {
    //        fmt.Println("count", i, count)
    //    }
    //}
}

func (this *Bot) ColorString() string {
    ForEachPoint(func(p Point) {
        if this.terrain.At(p).HasFriendlyHill() {
            this.hudCenter = p
        }
    })

    topLeftCorner := this.hudCenter.Plus(Point{-31, -85})
    if cols < 170 {
        topLeftCorner.col -= (cols - 170) / 2
    }

    return GridToColorString(func(p1 Point) ColorChar {
        p := p1.Plus(topLeftCorner)
        s := this.terrain.At(p)

        var cc ColorChar

        cc.background = BLACK
        if !s.IsVisible() {
            cc.background += BRIGHT
        }

        switch {
        case s.HasLand():
            switch {
            case s.HasFood():
                cc.symbol = "*"
                cc.foreground = YELLOW
            case s.HasHill():
                if s.HasAnt() {
                    cc.symbol = string('a' + byte(s.owner))
                } else {
                    cc.symbol = " "
                }
                cc.foreground = BLACK
                if s.IsFriendly() {
                    cc.background = BRIGHT + GREEN
                } else {
                    cc.background = BRIGHT + RED
                }
            case s.HasAnt():
                cc.symbol = string('a' + byte(s.owner))
                if s.IsFriendly() {
                    if this.army.IsBerzerkerAt(p) {
                        cc.foreground = MAGENTA
                    } else {
                        cc.foreground = GREEN
                    }
                    if this.rageVirus.InfectedAt(p) {
                        cc.foreground += BRIGHT
                    }
                } else {
                    cc.foreground = RED
                }
            default:
                cc.symbol = "."
                if this.potentialEnemy.At(p) {
                    cc.foreground = RED
                } else {
                    cc.foreground = BRIGHT + BLACK
                }
            }
        case s.HasWater():
            cc.symbol = "▒"
            cc.foreground = BLUE
            // ■ ▓ █
        default:
            cc.symbol = "?"
            cc.foreground = WHITE
        }
        return cc
    })
}
