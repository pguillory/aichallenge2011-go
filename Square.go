package main

const (
    _ = iota
    SQUARE_VISIBLE = byte(1<<iota)
    SQUARE_LAND
    SQUARE_WATER
    SQUARE_FOOD
    SQUARE_ANT
    SQUARE_HILL
)

type Square struct {
    contents byte
    owner Player
}

func (this Square) IsVisible() bool {
    return this.contents & SQUARE_VISIBLE > 0
}

func (this Square) HasLand() bool {
    return this.contents & SQUARE_LAND > 0
}

func (this Square) HasWater() bool {
    return this.contents & SQUARE_WATER > 0
}

func (this Square) HasFood() bool {
    return this.contents & SQUARE_FOOD > 0
}

func (this Square) HasAnt() bool {
    return this.contents & SQUARE_ANT > 0
}

func (this Square) HasHill() bool {
    return this.contents & SQUARE_HILL > 0
}

func (this Square) Owner() Player {
    return this.owner
}

func (this Square) IsFriendly() bool {
    return this.owner == 0
}

func (this Square) IsEnemy() bool {
    return this.owner > 0
}

func (this Square) HasFriendlyAnt() bool {
    return this.IsFriendly() && this.HasAnt()
}

func (this Square) HasFriendlyHill() bool {
    return this.IsFriendly() && this.HasHill()
}

func (this Square) HasEnemyAnt() bool {
    return this.IsEnemy() && this.HasAnt()
}

func (this Square) HasEnemyHill() bool {
    return this.IsEnemy() && this.HasHill()
}


func (this Square) PlusVisible() Square {
    this.contents |= SQUARE_VISIBLE
    return this
}

func (this Square) PlusLand() Square {
    this.contents |= SQUARE_LAND
    return this
}

func (this Square) PlusWater() Square {
    this.contents |= SQUARE_WATER
    return this
}

func (this Square) PlusFood() Square {
    this.contents |= SQUARE_FOOD
    return this
}

func (this Square) PlusAnt(owner Player) Square {
    this.contents |= SQUARE_ANT
    this.owner = owner
    return this
}

func (this Square) PlusHill(owner Player) Square {
    this.contents |= SQUARE_HILL
    this.owner = owner
    return this
}

func (this Square) MinusVisible() Square {
    this.contents &= ^SQUARE_VISIBLE
    return this
}

func (this Square) MinusLand() Square {
    this.contents &= ^SQUARE_LAND
    return this
}

func (this Square) MinusWater() Square {
    this.contents &= ^SQUARE_WATER
    return this
}

func (this Square) MinusFood() Square {
    this.contents &= ^SQUARE_FOOD
    return this
}

func (this Square) MinusAnt() Square {
    this.contents &= ^SQUARE_ANT
    return this
}

func (this Square) MinusHill() Square {
    this.contents &= ^SQUARE_HILL
    return this
}
