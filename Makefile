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
	HolyGround.go\
	Situation.go\
	Prediction.go\
	Mystery.go\
	Scent.go\
	Army.go\
	Move.go\
	MoveSet.go\
	Focus.go\
	OrderedMoveList.go\
	Command.go\
	Bot.go\
	MyBot.go\
	Timer.go\
	Log.go\
	toString.go\

include $(GOROOT)/src/Make.cmd
