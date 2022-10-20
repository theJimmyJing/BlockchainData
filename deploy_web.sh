#!/bin/bash
set -x
ps -A | grep fcc | awk '{print $1}' | xargs kill -9 $1
# git pull
go build
nohup ./fcc >fcc.log 2>&1 &
