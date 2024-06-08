#!/bin/sh
echo "beacon node with MAXPeers=${MAXPEERS} allpeer=${ALLPEERS} and EXECUTE=$EXECUTE, p2pkey=${P2PKEY}"

if [ "$BOOT_DELAY" != "" ] &&  [ "$BOOT_DELAY" != "0" ];then
	        sleep $BOOT_DELAY
fi

sleep 5 && /usr/bin/beacon-chain \
	--datadir=/root/beacondata \
	--min-sync-peers=0 \
	--genesis-state=/root/config/genesis.ssz \
	--bootstrap-node "${BOOTNODE}"\
	--interop-eth1data-votes \
	--chain-config-file=/root/config/config.yml \
	--p2p-max-peers=${MAXPEERS} \
	--p2p-priv-hex=${P2PKEY} \
	--contract-deployment-block=0 \
	--chain-id=${CHAIN_ID:-32382} \
	--rpc-host=0.0.0.0 \
	--grpc-gateway-host=0.0.0.0 \
	--execution-endpoint=http://${EXECUTE}:8551 \
	--accept-terms-of-use \
	--p2p-static-id=true \
	--jwt-secret=/root/config/jwtsecret \
	--suggested-fee-recipient=0x123463a4b065722e99115d6c222f267d9cabb524 \
	--minimum-peers-per-subnet=0 \
	--enable-debug-rpc-endpoints \
	${ALLPEERS} >> /root/beacondata/d.log 2>&1

