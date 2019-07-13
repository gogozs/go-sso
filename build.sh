#!/usr/bin/env bash

# git
echo "更新代码..."
git pull origin master

# docker
echo "docker构建..."
docker-compose -f docker-compose.prd.yml up -d
docker-compose -f docker-compose.prd.yml build
