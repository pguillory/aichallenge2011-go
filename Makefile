include $(GOROOT)/src/Make.inc

TARG=MyBot
GOFILES=\
    globals.go\
	Direction.go\
	Point.go\
	Player.go\
	Square.go\
	Map.go\
	Mystery.go\
	Scent.go\
	Moves.go\
	Bot.go\
	MyBot.go\
	Timer.go\
	Log.go\

include $(GOROOT)/src/Make.cmd
