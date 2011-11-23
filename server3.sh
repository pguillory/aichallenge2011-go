#!/bin/bash
cd `dirname $0`
bot=./GoBot
while [ -f $bot ]; do
    nice -16 python ../tcpclient.py tcpants.com 2081 "$bot.sh" pguillory Qjg7w96y6F5TSnrm 1
    echo "Game over" | growlnotify
    sleep 1
done
