include $(GOROOT)/src/Make.inc

TARG=MyBot
GOFILES=\
    globals.go\
	Direction.go\
	Point.go\
	Player.go\
	Square.go\
	Map.go\
	Scent.go\
	Moves.go\
	Bot.go\
	MyBot.go\

include $(GOROOT)/src/Make.cmd
