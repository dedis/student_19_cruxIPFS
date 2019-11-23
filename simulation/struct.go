package main

import (
	"github.com/dedis/student_19_cruxIPFS/gentree"
	"go.dedis.ch/onet/v3"
)

// IPFSSimulation only holds the BFTree simulation
type IPFSSimulation struct {
	onet.SimulationBFTree

	Nodes   gentree.LocalityNodes
	Parents map[string][]string

	NodePathLocal  string
	NodePathRemote string
}
