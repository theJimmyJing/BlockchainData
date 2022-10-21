#!/bin/bash

#! /bin/bash

while true 
do
	monitor=`ps -ef | grep fcc | grep -v grep| grep -v monitor | wc -l ` 
    echo $monitor
	if [ $monitor -eq 0 ] 
	then
		echo "Manipulator program is not running, restart Manipulator"
		./deploy_web.sh &
	else
		echo "Manipulator program is running"
	fi
	sleep 5
done
