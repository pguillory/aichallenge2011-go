#!/bin/bash
set -e
cd `dirname $0`
gomake
bot=./GoBot
while [ -f "$bot" ]; do
    nice -16 python ../tcpclient.py tcpants.com 2081 "$bot" pguillory Qjg7w96y6F5TSnrm 1
    echo "Game over" | growlnotify
    sleep 1
done
