package main

import "runtime"
import "os"
import "bufio"
import "strings"
import "strconv"
import "rand"
import "fmt"
import "flag"
//import "os/signal"
//import "runtime/debug"

func main() {
    flag.BoolVar(&debug, "debug", false, "debug mode")
    flag.Parse()

    //go func() {
    //    //log := NewLog("signals", "txt")
    //    for sig := range signal.Incoming {
    //        //log.WriteString(fmt.Sprintf("Caught signal %v\n%#v\n\n", sig, sig))
    //        panic(fmt.Sprintf("Caught signal %v\n%#v\n\n", sig, sig))
    //        //if s == os.SIGINT
    //        //os.Exit(1)
    //    }
    //}()

	var bot Bot

    stdin := bufio.NewReader(os.Stdin)
    for {
        bytes, _, _ := stdin.ReadLine()
        words := strings.Fields(string(bytes))
        switch len(words) {
        case 1:
            switch words[0] {
            case "go":
                bot.Go(func(row, col int, dir byte) {
                    fmt.Printf("o %v %v %c\n", row, col, dir)
                }, func() {
                    fmt.Println("go")
                    os.Stdout.Sync()
                	runtime.GC()
                })
            case "ready":
                bot.Ready()
                fmt.Println("go")
                os.Stdout.Sync()
            case "end":
                bot.End()
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
            case "player_seed":
    			player_seed, _ := strconv.Atoi64(words[1])
    			rand.Seed(player_seed)
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
