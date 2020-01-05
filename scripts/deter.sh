#!/bin/sh


cd ../simulation
go build
rm ../data/output.txt > /dev/null 2>&1
./simulation -platform deterlab -mport 10008 ipfs.toml > ../data/output.txt
