#!/bin/bash

folder='K3N20D150remoteO2000crdt/'
output_v=$folder'output_v.txt'
output_c=$folder'output_c.txt'
data=$folder'data/'

mkdir $data
mkdir $folder'graphs'

cp '../data/ipfs.toml' $folder
cp '../data/nodes.txt' $folder
cp '../data/details.txt' $folder

cat $output_v | grep "ping node" > $data'pings.txt'
cat $output_v | grep minoptime > $data'vanilla.txt'
cat $output_c | grep minoptime > $data'min.txt'
cat $output_c | grep maxoptime > $data'max.txt'
