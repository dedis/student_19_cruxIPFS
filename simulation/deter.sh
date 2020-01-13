#!/bin/bash

go build
./simulation -platform deterlab -mport 10008 ipfs.toml > output_c.txt

while ! grep -q "Done" output_c.txt; do
    sleep 15
done

echo "Done"
