#!/bin/bash
set -x
rm -rf /usr/share/nginx/fcc/*
cd /root/FreeChat-WEB || exit
# git reset --hard
git pull
npm install
npm run build
cp -rf build/* /usr/share/nginx/fcc/
nginx -t
nginx -s reload
