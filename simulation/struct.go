package main

import (
	"github.com/dedis/student_19_cruxIPFS/gentree"
	"go.dedis.ch/onet/v3"
)

// NODEPATHREMOTE path to remote nodes file
var NODEPATHREMOTE = ""

// NODEPATHLOCAL path to local nodes file
var NODEPATHLOCAL = ""

// IPFSSimulation only holds the BFTree simulation
type IPFSSimulation struct {
	onet.SimulationBFTree

	Nodes   gentree.LocalityNodes
	Parents map[string][]string
}
