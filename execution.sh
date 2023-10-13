#! /bin/bash

rm -rf gethdata/
rm -rf genesis.ssz

rm -rf beacon1/beacondata/
rm -rf beacon1/validatordata/
rm -rf beacon1/balance.txt
rm -rf beacon1/balance.json
rm -rf beacon1/p2p.txt

rm -rf beacon2/beacondata/
rm -rf beacon2/validatordata/
rm -rf beacon2/genesistime.txt
rm -rf beacon2/uatt.json
rm -rf beacon2/duties.json
rm -rf beacon2/att.json
rm -rf beacon2/proposer_slot.txt

./geth --datadir=gethdata init genesis.json
echo -e "\n\n" | ./geth --datadir=gethdata account import secret.json

./prysmctl testnet generate-genesis --num-validators=128 --output-ssz=genesis.ssz --chain-config-file=config.yml

./geth --http --http.api "eth,engine" --datadir=gethdata --allow-insecure-unlock --unlock="0x123463a4b065722e99115d6c222f267d9cabb524" --password=password.txt  --nodiscover console --syncmode=full --mine --miner.etherbase=0x123463a4b065722e99115d6c222f267d9cabb524