#!/bin/bash
set -e
cd `dirname $0`
gomake
bot=./GoBot
while [ -f "$bot" ]; do
    nice -16 python ../tcpclient.py ants.fluxid.pl 2081 "$bot" pguillory 8xR8UQBYWjb75Ut9 1
    echo "Game over" | growlnotify
    sleep 5
done
