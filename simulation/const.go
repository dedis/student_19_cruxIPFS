package main

import (
	"path/filepath"

	cruxIPFS "github.com/dedis/student_19_cruxIPFS"
)

const (
	rootFolder      = ".."
	installFolder   = "install"
	ipfsFile        = "ipfs"
	ipfsClusterFile = "ipfs-cluster-service"
	ipfsCtlFile     = "ipfs-cluster-ctl"
	prescriptFile   = "prescript.sh"
	nodesFile       = "nodes.txt"
	genFolder       = "gen"
	detailsFile     = "details.txt"
	gendetailsFile  = "gendetails.txt"
)

var dataLocation string
var scriptLocation string
var installLocation string
var ipfsLocation string
var ipfsClusterLocation string
var ipfsCtlLocation string
var prescriptLocation string
var nodesLocation string
var gendetailsLocation string
var simdetailsLocation string

var mode = "raft"
var computePings = true
var cruxified = true
var remote = true
var nOps = 100

func init() {
	dataLocation = filepath.Join(rootFolder, cruxIPFS.DataFolder)
	scriptLocation = filepath.Join(rootFolder, cruxIPFS.ScriptsFolder)
	installLocation = filepath.Join(dataLocation, installFolder)
	ipfsLocation = filepath.Join(installLocation, ipfsFile)
	ipfsClusterLocation = filepath.Join(installLocation, ipfsClusterFile)
	ipfsCtlLocation = filepath.Join(installLocation, ipfsCtlFile)
	prescriptLocation = filepath.Join(scriptLocation, prescriptFile)
	nodesLocation = filepath.Join(dataLocation, nodesFile)
	gendetailsLocation = filepath.Join(dataLocation, genFolder, gendetailsFile)
	simdetailsLocation = filepath.Join(dataLocation, detailsFile)
}
