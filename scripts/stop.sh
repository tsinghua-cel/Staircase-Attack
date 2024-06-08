#!/bin/bash
case=`cat .case`
docker compose -f "$case" down && rm .case
