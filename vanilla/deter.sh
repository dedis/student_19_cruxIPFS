#!/bin/sh

rm output.txt

go build
./vanilla -platform deterlab -mport 10008 ipfs.toml > output.txt
#./vanilla ipfs.toml > output.txt
