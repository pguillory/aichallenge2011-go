package main

type Order struct {
    row, col int
    dir byte
}

type Bot struct {
    m *Map
    u *Map
    s *Scent
}

func (this *Bot) Ready() {
    this.m = new(Map)
    this.s = NewScent(this.m)
}

func (this *Bot) Turn() {
    this.u = new(Map)
}

func (this *Bot) SeeWater(row, col int) {
    this.u.SeeWater(Point{row, col})
}

func (this *Bot) SeeFood(row, col int) {
    this.u.SeeFood(Point{row, col})
}

func (this *Bot) SeeAnt(row, col, player int) {
    this.u.SeeAnt(Point{row, col}, Player(player))
}

func (this *Bot) SeeHill(row, col, player int) {
    this.u.SeeHill(Point{row, col}, Player(player))
}

func (this *Bot) SeeDeadAnt(row, col, player int) {
}

func (this *Bot) Go() (orders []Order) {
    this.m.Update(this.u)
    for i := 0; i < 25; i++ {
        this.s = this.s.Iterate()
    }
    moves := NewMoves(this.m, this.s)
    ForEachPoint(func(p Point) {
        c := moves.At(p).Char()
        if (c == 'N' || c == 'E' || c == 'S' || c == 'W') {
            orders = append(orders, Order{p.row, p.col, c})
        }
    })
    orders = append(orders, Order{0, 0, 'N'})
    return
}
