#! /bin/bash

ps -aux|grep geth| grep -v grep| awk '{print $2}'| xargs kill -9
ps -aux|grep beacon| grep -v grep| awk '{print $2}'| xargs kill -9
ps -aux|grep validator| grep -v grep| awk '{print $2}'| xargs kill -9

rm -rf gethdata/
rm -rf genesis.ssz

rm -rf beacon1/beacondata/
rm -rf beacon1/validatordata/
rm -rf balance.txt
rm -rf p2p.txt

rm -rf beacon2/beacondata/
rm -rf beacon2/validatordata/
rm -rf genesistime.txt
rm -rf uatt.json
rm -rf duties.json
rm -rf att.json
rm -rf proposer_slot.txt

