#!/bin/bash

go build
./simulation -platform deterlab -mport 10008 ipfs.toml > output_c.txt

echo "Done"
