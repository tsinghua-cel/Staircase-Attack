#! /bin/bash

rm -rf gethdata/
rm -rf genesis.ssz

rm -rf beacondata/
rm -rf validatordata/
rm -rf balance.txt
rm -rf balance.json
rm -rf p2p.txt

rm -rf beacondata2/
rm -rf validatordata2/
rm -rf genesistime.txt
rm -rf uatt.json
rm -rf duties.json
rm -rf att.json
rm -rf proposer_slot.txt


rm -rf out
mkdir out

sleep 3

echo "===========Start execution layer================"
./geth --datadir=gethdata init genesis.json
./geth --datadir=gethdata account import secret.json <<EOM
123
123
EOM
./prysmctl \
   testnet \
   generate-genesis \
   --fork=capella \
   --num-validators=16 \
   --genesis-time-delay=15 \
   --output-ssz=genesis.ssz \
   --chain-config-file=config.yml \
   --geth-genesis-json-in=genesis.json \
   --geth-genesis-json-out=genesis.json
nohup ./geth \
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
   --syncmode=full > out/execution.log 2>&1 &

sleep 10

echo "===========Start honest validators=============="
nohup ./beacon-chain \
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
   --enable-debug-rpc-endpoints > out/beacon_hon.log 2>&1 &
nohup ./validator   --datadir=validatordata   --accept-terms-of-use   --interop-num-validators=8   --interop-start-index=0   --force-clear-db   --chain-config-file=config.yml   --config-file=config.yml > out/validator_hon.log 2>&1 &


sleep 10

echo "===========Start Byzantine validators==========="
script_dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
file_path="$script_dir/p2p.txt"
file_content=$(cat "$file_path")
export PEER=$file_content

nohup ./beacon-chain_adv \
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
   --peer=$PEER > out/beacon_adv.log 2>&1 &
nohup ./validator_adv   --datadir=validatordata2   --accept-terms-of-use   --interop-start-index=8   --interop-num-validators=8   --force-clear-db   --beacon-rpc-provider=127.0.0.1:4001   --grpc-gateway-port=3501   --chain-config-file=config.yml   --config-file=config.yml > out/validator_adv.log 2>&1 &

sleep 10

nohup python3 execution.py > out/python.log 2>&1 &