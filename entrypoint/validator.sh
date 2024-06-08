#!/bin/sh

echo "bootstart validator with validators_num=${VALIDATORS_NUM} and index=${VALIDATORS_INDEX} and beacon $BEACONRPC"

if [ "$BOOT_DELAY" != "" ] &&  [ "$BOOT_DELAY" != "0" ];then
	        sleep $BOOT_DELAY
fi

/usr/bin/validator   --datadir=/root/validatordata   \
	--accept-terms-of-use   \
	--interop-num-validators=${VALIDATORS_NUM}   \
	--interop-start-index=${VALIDATORS_INDEX}   \
	--force-clear-db   \
	--chain-config-file=/root/config/config.yml   \
	--config-file=/root/config/config.yml \
	--beacon-rpc-provider "$BEACONRPC" >> /root/validatordata/d.log 2>&1 
