#!/bin/bash
curdir=${PWD}
./bin/bootnode -discv5-port 39000 -metrics-port 5100 -private 37871f094f89fc95599e73d72d3008c828bb76d0864759011fa14dbaf6cf3c06 -fork-version 20000089  -debug true > boot.log 2>&1 &

#docker run -it -d --restart=always --name bootnode -p 39000:39000 -p 5000:5000 --entrypoint /usr/bin/bootnode ethnet/tools -discv5-port 39000 -private 37871f094f89fc95599e73d72d3008c828bb76d0864759011fa14dbaf6cf3c06 -fork-version 20000089 -external-ip "172.17.0.1" -debug true
