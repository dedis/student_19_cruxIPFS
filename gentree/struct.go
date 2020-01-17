package gentree

import (
	"go.dedis.ch/onet/v3"
	"go.dedis.ch/onet/v3/network"
)

// GraphTree Represents the actual graph that will be linked to the Binary
// Tree of the Protocol
type GraphTree struct {
	Tree        *onet.Tree
	ListOfNodes []*onet.TreeNode
	Parents     map[*onet.TreeNode][]*onet.TreeNode
	Radius      float64
}

// TreeConverter is a structure for converting between a recursive tree (graph)
// and a binary tree.
type TreeConverter struct {
	BinaryTree    *onet.Tree
	RecursiveTree *onet.Tree
}

// LocalityNode represents a locality preserving node.
type LocalityNode struct {
	Name string
	IP   map[string]bool
	/*
		X              float64
		Y              float64
	*/
	Level          int
	ADist          []float64 // ADist[l] - minimum distance to level l
	PDist          []string  // pDist[l] - the node at level l whose distance from the crt Node isADist[l]
	Cluster        map[string]bool
	Bunch          map[string]bool
	OptimalCluster map[string]bool
	OptimalBunch   map[string]bool
	Rings          []string
	NrOwnRings     int
	ServerIdentity *network.ServerIdentity
	/*
		AvailablePortsStart int
		AvailablePortsEnd   int
		NextPort            int
		NextPortMtx         sync.Mutex
	*/
}

// LocalityNodes is a list of LocalityNode
type LocalityNodes struct {
	All                   []*LocalityNode
	ServerIdentityToName  map[network.ServerIdentityID]string
	ClusterBunchDistances map[*LocalityNode]map[*LocalityNode]float64
	Links                 map[*LocalityNode]map[*LocalityNode]map[*LocalityNode]bool
}
