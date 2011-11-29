#!/bin/bash
set -e
cd `dirname $0`
gomake
bot=./GoBot
while [ -f "$bot" ]; do
    nice -16 python ../tcpclient.py bhickey.net  2081 "$bot" pguillory_jr 5SPT7pcMrjfDBqJT 1
    echo "Game over" | growlnotify
    sleep 1
done
