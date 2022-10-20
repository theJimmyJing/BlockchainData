#!/bin/bash
set -x

cd /root/ || exit
# git reset --hard
rm -rf /root/FreeChat-WEB/
git clone https://github.com/FreeChatDevelopment/FreeChat-WEB.git

cd /root/FreeChat-WEB/
npm install
npm run build

rm -rf /usr/share/nginx/fcc/*
cp -rf build/* /usr/share/nginx/fcc/
nginx -t
nginx -s reload
