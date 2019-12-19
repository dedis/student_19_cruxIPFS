#!/bin/sh

go build
./vanilla -platform deterlab -mport 10008 ipfs.toml | grep min > results.txt

cat results.txt | cut -d ' ' -f11 > min.txt
cat results.txt | cut -d ' ' -f13 > max.txt
