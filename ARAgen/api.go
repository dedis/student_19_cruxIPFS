package ARAgen

import (
	"github.com/dedis/student_19_cruxIPFS/gentree"
	"go.dedis.ch/onet/v3"
	"go.dedis.ch/onet/v3/log"
	"go.dedis.ch/onet/v3/network"
)

const RND_NODES = false
const NR_LEVELS = 3
const OPTIMIZED = false
const OPTTYPE = 1
const MIN_BUNCH_SIZE = 12

// Runs the compact routing algorithm to generate ARAs.
// It returns "AuxNodes", a LocalityNodes structure with all compact routing relevant fields filled in (bunch, cluster etc), as well
// as all fields that "Nodes" has; the caller can use it to replace the initial "Nodes" structure
// "dist2" the compact distance matrix relevant for the node "rootNodeName",
// "ARATreeStruct" a map with a graphTree slice for every node in "Nodes"; each slice contains a graphTree for each ARA
// created by that node,
// and "ARAOnetTrees", a map with a onet tree slice for every node in "Nodes"; each slice contains an onet tree for
// each ARA created by that node, that the caller can use to run protocols on.

func GenARAs(Nodes gentree.LocalityNodes, rootNodeName string, PingDistances map[string]map[string]float64,
	NrLevels int) (AuxNodes gentree.LocalityNodes, dist2 map[*gentree.LocalityNode]map[*gentree.LocalityNode]float64,
	ARATreeStruct map[string][]GraphTree, ARAOnetTrees map[string][]*onet.Tree) {

	AuxNodes.All = make([]*gentree.LocalityNode, len(Nodes.All))
	AuxNodes.ServerIdentityToName = make(map[network.ServerIdentityID]string)
	AuxNodes.ClusterBunchDistances = make(map[*gentree.LocalityNode]map[*gentree.LocalityNode]float64)
	AuxNodes.Links = make(map[*gentree.LocalityNode]map[*gentree.LocalityNode]map[*gentree.LocalityNode]bool)

	ARATreeStruct = make(map[string][]GraphTree)
	ARAOnetTrees = make(map[string][]*onet.Tree)

	for i, n := range Nodes.All {
		AuxNodes.All[i].IP = n.IP
		AuxNodes.All[i].Name = n.Name
		AuxNodes.All[i].AvailablePortsStart = n.AvailablePortsStart
		AuxNodes.All[i].ServerIdentity = n.ServerIdentity
		AuxNodes.All[i].Level = n.Level
		AuxNodes.All[i].X = n.X
		AuxNodes.All[i].Y = n.Y
	}

	for k, v := range Nodes.ServerIdentityToName {
		AuxNodes.ServerIdentityToName[k] = v
	}

	gentree.CreateLocalityGraph(AuxNodes, RND_NODES, RND_NODES, NrLevels, PingDistances)

	// we are rooting trees here
	myname := rootNodeName

	if OPTIMIZED {
		gentree.OptimizeGraph(AuxNodes, myname, MIN_BUNCH_SIZE, OPTTYPE)
	}

	dist2 = gentree.AproximateDistanceOracle(AuxNodes)

	// Generate trees for all nodes
	for _, crtRoot := range AuxNodes.All {
		crtRootName := crtRoot.Name

		tree, NodesList, Parents, TreeRadiuses := gentree.CreateOnetRings(AuxNodes, crtRootName, dist2)

		for i, n := range tree {
			ARATreeStruct[crtRootName] = append(ARATreeStruct[crtRootName], GraphTree{
				n,
				NodesList[i],
				Parents[i],
				TreeRadiuses[i],
			})
		}
	}

	for rootName, graphTrees := range ARATreeStruct {
		for _, n := range graphTrees {

			rosterNames := make([]string, 0)
			for _, si := range n.Tree.Roster.List {
				rosterNames = append(rosterNames, AuxNodes.GetServerIdentityToName(si))
			}

			log.Lvl2("generation node ", rootNodeName, "rootName x", rootName, "creates binary with roster", rosterNames)

			ARAOnetTrees[rootName] = append(ARAOnetTrees[rootName], createBinaryTreeFromGraphTree(n))
		}
	}

	log.Lvl2("done")

	return AuxNodes, dist2, ARATreeStruct, ARAOnetTrees
}

//Computes an onet binary tree based on a GraphTree
func createBinaryTreeFromGraphTree(GraphTree GraphTree) *onet.Tree {

	BinaryTreeRoster := GraphTree.Tree.Roster
	Tree := BinaryTreeRoster.GenerateBinaryTree()

	return Tree
}
