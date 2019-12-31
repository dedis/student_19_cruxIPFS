package gentree

import (
	"math"

	"go.dedis.ch/onet/v3"
	"go.dedis.ch/onet/v3/log"
	"go.dedis.ch/onet/v3/network"
)

const RND_NODES = false
const RND_LVLS = RND_NODES
const NR_LEVELS = 3
const OPTIMIZED = false
const OPTTYPE = 1
const MIN_BUNCH_SIZE = 12

// GenARAs :
// Runs the compact routing algorithm to generate ARAs.
// It returns "AuxNodes", a LocalityNodes structure with all compact routing
// relevant fields filled in (bunch, cluster etc), as well as all fields that
// "Nodes" has; the caller can use it to replace the initial "Nodes" structure
// "dist2" the compact distance matrix relevant for the node "rootNodeName",
// "ARATreeStruct" a map with a graphTree slice for every node in "Nodes"; each
// slice contains a graphTree for each ARA created by that node,
// and "ARAOnetTrees", a map with a onet tree slice for every node in "Nodes";
// each slice contains an onet tree for each ARA created by that node, that the
// caller can use to run protocols on.
func GenARAs(Nodes LocalityNodes, rootNodeName string,
	PingDistances map[string]map[string]float64, NrLevels int) (
	AuxNodes LocalityNodes, dist2 map[*LocalityNode]map[*LocalityNode]float64,
	ARATreeStruct map[string][]GraphTree,
	ARAOnetTrees map[string][]*onet.Tree) {

	AuxNodes.All = make([]*LocalityNode, len(Nodes.All))
	AuxNodes.ServerIdentityToName = make(map[network.ServerIdentityID]string)
	AuxNodes.ClusterBunchDistances =
		make(map[*LocalityNode]map[*LocalityNode]float64)
	AuxNodes.Links =
		make(map[*LocalityNode]map[*LocalityNode]map[*LocalityNode]bool)

	ARATreeStruct = make(map[string][]GraphTree)
	ARAOnetTrees = make(map[string][]*onet.Tree)

	for i, n := range Nodes.All {

		IPlist := ""
		for IPaddr, exists := range n.IP {
			if exists {
				if IPlist != "" {
					IPlist = IPlist + "," + IPaddr
				} else {
					IPlist = IPaddr
				}
			}
		}

		AuxNodes.All[i] = CreateNode(n.Name, n.Level)
		//AuxNodes.All[i].AvailablePortsStart = n.AvailablePortsStart
		AuxNodes.All[i].ServerIdentity = n.ServerIdentity

		AuxNodes.All[i].ADist = make([]float64, 0)
		AuxNodes.All[i].PDist = make([]string, 0)
		AuxNodes.All[i].OptimalCluster = make(map[string]bool)
		AuxNodes.All[i].OptimalBunch = make(map[string]bool)
		AuxNodes.All[i].Cluster = make(map[string]bool)
		AuxNodes.All[i].Bunch = make(map[string]bool)
		AuxNodes.All[i].Rings = make([]string, 0)

	}

	for _, node := range AuxNodes.All {
		AuxNodes.ClusterBunchDistances[node] =
			make(map[*LocalityNode]float64)
		AuxNodes.Links[node] =
			make(map[*LocalityNode]map[*LocalityNode]bool)

		for _, node2 := range AuxNodes.All {
			AuxNodes.ClusterBunchDistances[node][node2] = math.MaxFloat64
			AuxNodes.Links[node][node2] = make(map[*LocalityNode]bool)

			if node == node2 {
				AuxNodes.ClusterBunchDistances[node][node2] = 0
			}
		}
	}

	for k, v := range Nodes.ServerIdentityToName {
		AuxNodes.ServerIdentityToName[k] = v
	}

	CreateLocalityGraph(AuxNodes, RND_NODES,
		RND_LVLS, NrLevels, PingDistances)

	// we are rooting trees here
	myname := rootNodeName

	if OPTIMIZED {
		OptimizeGraph(AuxNodes, myname, MIN_BUNCH_SIZE, OPTTYPE)
	}

	dist2 = AproximateDistanceOracle(AuxNodes)

	// Generate trees for all nodes
	for _, crtRoot := range AuxNodes.All {
		crtRootName := crtRoot.Name

		tree, NodesList, Parents, TreeRadiuses :=
			CreateOnetRings(AuxNodes, crtRootName, dist2)

		for i, n := range tree {
			ARATreeStruct[crtRootName] =
				append(ARATreeStruct[crtRootName], GraphTree{
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
				rosterNames =
					append(rosterNames, AuxNodes.GetServerIdentityToName(si))
			}

			log.Lvl3("generation node ", rootNodeName, "rootName ", rootName,
				"creates binary with roster", rosterNames)

			ARAOnetTrees[rootName] =
				append(ARAOnetTrees[rootName], createBinaryTreeFromGraphTree(n))
		}
	}

	log.Lvl3("done")

	return AuxNodes, dist2, ARATreeStruct, ARAOnetTrees
}

//Computes an onet binary tree based on a GraphTree
func createBinaryTreeFromGraphTree(GraphTree GraphTree) *onet.Tree {

	BinaryTreeRoster := GraphTree.Tree.Roster
	Tree := BinaryTreeRoster.GenerateBinaryTree()

	return Tree
}

// CreateNode with the given parameters
func CreateNode(Name string, level int) *LocalityNode {
	//func CreateNode(Name string, IP string, level int) *LocalityNode {

	var myNode LocalityNode

	/*
		myNode.X = x
		myNode.Y = y
	*/
	myNode.Name = Name
	myNode.IP = make(map[string]bool)

	/*
		tokens := strings.Split(IP, ",")
		for _, t := range tokens {
			myNode.IP[t] = true
		}
	*/
	myNode.Level = level
	myNode.ADist = make([]float64, 0)
	myNode.PDist = make([]string, 0)
	myNode.Cluster = make(map[string]bool)
	myNode.Bunch = make(map[string]bool)
	myNode.Rings = make([]string, 0)
	return &myNode
}
