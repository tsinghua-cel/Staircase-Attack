#!/bin/bash
curdir=${PWD}
echo "build ethtools image"
docker build -t ethnettools -f dockerfile/ethtools.Dockerfile .
docker run -it --rm -v "${curdir}/config:/root/config" --name generate --entrypoint /usr/bin/prysmctl ethnettools \
	testnet \
	generate-genesis \
	--fork=capella \
	--num-validators=1000 \
	--genesis-time-delay=15 \
	--output-ssz=/root/config/genesis.ssz \
	--chain-config-file=/root/config/config.yml \
	--geth-genesis-json-in=/root/config/genesis.json \
	--geth-genesis-json-out=/root/config/genesis.json

