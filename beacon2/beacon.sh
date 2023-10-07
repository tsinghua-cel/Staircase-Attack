script_dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
file_path="$script_dir/../beacon1/p2p.txt"
file_content=$(cat "$file_path")

export PEER=$file_content
echo "$PEER"
./beacon-chain   --datadir=beacondata   --min-sync-peers=1   --genesis-state=../genesis.ssz   --bootstrap-node=   --chain-config-file=../config.yml   --config-file=../config.yml   --chain-id=32382   --execution-endpoint=http://localhost:8551   --accept-terms-of-use   --rpc-port=4001   --p2p-tcp-port=13001   --p2p-udp-port=12001   --grpc-gateway-port=3501   --monitoring-port=8001   --jwt-secret=../gethdata/geth/jwtsecret   --peer=$PEER
