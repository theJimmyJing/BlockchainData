#!/bin/bash

#! /bin/bash

# 使用crontab定时执行
# crontab -e 然后输入
# * * * * *  /..路径/Tokenlist/monitor_fcc.sh 2&>1 >> /路径../check-alive.log  #每分钟检查一次
# */5 * * * *  /..路径/Tokenlist/monitor_fcc.sh 2&>1 >> /路径../check-alive.log  #每5分钟检查一次

# while true 
# do
	monitor=`ps -ef | grep fcc | grep -v grep| grep -v monitor | wc -l ` 
	if [ $monitor -eq 0 ] 
	then
		echo "fcc program is not running, restart"
		./deploy_web.sh &
	# else
	# 	echo "fcc program is running"
	fi
	# sleep 10
# done
