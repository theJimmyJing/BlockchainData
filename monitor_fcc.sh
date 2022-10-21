#!/bin/bash

#! /bin/bash

while true 
do
	monitor=`ps -ef | grep fcc | grep -v grep| grep -v monitor | wc -l ` 
	if [ $monitor -eq 0 ] 
	then
		echo "fcc program is not running, restart"
		./deploy_web.sh &
	else
		echo "fcc program is running"
	fi
	sleep 5
done
