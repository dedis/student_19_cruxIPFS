package service

import (
	"sync"

	"github.com/dedis/paper_crux/dsn_exp/gentree"
	"go.dedis.ch/onet"
	"go.dedis.ch/onet/network"
)

// Service is our template-service
type Service struct {
	// We need to embed the ServiceProcessor, so that incoming messages
	// are correctly handled.
	*onet.ServiceProcessor

	Nodes        gentree.LocalityNodes
	LocalityTree *onet.Tree
	Parents      []*onet.TreeNode
	GraphTree    map[string][]GraphTree
	BinaryTree   map[string][]*onet.Tree
	alive        bool
	Distances    map[*gentree.LocalityNode]map[*gentree.LocalityNode]float64

	OwnPings      map[string]float64
	DonePing      bool
	PingDistances map[string]map[string]float64
	NrPingAnswers int
	PingAnswerMtx sync.Mutex
	PingMapMtx    sync.Mutex

	storage *storage
}

// GraphTree structure of the tree made of nodes
type GraphTree struct {
	Tree        *onet.Tree
	ListOfNodes []*onet.TreeNode
	Parents     map[*onet.TreeNode][]*onet.TreeNode
	Radius      float64
}

// InitRequest to generate trees
type InitRequest struct {
	Nodes                []*gentree.LocalityNode
	ServerIdentityToName map[*network.ServerIdentity]string
	NrOps                int
	OpIdxStart           int
	Roster               *onet.Roster
}

// InitResponse to init request
type InitResponse struct {
}

// storage is used to save our data.
type storage struct {
	Count int
	sync.Mutex
}
