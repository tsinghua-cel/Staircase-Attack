#! /bin/bash

./prysmctl \
   testnet \
   generate-genesis \
   --fork=capella \
   --num-validators=1500 \
   --genesis-time-delay=15 \
   --output-ssz=genesis.ssz \
   --chain-config-file=config.yml \
   --geth-genesis-json-in=genesis.json \
   --geth-genesis-json-out=genesis.json
 
./geth --datadir=gethdata init genesis.json
./geth --datadir=gethdata account import secret.json <<EOM
123
123
EOM

./geth \
   --http \
   --http.api=eth,net,web3 \
   --http.addr=0.0.0.0 \
   --http.corsdomain=* \
   --ws \
   --ws.api=eth,net,web3 \
   --ws.addr=0.0.0.0 \
   --ws.origins=* \
   --authrpc.vhosts=* \
   --authrpc.addr=0.0.0.0 \
   --authrpc.jwtsecret=jwtsecret \
   --datadir=gethdata \
   --allow-insecure-unlock \
   --unlock=0x123463a4b065722e99115d6c222f267d9cabb524 \
   --password=password.txt \
   --nodiscover \
   --syncmode=full
