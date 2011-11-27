include $(GOROOT)/src/Make.inc

TARG=GoBot
GOFILES=\
    globals.go\
	Direction.go\
	Point.go\
	PointSet.go\
	PointQueue.go\
	Player.go\
	Square.go\
	Terrain.go\
	Mystery.go\
	PotentialEnemy.go\
	Scrum.go\
	Distance.go\
	Situation.go\
	Prediction.go\
	AdjacentWater.go\
	Army.go\
	Move.go\
	MoveSet.go\
	Repulsion.go\
	Focus.go\
	OrderedMoveList.go\
	RageVirus.go\
	Command.go\
	Bot.go\
	MyBot.go\
	Timer.go\
	Log.go\
	toString.go\

include $(GOROOT)/src/Make.cmd
