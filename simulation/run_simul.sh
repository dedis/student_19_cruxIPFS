#!/bin/bash

go build

N=${1-20}
printf 'Simulation = "IPFS"\nServers = '$N'\nBf = '$(($N-1))'\nRounds = 1\nSuite = "Ed25519"\nPrescript = "clean.sh"\n\nDepth\n1' > ipfs.toml
./simulation ipfs.toml
