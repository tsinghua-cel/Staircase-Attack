#! /bin/bash

rm -rf gethdata/
rm -rf genesis.ssz

rm -rf beacon1/beacondata/
rm -rf beacon1/validatordata/
rm -rf balance.txt
rm -rf balance.json
rm -rf p2p.txt

rm -rf beacon2/beacondata/
rm -rf beacon2/validatordata/
rm -rf genesistime.txt
rm -rf uatt.json
rm -rf duties.json
rm -rf att.json
rm -rf proposer_slot.txt

rm -rf out

echo "===========Start execution layer================"
./geth --datadir=gethdata init genesis.json
echo -e "\n\n" | ./geth --datadir=gethdata account import secret.json
./prysmctl testnet generate-genesis --num-validators=1000 --output-ssz=genesis.ssz --chain-config-file=config.yml
nohup ./geth --http --http.api "eth,engine" --datadir=gethdata --allow-insecure-unlock --unlock="0x123463a4b065722e99115d6c222f267d9cabb524" --password=password.txt  --nodiscover --syncmode=full --mine --miner.etherbase=0x123463a4b065722e99115d6c222f267d9cabb524 > out/execution.log 2>&1 &

sleep 3

echo "===========Start honest validators=============="
nohup ./beacon1/beacon-chain --datadir=beacon1/beacondata --min-sync-peers=0 --genesis-state=genesis.ssz --bootstrap-node= --chain-config-file=config.yml --config-file=config.yml --chain-id=32382 --execution-endpoint=http://localhost:8551 -accept-terms-of-use --jwt-secret=gethdata/geth/jwtsecret > out/beacon1.log 2>&1 &
nohup ./beacon1/validator --datadir=beacon1/validatordata --accept-terms-of-use --interop-num-validators=666 --interop-start-index=0 --force-clear-db --chain-config-file=config.yml --config-file=config.yml > out/validator1.log 2>&1 &
nohup python3 ./beacon1/balance_monitor.py > out/python.log 2>&1 &


sleep 3

echo "===========Start Byzantine validators==========="
script_dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
file_path="$script_dir/p2p.txt"
file_content=$(cat "$file_path")

export PEER=$file_content
nohup ./beacon2/beacon-chain   --datadir=beacon2/beacondata   --min-sync-peers=1   --genesis-state=genesis.ssz   --bootstrap-node=   --chain-config-file=config.yml   --config-file=config.yml   --chain-id=32382   --execution-endpoint=http://localhost:8551   --accept-terms-of-use   --rpc-port=4001   --p2p-tcp-port=13001   --p2p-udp-port=12001   --grpc-gateway-port=3501   --monitoring-port=8001   --jwt-secret=gethdata/geth/jwtsecret   --peer=$PEER > out/beacon2.log 2>&1 &
nohup ./beacon2/validator   --datadir=beacon2/validatordata   --accept-terms-of-use   --interop-start-index=666   --interop-num-validators=334   --force-clear-db   --beacon-rpc-provider=127.0.0.1:4001   --grpc-gateway-port=3501   --chain-config-file=config.yml   --config-file=config.yml > out/validator2.log 2>&1 &
nohup python3 ./beacon2/role_monitor.py > out/python.log 2>&1 &

rm -rf out