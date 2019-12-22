#!/bin/sh

rm output.txt

go build
#./simulation -platform deterlab -mport 10008 ipfs.toml > output.txt
./simulation ipfs.toml > output.txt
