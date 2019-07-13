#!/usr/bin/env bash

# git
cd $GO_WEIXIN_WORKDIR
echo "更新代码..."
git pull origin master

# docker
echo "docker构建..."
docker-compose -f docker-compose.prd.yml up -d
docker-compose -f docker-compose.prd.yml build
