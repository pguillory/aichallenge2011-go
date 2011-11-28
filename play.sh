#!/usr/bin/env sh
set -e
turns=$1
shift

cd `dirname $0`
gomake

cd ../tools
# --player_seed 42 
# --map_file maps/test.map \
nice -n 15 ./playgame.py --end_wait=0.25 --verbose --log_dir game_logs --turns $turns --turntime 25000 \
 --fill \
 --strict \
 --serial \
 --map_file maps/symmetric_random_walk/random_walk_01.map \
 "$@" \
 "../go/GoBot" \
 "../c/CBot" \
 "python sample_bots/python/LeftyBot.py" \


 # "python sample_bots/python/HunterBot.py" \
 # "python sample_bots/python/LeftyBot.py" \
 # "python sample_bots/python/GreedyBot.py"
 # --map_file maps/symmetric_random_walk/random_walk_01.map \
 # --map_file maps/maze/maze_5.map \
