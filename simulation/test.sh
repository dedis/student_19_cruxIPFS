#!/bin/sh

go build
./simulation -platform deterlab -mport 10008 ipfs.toml
