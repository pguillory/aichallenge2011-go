#!/usr/bin/env sh
set -e
#turns=$1
#shift

cd `dirname $0`
gomake

cd ../tools
# --player_seed 42 
# --map_file maps/test.map \
./playgame.py --verbose --log_dir game_logs --turntime 1500 \
 --fill \
 --serial \
 --map_file maps/random_walk/random_walk_03p_01.map \
 "$@" \
 "../go/GoBot -debug" \
 "../c/CBot" \
 "python sample_bots/python/LeftyBot.py" \


 # --strict \


 # --map_file maps/cell_maze/cell_maze_p03_01.map \
 # --map_file maps/random_walk/random_walk_03p_03.map \
 # "python sample_bots/python/HunterBot.py" \
 # "python sample_bots/python/LeftyBot.py" \
 # "python sample_bots/python/GreedyBot.py"
 # --map_file maps/symmetric_random_walk/random_walk_01.map \
 # --map_file maps/maze/maze_5.map \
