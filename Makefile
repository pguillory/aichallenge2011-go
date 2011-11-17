include $(GOROOT)/src/Make.inc

TARG=MyBot
GOFILES=\
    globals.go\
	Direction.go\
	Point.go\
	Player.go\
	Square.go\
	Terrain.go\
	Mystery.go\
	Scent.go\
	Army.go\
	MoveSet.go\
	OrderedMoveList.go\
	Command.go\
	Bot.go\
	MyBot.go\
	Timer.go\
	Log.go\
	toString.go\

include $(GOROOT)/src/Make.cmd
