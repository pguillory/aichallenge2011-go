package main

import "runtime"
import "os"
import "bufio"
import "strings"
import "strconv"
import "fmt"

type Order struct {
    row, col int
    dir byte
}

var debugMode = true

func main() {
	var bot Bot

    stdin := bufio.NewReader(os.Stdin)
    for {
        bytes, _, _ := stdin.ReadLine()
        words := strings.Fields(string(bytes))
        switch len(words) {
        case 1:
            switch words[0] {
            case "go":
                for _, order := range bot.Go() {
                    fmt.Printf("o %v %v %c\n", order.row, order.col, order.dir)
                }
                fmt.Println("go")
                os.Stdout.Sync()
            	runtime.GC()
            case "ready":
                bot.Ready()
                fmt.Println("go")
                os.Stdout.Sync()
            case "end":
                os.Exit(0)
            }
        case 2:
            switch words[0] {
            case "turn":
                turn, _ = strconv.Atoi(words[1])
                if turn > 0 {
                    bot.Turn()
                }
            case "loadtime":
                loadtime, _ = strconv.Atoi(words[1])
            case "turntime":
                turntime, _ = strconv.Atoi(words[1])
            case "rows":
                rows, _ = strconv.Atoi(words[1])
            case "cols":
                cols, _ = strconv.Atoi(words[1])
            case "turns":
                turns, _ = strconv.Atoi(words[1])
            case "viewradius2":
                viewradius2, _ = strconv.Atoi(words[1])
            case "attackradius2":
                attackradius2, _ = strconv.Atoi(words[1])
            case "spawnradius2":
                spawnradius2, _ = strconv.Atoi(words[1])
            }
        case 3:
            switch words[0] {
            case "w":
                row, _ := strconv.Atoi(words[1])
                col, _ := strconv.Atoi(words[2])
                bot.SeeWater(row, col)
            case "f":
                row, _ := strconv.Atoi(words[1])
                col, _ := strconv.Atoi(words[2])
                bot.SeeFood(row, col)
            }
        case 4:
            switch words[0] {
            case "a":
                row, _ := strconv.Atoi(words[1])
                col, _ := strconv.Atoi(words[2])
                player, _ := strconv.Atoi(words[3])
                bot.SeeAnt(row, col, player)
            case "h":
                row, _ := strconv.Atoi(words[1])
                col, _ := strconv.Atoi(words[2])
                player, _ := strconv.Atoi(words[3])
                bot.SeeHill(row, col, player)
            case "d":
                row, _ := strconv.Atoi(words[1])
                col, _ := strconv.Atoi(words[2])
                player, _ := strconv.Atoi(words[3])
                bot.SeeDeadAnt(row, col, player)
            }
        }
    }
}
