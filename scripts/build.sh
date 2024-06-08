#!/bin/bash
export  DOCKER_BUILDKIT=1
docker build -t geth:v1.13-base -f dockerfile/geth.Dockerfile .
docker build -t attacker:latest -f dockerfile/attacker.Dockerfile .
docker build -t beacon:base -f dockerfile/beacon.Dockerfile .
docker build -t beacon:modified -f dockerfile/beacon.modify.Dockerfile .
docker build -t validator:base -f dockerfile/validator.Dockerfile .
docker build -t validator:modified -f dockerfile/validator.modify.Dockerfile .
