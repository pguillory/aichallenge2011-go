package main

import "fmt"

type Observation struct {
    situation *Situation
    dir Direction
    timesSeen uint32
}

func NewObservation(situation *Situation, dir Direction) *Observation {
    this := new(Observation)
    this.situation = situation
    this.dir = dir
    this.timesSeen = 1
    return this
}

func (this *Observation) Confident() bool {
    return (this.timesSeen >= 2 && this.dir.IsSingle())
}

func (this *Observation) String() string {
    return fmt.Sprintf("%v, seen %v times\n%v\n", this.dir, this.timesSeen, this.situation)
}



type ObservationLinkedListNode struct {
    observation *Observation
    next *ObservationLinkedListNode
}

func NewObservationLinkedListNode(observation *Observation, next *ObservationLinkedListNode) *ObservationLinkedListNode {
    this := new(ObservationLinkedListNode)
    this.observation = observation
    this.next = next
    return this
}

type ObservationLinkedList struct {
    head *ObservationLinkedListNode
    length int
}

func (this *ObservationLinkedList) ForEach(f func(*Observation)) {
    for n := this.head; n != nil; n = n.next {
        f(n.observation)
    }
}

func (this *ObservationLinkedList) Add(situation *Situation, dir Direction) {
    var previous *ObservationLinkedListNode
    for n := this.head; n != nil; n = n.next {
        if n.observation.situation.Matches(situation) {
            n.observation.dir |= dir
            n.observation.timesSeen += 1
            if previous != nil {
                previous.next = n.next
                n.next = this.head
                this.head = n
            }
            return
        }
        previous = n
    }

    this.head = NewObservationLinkedListNode(NewObservation(situation, dir), this.head)
    this.length += 1
}

func (this *ObservationLinkedList) Matching(situation *Situation) (result *Observation, found bool) {
    for n := this.head; n != nil; n = n.next {
        if n.observation.situation.Matches(situation) {
            return n.observation, true
        }
    }
    return nil, false
}



const OBSERVATIONSLAB_SIZE = 4096
//const OBSERVATIONSLAB_BITSIZE = (OBSERVATIONSLAB_SIZE + 31) / 32

type ObservationSlab struct {
    members [OBSERVATIONSLAB_SIZE]ObservationLinkedListNode
    head *ObservationLinkedListNode
    length int
}

func (this *ObservationSlab) ForEach(f func(*Observation)) {
    for n := this.head; n != nil; n = n.next {
        f(n.observation)
    }
}

func (this *ObservationSlab) Add(situation *Situation, dir Direction) {
    var last, last2 *ObservationLinkedListNode

    for n := this.head; n != nil; n = n.next {
        if n.observation.situation.Matches(situation) {
            n.observation.dir |= dir
            n.observation.timesSeen += 1
            if last != nil {
                last.next = n.next
                n.next = this.head
                this.head = n
            }
            return
        }
        last2 = last
        last = n
    }

    if this.length < OBSERVATIONSLAB_SIZE {
        last = &this.members[this.length]
        this.length += 1
    } else {
        last2.next = nil
    }

    last.observation = NewObservation(situation, dir)
    last.next = this.head
    this.head = last
}

func (this *ObservationSlab) Matching(situation *Situation) (result *Observation, found bool) {
    for n := this.head; n != nil; n = n.next {
        if n.observation.situation.Matches(situation) {
            return n.observation, true
        }
    }
    return nil, false
}




type Predictions struct {
    time int64
    turn int
    terrain *Terrain
    oldTerrain Terrain
    observations [MAX_PLAYERS]ObservationSlab
}

func NewPredictions(terrain *Terrain) *Predictions {
    this := new(Predictions)
    this.terrain = terrain

    this.Calculate()
    return this
}

func (this *Predictions) Calculate() {
    if this.turn == turn {
        return
    }
    startTime := now()

    //log := NewTurnLog("predictions", "txt")

    ForEachEnemyPlayer(func(player Player) {
        //log.WriteString(fmt.Sprintf("player %v\n", player))

        found := false

        moves := new(MoveSet)
        ForEachPoint(func(p Point) {
            if this.oldTerrain.At(p).HasAntBelongingTo(player) {
                ForEachDirection(func(dir Direction) {
                    if !(dir.Includes(EAST | SOUTH) && this.terrain.At(p).HasAnt()) {
                        p2 := p.Neighbor(dir)
                        if this.terrain.At(p2).HasAntBelongingTo(player) {
                            found = true
                            moves.Include(Move{p, dir})
                        }
                    }
                })
            }
        })

        if found {
            //log.WriteString(fmt.Sprintf("before elimination, %v moves\n", moves.Cardinality()))

            //moves.EliminateLoops()

            moves.ForEach(func(move Move) {
                if move.dir.IsSingle() {
                    moves.Select(move)
                }
            })

            //log.WriteString(fmt.Sprintf("after elimination, %v moves\n", moves.Cardinality()))

            moves.ForEach(func(move Move) {
                if move.dir.IsSingle() {
                    //log.WriteString(fmt.Sprintf("move: %v\n", move))

                    situation, friendlyNearby := NewSituation(&this.oldTerrain, move.from)

                    //log.WriteString(fmt.Sprintf("%v\n", situation))

                    if friendlyNearby {
                        this.observations[player].Add(situation, move.dir)

                        //log.WriteString(fmt.Sprintf("added\n\n"))
                    } else {
                        //log.WriteString(fmt.Sprintf("not relevant\n\n"))
                    }
                }
            })

            //log.WriteString(fmt.Sprintf("player %v has %v observations\n", player, this.observations[player].length))
            //this.observations[player].ForEach(func(observation *Observation) {
            //    log.WriteString(fmt.Sprintf("%v\n", observation))
            //})
        }
    })

    //log.WriteString(fmt.Sprintf("saving terrain\n"))
    this.oldTerrain = *this.terrain

    //log.WriteString(fmt.Sprintf("done\n"))

    this.time = now() - startTime
    this.turn = turn
}

func (this *Predictions) At(p Point) Direction {
    //log := NewTurnLog("Predictions.At", "txt")

    //log.WriteString(fmt.Sprintf("point %v\n", p))

    situation, friendlyNearby := NewSituation(this.terrain, p)

    if friendlyNearby {
        //log.WriteString(fmt.Sprintf("%v\n", situation))

        player := this.terrain.At(p).owner
        //log.WriteString(fmt.Sprintf("player %v has %v observations\n", player, this.observations[player].length))

        observation, found := this.observations[player].Matching(situation)
        if found {
            //log.WriteString(fmt.Sprintf("found\n"))

            if observation.Confident() {
                //log.WriteString(fmt.Sprintf("return %v\n", observation.dir))

                return observation.dir
            } else {
                //log.WriteString(fmt.Sprintf("ignoring, timesSeen=%v, dir=%v\n", observation.timesSeen, observation.dir))
            }
        } else {
            //log.WriteString(fmt.Sprintf("not found\n"))
        }
    } else {
        //log.WriteString(fmt.Sprintf("not relevant\n"))
    }

    //log.WriteString(fmt.Sprintf("return ALL_DIRECTIONS\n"))

    return ALL_DIRECTIONS
}
