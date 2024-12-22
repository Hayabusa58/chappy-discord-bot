#!/bin/bash

# Docker コンテナをビルドし起動するスクリプト

# コンテナが稼働していれば止める
cid=`docker ps | grep chappy-discord-bot | awk '{ print $1}'`

if [ -n "$cid" ]; then
  echo "stopping contaiter: $cid"
  docker stop $cid
fi

docker build -t chappy-discord-bot .
docker run -d --rm --name chappy-discord-bot chappy-discord-bot

cid=`docker ps | grep chappy-discord-bot | awk '{ print $1}'`
echo "started contaier: $cid"



