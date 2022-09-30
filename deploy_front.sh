#!/bin/bash
set -x
rm -rf /usr/share/nginx/fcc/*
cd /root/FreeChat-WEB || exit
# git reset --hard
git pull
cp -rf build/* /usr/share/nginx/fcc/*
nginx -t
nginx -s reload
