package main

import "path/filepath"

const (
	installFolder = "install"
)

var ipfsLocation string
var ipfsClusterLocation string

func init() {
	ipfsLocation = filepath.Join(installFolder, "ipfs")
	ipfsClusterLocation = filepath.Join(installFolder, "ipfs-cluster-service")
}
