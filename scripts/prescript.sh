#!/bin/sh

# kill running ipfs & ipfs-cluster-service instances
nohup killall -9 ipfs >/dev/null 2>&1 &
nohup killall -9 ipfs-cluster-se >/dev/null 2>&1 &

# install ipfs if not already installed
if ! [ -x "$(command -v ipfs)" ]; then
  echo 'Installing ipfs'
  sudo cp ipfs /usr/local/bin
fi

# install ipfs-cluster-service if not already installed
if ! [ -x "$(command -v ipfs-cluster-service)" ]; then
  echo 'Installing ipfs-cluster-service'
  sudo cp ipfs-cluster-service /usr/local/bin
fi

# install ipfs-cluster-ctl if not already installed
#if ! [ -x "$(command -v ipfs-cluster-ctl)" ]; then
#  echo 'Installing ipfs-cluster-ctl'
#  sudo cp ipfs-cluster-ctl /usr/local/bin
#fi
