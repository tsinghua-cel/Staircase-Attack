#! /bin/bash

rm -rf gethdata/
rm -rf beacon1/beacondata/
rm -rf beacon1/validatordata/
rm -rf beacon2/beacondata/
rm -rf beacon2/validatordata2/

./geth --datadir=gethdata init genesis.json
./geth --datadir=gethdata account import secret.json

./prysmctl testnet generate-genesis --num-validators=12 --output-ssz=genesis.ssz --chain-config-file=config.yml

./geth --http --http.api "eth,engine" --datadir=gethdata --allow-insecure-unlock --unlock="0x123463a4b065722e99115d6c222f267d9cabb524" --password="" --nodiscover console --syncmode=full --mine --miner.etherbase=0x123463a4b065722e99115d6c222f267d9cabb524
