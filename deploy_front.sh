#!/bin/bash
set -x

cd /root/FreeChat-WEB || exit
# git reset --hard
git pull
npm install
npm run build
rm -rf /usr/share/nginx/fcc/*
cp -rf build/* /usr/share/nginx/fcc/
nginx -t
nginx -s reload
