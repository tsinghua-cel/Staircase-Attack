#!/bin/bash
export  DOCKER_BUILDKIT=1
BUILDFLAG=""
if [ "$REBUILD" = "1" ];then
	BUILDFLAG="$BUILDFLAG --no-chace"
fi
docker build $BUILDFLAG -t geth:v1.13-base -f dockerfile/geth.Dockerfile .
docker build $BUILDFLAG -t attacker:latest -f dockerfile/attacker.Dockerfile .
docker build $BUILDFLAG -t beacon:base -f dockerfile/beacon.Dockerfile .
docker build $BUILDFLAG -t beacon:modified -f dockerfile/beacon.modify.Dockerfile .
docker build $BUILDFLAG -t validator:base -f dockerfile/validator.Dockerfile .
docker build $BUILDFLAG -t validator:modified -f dockerfile/validator.modify.Dockerfile .
