package main

type Player byte

const MAX_PLAYERS = 10

func ForEachEnemyPlayer(f func(Player)) {
    for player := Player(1); player < MAX_PLAYERS; player++ {
        f(player)
    }
}
