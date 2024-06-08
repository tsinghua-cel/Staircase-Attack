#!/bin/bash

echo "build ethtools image at $PWD"
docker build -t ethnettools -f dockerfile/ethtools.Dockerfile .
docker run -it --rm -v "${PWD}/ethnet/config:/root/config" --name generate --entrypoint /usr/bin/prysmctl ethnettools \
	testnet \
	generate-genesis \
	--fork=capella \
	--num-validators=1000 \
	--genesis-time-delay=15 \
	--output-ssz=/root/config/genesis.ssz \
	--chain-config-file=/root/config/config.yml \
	--geth-genesis-json-in=/root/config/genesis.json \
	--geth-genesis-json-out=/root/config/genesis.json

