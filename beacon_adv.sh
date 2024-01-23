#! /bin/bash

script_dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
file_path="$script_dir/p2p.txt"
file_content=$(cat "$file_path")

export PEER=$file_content

./beacon-chain_adv \
   --datadir=beacondata2 \
   --min-sync-peers=1 \
   --genesis-state=genesis.ssz \
   --bootstrap-node= \
   --interop-eth1data-votes \
   --chain-config-file=config.yml \
   --contract-deployment-block=0 \
   --chain-id=${CHAIN_ID:-32382} \
   --rpc-host=0.0.0.0 \
   --grpc-gateway-host=0.0.0.0 \
   --execution-endpoint=http://localhost:8551 \
   --accept-terms-of-use \
   --rpc-port=4001 \
   --p2p-tcp-port=13001 \
   --p2p-udp-port=12001 \
   --grpc-gateway-port=3501 \
   --monitoring-port=8001 \
   --jwt-secret=jwtsecret \
   --suggested-fee-recipient=0x123463a4b065722e99115d6c222f267d9cabb524 \
   --minimum-peers-per-subnet=0 \
   --enable-debug-rpc-endpoints \
   --peer=$PEER