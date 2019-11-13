#!/usr/bin/env bash

# git
cd $GO_SSO_WORKDIR
echo "更新代码..."
git fetch --all
# 强制远程分支覆盖本地分支
git reset --hard origin/master

# docker
echo "docker构建..."
docker-compose -f docker-compose.prd.yml build
docker-compose -f docker-compose.prd.yml up -d
