#! /bin/bash

./beacon-chain \
   --datadir=beacondata \
   --min-sync-peers=0 \
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
   --jwt-secret=jwtsecret \
   --suggested-fee-recipient=0x123463a4b065722e99115d6c222f267d9cabb524 \
   --minimum-peers-per-subnet=0 \
   --enable-debug-rpc-endpoints