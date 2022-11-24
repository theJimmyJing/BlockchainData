#!/bin/bash
#set -x
#ps -A | grep fcc | grep -v monitor | awk '{print $1}' | xargs kill -9 $1
#git pull
#go build
#nohup ./fcc >fcc.log 2>&1 &

git pull && cd ../Private_Config
git pull && cp -rv ./blockchain-config/config.yaml ../BlockchainData/config/ && cd ../BlockchainData
go build && sudo systemctl restart blockchaindata
