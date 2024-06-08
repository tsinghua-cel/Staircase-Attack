#!/bin/bash

CASEFILE=${1:-""}

if [ "$CASEFILE" ==  "" ];then
  echo "invalid param, case file is empty"
  exit 1
fi

# check $CASEFILE is exist in testcase dir
if [ ! -f "$CASEFILE" ];then
  echo "invalid param, case file is not exist"
  exit 1
fi

# build docker images
./scripts/build.sh || exit 1

# clear generated dir
./scripts/clear.sh || exit 1

# generate genesis
./scripts/genesis.sh || exit 1

# start docker compose
docker compose -f $CASEFILE up -d && echo "$CASEFILE" > .case

echo "start testcase success"
