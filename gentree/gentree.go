package gentree

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"math"
	"math/rand"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/satori/go.uuid"
	"go.dedis.ch/onet/v3"
	"go.dedis.ch/onet/v3/log"
	"go.dedis.ch/onet/v3/network"
)

// ToRecursiveTreeNode finds the equivalent tree node in the recursive tree.
func (t *TreeConverter) ToRecursiveTreeNode(target *onet.TreeNode) (*onet.TreeNode, error) {
	return findTreeNode(t.RecursiveTree, target)
}

// ToBinaryTreeNode finds the equivalent tree node in the binary tree.
func (t *TreeConverter) ToBinaryTreeNode(target *onet.TreeNode) (*onet.TreeNode, error) {
	return findTreeNode(t.BinaryTree, target)
}

// GetByIP gets the node by IP.
func (ns LocalityNodes) GetByIP(ip string) *LocalityNode {

	for _, n := range ns.All {
		if n.IP[ip] {
			return n
		}
	}
	return nil
}

func (ns LocalityNodes) GetByServerIdentityIP(ip string) *LocalityNode {

	for _, n := range ns.All {
		//if strings.Contains(n.ServerIdentity.String(), ip) {
		if n.IP[ip] {
			return n
		}
	}
	return nil
}

/*
func (ns LocalityNodes) OccupyNextPort(ip string) int {

	port := -1
	for _, n := range ns.All {
		if strings.Contains(n.ServerIdentity.String(), ip) {
			n.NextPortMtx.Lock()
			if n.NextPort != n.AvailablePortsEnd {
				port = n.NextPort
				n.NextPort++
			}
			n.NextPortMtx.Unlock()
		}
	}
	return port
}
*/

/*
func (ns LocalityNodes) OccupyNextPortByName(name string) int {

	port := -1
	for _, n := range ns.All {
		if n.Name == name {
			if n.NextPort != n.AvailablePortsEnd {
				port = n.NextPort
				n.NextPort++
			}
		}
	}
	return port
}
*/

// GetByName gets the node by name.
func (ns LocalityNodes) GetByName(name string) *LocalityNode {
	nodeIdx := NodeNameToInt(name)

	//log.LLvl1("name here is", name)

	//log.LLvl1("ns length", len(ns.All), "nodeIdx", nodeIdx)
	if len(ns.All) < nodeIdx {
		//log.LLvl1("returning NOT fine")
		return nil
	}
	//log.LLvl1("returning fine", ns.All[nodeIdx])
	//log.LLvl1(ns.All)
	return ns.All[nodeIdx%len(ns.All)]
	//return ns.All[nodeIdx]
}

// NameToServerIdentity gets the server identity by name.
func (ns LocalityNodes) NameToServerIdentity(name string) *network.ServerIdentity {
	node := ns.GetByName(name)

	if node != nil {

		if node.ServerIdentity == nil {
			log.Error("nil 1", node, node.ServerIdentity)
		}

		if node.ServerIdentity.Address == "" {
			log.Error("nil 2", node.ServerIdentity.Address)
		}

		return node.ServerIdentity
	}

	return nil
}

// ServerIdentityToName gets the name by server identity.
func (ns LocalityNodes) GetServerIdentityToName(sid *network.ServerIdentity) string {
	return ns.ServerIdentityToName[sid.ID]
}

func findTreeNode(tree *onet.Tree, target *onet.TreeNode) (*onet.TreeNode, error) {

	for _, node := range tree.List() {
		if node.ID.Equal(target.ID) {
			return node, nil
		}
	}
	return nil, errors.New("not found")
}

//First argument is all the nodes
//Second argument is if the coordinates are random
//Third arguments is if the levels are random
//Fourth argument is how many levels there should be if they are random
//Second and Third arguments should always be the same
func CreateLocalityGraph(all LocalityNodes, randomCoords, randomLevels bool, levels int, pingDist map[string]map[string]float64) {
	nodes := all.All

	randSrc := rand.New(rand.NewSource(time.Now().UnixNano()))

	/*
		if randomCoords {
			//Computes random coordinates
			for _, n := range nodes {
				n.X = randSrc.Float64() * 500
				n.Y = randSrc.Float64() * 500
			}
		}
	*/

	if randomLevels {
		//Computes random levels
		probability := 1.0 / math.Pow(float64(len(nodes)), 1.0/float64(levels))
		for i := 1; i < levels; i++ {
			for _, n := range nodes {
				if n.Level == i-1 && randSrc.Float64() < probability {
					n.Level = i
				}
			}
		}
	}

	for i := 0; i < levels; i++ {

		for _, v1 := range nodes {

			var node *LocalityNode
			var distance float64
			min := math.MaxFloat64

			for _, v2 := range nodes {

				if v2.Level >= i {

					distance = pingDist[v1.Name][v2.Name]

					if distance < min {
						min = distance
						node = v2
					}
				}
			}
			v1.PDist = append(v1.PDist, node.Name)
			v1.ADist = append(v1.ADist, min)
		}
	}

	for _, v1 := range nodes {

		for i := 0; i < levels-1; i++ {

			for _, v2 := range nodes {

				d := pingDist[v1.Name][v2.Name]

				//log.LLvl1("assigning", v1.Name, v2.Name)

				if v2.Level >= i {

					if checkDistance(d, i, levels, v1.ADist) && v2 != v1 {

						v1.Bunch[v2.Name] = true
						v1.OptimalBunch[v2.Name] = true
						all.ClusterBunchDistances[v1][v2] = d
						all.ClusterBunchDistances[v2][v1] = d

					}

				}

				if v2.Level == levels-1 && v2 != v1 {

					v1.Bunch[v2.Name] = true
					v1.OptimalBunch[v2.Name] = true
					all.ClusterBunchDistances[v1][v2] = d
					all.ClusterBunchDistances[v2][v1] = d
				}
			}
		}
	}

	for _, v1 := range nodes {

		for _, v2 := range nodes {

			if v2.Bunch[v1.Name] && v2.Name != v1.Name {

				v1.OptimalCluster[v2.Name] = true
				v1.Cluster[v2.Name] = true
				all.ClusterBunchDistances[v1][v2] = pingDist[v1.Name][v2.Name]
				all.ClusterBunchDistances[v2][v1] = pingDist[v1.Name][v2.Name]
			}
		}
	}
	/*

			// write to file
			file, _ := os.Create("Specs/original.txt")
			w := bufio.NewWriter(file)
			w.WriteString(strconv.Itoa(len(all.All)) + "\n")
			for _, node := range all.All {
				w.WriteString(fmt.Sprint(node.X) + " " + fmt.Sprint(node.Y) + "\n")

			}

			for _, node := range all.All {
				for clusterNodeName, exists := range node.Cluster {
					if exists {
						name1 := strings.Split(node.Name, "_")[1]
						name2 := strings.Split(clusterNodeName, "_")[1]

						w.WriteString(name1 + " " + name2 + "\n")
					}
				}
			}

			w.Flush()
			file.Close()

				file, _ = os.Create("nodes_read.txt")
				w = bufio.NewWriter(file)
					for _, node := range all.All {
						w.WriteString(node.Name + " " + fmt.Sprint(node.X) + "," + fmt.Sprint(node.Y) + " 127.0.0.1 " + strconv.Itoa(node.Level) + "\n")

					}

				w.Flush()
		file.Close()
	*/

}

//Checks if a Node is suitable to be another Node's bunch depending on its distance to it
func checkDistance(distance float64, lvl int, lvls int, Adist []float64) bool {

	for i := lvl + 1; i < lvls; i++ {

		if distance <= Adist[i] {
			// do nothing
		} else {
			return false
		}
	}
	return true
}

/*
//Computes the Euclidian distance between two nodes
func ComputeDist(v1 *LocalityNode, v2 *LocalityNode, pingDist map[string]map[string]float64) float64 {
	if len(pingDist) == 0 {
		//panic("aaa")
		dist := math.Sqrt(math.Pow(v1.X-v2.X, 2) + math.Pow(v1.Y-v2.Y, 2))
		return dist
	}
	return pingDist[v1.Name][v2.Name]
}
*/

func GenerateRadius(maxDist float64) []float64 {
	multiplier := 1.0
	//base := math.Sqrt(2)
	base := 2.0
	radiuses := make([]float64, 0)
	for i := 0; ; i++ {
		crtMaxRadius := multiplier * math.Pow(base, float64(i))
		/*
			prevRadius := 0.0
			if i != 0 {
				prevRadius = multiplier * math.Pow(base, float64(i - 1))
			}
			if crtMaxRadius > maxDist && prevRadius > maxDist {
		*/
		if crtMaxRadius > 256 {
			break
		}
		radiuses = append(radiuses, crtMaxRadius)
	}

	radiuses = append(radiuses, maxDist)

	return radiuses
}

func generateRingID(node string, ringNumber int) string {
	return node + "_" + strconv.Itoa(ringNumber)
}

func getRingIDFromDistance(distance float64) int {
	radiuses := GenerateRadius(distance)
	return len(radiuses) - 1
}

// Called on a node, It will add all the coresponding children depending on the optimisation stated previously
// AllowedNodes are the nodes that remain in the tree after the optimisation and the filter by radius (Rings)
// treeNode is the node that we are about to set the parents/childrens of
func CreateAndSetChildren(Rings bool, AllowedNodes map[string]bool, file *os.File, all LocalityNodes, treeNode *onet.TreeNode, NodeList map[string]*onet.TreeNode, parents map[*onet.TreeNode][]*onet.TreeNode) map[*onet.TreeNode][]*onet.TreeNode {

	//Ranges through the nodes Cluster
	for clusterNodeName, exists := range all.All[treeNode.RosterIndex].OptimalCluster {

		//Continues if the node is not in the range of the radius if we have activated the Ring mode
		if Rings && !AllowedNodes[clusterNodeName] {
			continue
		}
		//Continues if the node doesnt exist in the tree
		if !exists {

			continue

		}

		var clusterTreeNode *onet.TreeNode

		// check if clusternode is not already a tree node and creates it if doesn't exist yet
		if clusterAuxTreeNode, ok := NodeList[all.NameToServerIdentity(clusterNodeName).String()]; !ok {

			ClusterNodeID := all.NameToServerIdentity(clusterNodeName)
			clusterTreeNode = onet.NewTreeNode(NodeNameToInt(clusterNodeName), ClusterNodeID)
			clusterTreeNode.Parent = treeNode
			NodeList[clusterTreeNode.ServerIdentity.String()] = clusterTreeNode

		} else {

			clusterTreeNode = clusterAuxTreeNode
		}

		childExists := false

		//Checks if the node that is about to be added as a children is already a children
		for _, child := range treeNode.Children {

			if child.RosterIndex == clusterTreeNode.RosterIndex {

				childExists = true
				break
			}
		}

		//Adds node as a children
		if !childExists {

			nodeName := all.GetServerIdentityToName(treeNode.ServerIdentity)

			clusterTreeNodeName := all.GetServerIdentityToName(clusterTreeNode.ServerIdentity)

			fmt.Fprintf(file, strconv.Itoa(NodeNameToInt(nodeName))+" "+strconv.Itoa(NodeNameToInt(clusterTreeNodeName))+"\n")

			treeNode.Children = append(treeNode.Children, clusterTreeNode)

		}
	}

	//Does the same as the first loop but with the Bunch
	for bunchNodeName, exists := range all.All[treeNode.RosterIndex].OptimalBunch {

		if Rings && !AllowedNodes[bunchNodeName] {
			continue
		}

		if !exists {

			continue

		}

		var bunchTreeNode *onet.TreeNode

		// check if clusternode is not already a tree node
		if bunchAuxTreeNode, ok := NodeList[all.NameToServerIdentity(bunchNodeName).String()]; !ok {
			BunchNodeID := all.NameToServerIdentity(bunchNodeName)
			bunchTreeNode = onet.NewTreeNode(NodeNameToInt(bunchNodeName), BunchNodeID)
			bunchTreeNode.Parent = treeNode
			NodeList[bunchTreeNode.ServerIdentity.String()] = bunchTreeNode
			parents[treeNode] = append(parents[treeNode], NodeList[all.NameToServerIdentity(bunchNodeName).String()])

		} else {
			bunchTreeNode = bunchAuxTreeNode
			parents[treeNode] = append(parents[treeNode], NodeList[all.NameToServerIdentity(bunchNodeName).String()])
		}

		childExists := false
		for _, child := range treeNode.Children {
			if child.RosterIndex == bunchTreeNode.RosterIndex {
				childExists = true
				break
			}
		}

		if !childExists {

			treeNode.Children = append(treeNode.Children, bunchTreeNode)

		}
	}

	return parents
}

//Converts a Node to it's index
func NodeNameToInt(nodeName string) int {
	separation := strings.Split(nodeName, "_")
	if len(separation) != 2 {
		log.LLvl1(separation)
	}
	idx, err := strconv.Atoi(separation[1])
	if err != nil {
		panic(err.Error())
	}
	return idx
}

//Filters Nodes depending on their distance to the root given as an argument
//distances is the distance between two nodes in the current graph
func Filter(all LocalityNodes, root *LocalityNode, radius float64, distances map[*LocalityNode]map[*LocalityNode]float64) []*LocalityNode {

	//Childrend That are in the radius
	ChildrenInRange := make(map[*LocalityNode]bool)

	StartPointEndPoint := make(map[*LocalityNode][]*LocalityNode)

	//Node used to go from node A to node B if needed
	MiddlePoint := make(map[*LocalityNode]map[*LocalityNode]*LocalityNode)

	MiddlePoint[root] = make(map[*LocalityNode]*LocalityNode)

	for _, n := range all.All {

		MiddlePoint[n] = make(map[*LocalityNode]*LocalityNode)
	}

	for n, _ := range root.Cluster {

		//checks if node is inside the radius
		if distances[root][all.GetByName(n)] <= radius {
			//Adds it to the final nodes if it is inside of the radius
			ChildrenInRange[all.GetByName(n)] = true

			//ranges through the links that connect the root to node n if there are any
			for k, _ := range all.Links[root][all.GetByName(n)] {

				//Adds them to the final nodes
				ChildrenInRange[k] = true
				//Adds n as present indirect connection from the root to n
				StartPointEndPoint[root] = append(StartPointEndPoint[root], all.GetByName(n))
				//Adds the link as the middle point between the root and node n
				MiddlePoint[root][all.GetByName(n)] = k

			}

		}

	}

	for {

		CopyStartPointEndPoint := make(map[*LocalityNode][]*LocalityNode)
		CopyMiddlePoint := make(map[*LocalityNode]map[*LocalityNode]*LocalityNode)
		CopyMiddlePoint[root] = make(map[*LocalityNode]*LocalityNode)

		i := 0
		j := 0

		//Ranges through all the nodes that are currently used in a path and that are not directly connected to the root to check for possible links
		for Root, EndPoints := range StartPointEndPoint {

			//Ranges throught the nodes that are not directly connected to the root
			for _, EndPoint := range EndPoints {

				j++

				if len(all.Links[Root][MiddlePoint[root][EndPoint]]) == 0 {
					i++
				}

				//Ranges through the links that connect the root to the node that is currently the furthest from the node in the cluster that is not directly connected to the root
				//In other words you can imagine this as two nodes wich are not directly connecte, the root and another node, and we try to reconstruct the path starting by the node that is not the root
				//So endpoint would represent the second that is the furthest starting from the node that is note the root, and middlepoint of that node and the root would be the node that is the
				// furthes from the node that is not the root, the last node of the path that is being built towards the root
				//In this loop we range through that last node and the root links to see what nodes we can use to connect them if there is any
				for a, _ := range all.Links[Root][MiddlePoint[root][EndPoint]] {
					ChildrenInRange[a] = true
					CopyStartPointEndPoint[Root] = append(CopyStartPointEndPoint[Root], MiddlePoint[root][EndPoint])
					MiddlePoint[Root][MiddlePoint[root][EndPoint]] = a

				}
				if len(all.Links[MiddlePoint[root][EndPoint]][EndPoint]) == 0 {
					i++
				}
				//In this loop we do the same as the previous loop but instead of starting from the node that is not the root we start from the root
				for a, _ := range all.Links[MiddlePoint[root][EndPoint]][EndPoint] {

					ChildrenInRange[a] = true
					CopyStartPointEndPoint[MiddlePoint[root][EndPoint]] = append(CopyStartPointEndPoint[MiddlePoint[root][EndPoint]], EndPoint)
					MiddlePoint[MiddlePoint[root][EndPoint]][EndPoint] = a

				}

			}

		}
		//It breaks when every time links is empty, meaning the path is completed
		if i == j*2 {
			break
		}

		StartPointEndPoint = CopyStartPointEndPoint
		MiddlePoint = CopyMiddlePoint
	}

	FinalNodes := ChildrenInRange

	FinalNodes[root] = true

	Nodes := make([]*LocalityNode, 0)

	for k, _ := range FinalNodes {

		Nodes = append(Nodes, k)
	}

	return Nodes
}

func NodesInARA(all LocalityNodes, root *LocalityNode, radius float64, distances map[*LocalityNode]map[*LocalityNode]float64) []*LocalityNode {

	//Childrend That are in the radius
	ChildrenInRange := make(map[*LocalityNode]bool)

	for n, _ := range root.Cluster {

		//checks if node is inside the radius
		if distances[root][all.GetByName(n)] <= radius {
			//Adds it to the final nodes if it is inside of the radius
			ChildrenInRange[all.GetByName(n)] = true
		}

	}

	FinalNodes := ChildrenInRange
	FinalNodes[root] = true

	Nodes := make([]*LocalityNode, 0)

	for k, _ := range FinalNodes {
		Nodes = append(Nodes, k)
	}

	return Nodes
}

//Root is the root of the graph
//Optimization is the upperBound set on the bunch of each node of the graph
func OptimizeGraph(all LocalityNodes, rootName string, Optimization int, OptType int) {
	RemoveLinks(all, all.GetByName(rootName), Optimization, OptType)
}

// CreateOnetLPTree TODO add documentation
// Will Build The Tree calling different functions for different purposes
// It's the main function, all functions created are called directly or indirectly through this one
func CreateOnetLPTree(all LocalityNodes, rootName string, BunchLowerBound int) ([]*onet.Tree, [][]*onet.TreeNode, []map[*onet.TreeNode][]*onet.TreeNode, map[*LocalityNode]map[*LocalityNode]float64) {

	//Slice of Trees to be returned
	Trees := make([]*onet.Tree, 0)

	//Slice of Lists of all the nodes for each tree
	Lists := make([][]*onet.TreeNode, 0)

	//Slice of Parents for each node od each Tree
	Parents := make([]map[*onet.TreeNode][]*onet.TreeNode, 0)

	//Creates a file where we can write
	file, _ := os.Create("Specs/optimized.txt")
	fmt.Fprintf(file, strconv.Itoa(len(all.All))+"\n")

	/*
		//Prints coordinates of all nodes into the file
		for _, n := range all.All {

			fmt.Fprintf(file, fmt.Sprint(n.X)+" "+fmt.Sprint(n.Y)+"\n")
		}
	*/

	//Distance between nodes after Optimisation
	Dist2 := AproximateDistanceOracle(all)

	// AllowedNodes are the nodes that remain in the tree after the optimisation and the filter by radius (Rings)
	AllowedNodes := make(map[string]bool)

	parents := make(map[*onet.TreeNode][]*onet.TreeNode)

	var rootIdxInRoster int
	nrProcessedNodes := 0

	//string represents ServerIdentity
	nodesInTree := make(map[string]*onet.TreeNode)
	roster := make([]*network.ServerIdentity, 0)

	var rootID *network.ServerIdentity

	// create root node
	rootID = all.NameToServerIdentity(rootName)
	root := onet.NewTreeNode(NodeNameToInt(rootName), rootID)
	nodesInTree[root.ServerIdentity.String()] = root

	//Creates Roster*
	for i, k := range all.All {

		if k.Name == rootName {
			rootIdxInRoster = i
		}
		roster = append(roster, all.NameToServerIdentity(k.Name))
	}

	//Remplacing firts element of Roster by the root
	replacement := roster[0]
	roster[0] = rootID
	roster[rootIdxInRoster] = replacement

	// set root children
	parents = CreateAndSetChildren(false, AllowedNodes, file, all, root, nodesInTree, parents)

	nrProcessedNodes++

	nextLevelNodes := make(map[*onet.TreeNode]bool)
	nextLevelNodesAux := make(map[*onet.TreeNode]bool)

	for _, childNode := range root.Children {

		parents = CreateAndSetChildren(false, AllowedNodes, file, all, childNode, nodesInTree, parents)

		nrProcessedNodes++

		for _, child := range childNode.Children {
			nextLevelNodesAux[child] = true
		}
	}

	nextLevelNodes = nextLevelNodesAux
	for nrProcessedNodes <= len(all.All) {
		if len(root.Children) == 0 {
			break
		}
		for childNode := range nextLevelNodes {

			parents = CreateAndSetChildren(false, AllowedNodes, file, all, childNode, nodesInTree, parents)
			nrProcessedNodes++
			for _, child := range childNode.Children {

				nextLevelNodesAux[child] = true
			}
		}
		nextLevelNodes = nextLevelNodesAux
	}

	//Computes Final Roster
	finalRoster := onet.NewRoster(roster)

	//Hashing
	h := sha256.New()
	for _, r := range roster {
		r.Public.MarshalTo(h)
	}

	url := network.NamespaceURL + "tree/" + finalRoster.ID.String() + hex.EncodeToString(h.Sum(nil))

	//Creates Tree
	t := &onet.Tree{
		Roster: finalRoster,
		Root:   root,
		ID:     onet.TreeID(uuid.NewV5(uuid.NamespaceURL, url)),
	}

	//Creates List of Nodes In Tree
	list := make([]*onet.TreeNode, 0)

	for _, v := range nodesInTree {

		list = append(list, v)

	}

	file.Close()

	Lists = append(Lists, list)
	Trees = append(Trees, t)
	Parents = append(Parents, parents)

	return Trees, Lists, Parents, Dist2
}

type ByServerIdentityAlphabetical []*network.ServerIdentity

func (a ByServerIdentityAlphabetical) Len() int {
	return len(a)
}

func (a ByServerIdentityAlphabetical) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a ByServerIdentityAlphabetical) Less(i, j int) bool {
	return a[i].String() < a[j].String()
}

func CreateARAs(all LocalityNodes, rootName string, dist2 map[*LocalityNode]map[*LocalityNode]float64) ([][]*LocalityNode, []float64) {

	//Slice of Lists of all the nodes for each tree
	Lists := make([][]*LocalityNode, 0)

	TreeRadiuses := make([]float64, 0)

	radiuses := GenerateRadius(10000)

	prevARASize := 0

	countt := 0

	for {

		AllowedNodes := make(map[string]bool)

		var Filterr []*LocalityNode

		Filterr = NodesInARA(all, all.GetByName(rootName), radiuses[countt], dist2)

		for _, n := range Filterr {
			AllowedNodes[n.Name] = true
		}

		for name, allowed := range AllowedNodes {
			if allowed && !all.GetByName(rootName).Cluster[name] && name != rootName {
				log.Panic("adding node to ring that is not in the cluster!", "root", rootName, "child", name, "but cluster is", all.GetByName(rootName).Cluster)
			}
		}

		if len(AllowedNodes) == prevARASize {
			countt++
			log.Lvl2("it's a repeat tree, skip it!")
			if countt == len(radiuses) {
				break
			}
			continue
		}

		prevARASize = len(AllowedNodes)

		if prevARASize == 1 {

			log.Lvl3("Only Root in tree number:", countt)
		}

		Lists = append(Lists, Filterr)
		TreeRadiuses = append(TreeRadiuses, radiuses[countt])

		countt++
		if countt == len(radiuses) {
			break
		}
	}

	return Lists, TreeRadiuses
}

func CreateOnetRings(all LocalityNodes, rootName string, dist2 map[*LocalityNode]map[*LocalityNode]float64) ([]*onet.Tree, [][]*onet.TreeNode, []map[*onet.TreeNode][]*onet.TreeNode, []float64) {

	//Slice of Trees to be returned
	Trees := make([]*onet.Tree, 0)

	//Slice of Lists of all the nodes for each tree
	Lists := make([][]*onet.TreeNode, 0)

	//Slice of Parents for each node od each Tree
	Parents := make([]map[*onet.TreeNode][]*onet.TreeNode, 0)

	TreeRadiuses := make([]float64, 0)

	radiuses := GenerateRadius(10000)

	//Distance between nodes after Optimisation

	/*
		var Links map[*LocalityNode]map[*LocalityNode]map[*LocalityNode]bool
		Dist2, _ := AproximateDistanceOracle(all)
		//Returns maps of link nodes between two nodes ([NodeA][NodeB][NodeC]Returns True if NodeC is a link between NodeA and NodeB )
		_, Links = AproximateDistanceOracle(all)
	*/

	prevRosterLen := 0

	countt := 0

	for {

		//Creates a file where we can write
		file, _ := os.Create("Specs/result" + strconv.Itoa(countt) + ".txt")
		/*

			fmt.Fprintf(file, strconv.Itoa(len(all.All))+"\n")

			//Prints coordinates of all nodes into the file
			for _, n := range all.All {

				fmt.Fprintf(file, fmt.Sprint(n.X)+" "+fmt.Sprint(n.Y)+"\n")
			}
		*/

		AllowedNodes := make(map[string]bool)
		var Filterr []*LocalityNode

		Filterr = Filter(all, all.GetByName(rootName), radiuses[countt], dist2)

		for _, n := range Filterr {
			AllowedNodes[n.Name] = true
		}

		for name, allowed := range AllowedNodes {
			if allowed && !all.GetByName(rootName).Cluster[name] && name != rootName {
				log.Panic("adding node to ring that is not in the cluster!", "root", rootName, "child", name, "but cluster is", all.GetByName(rootName).Cluster)
			}
		}

		parents := make(map[*onet.TreeNode][]*onet.TreeNode)

		var rootIdxInRoster int
		nrProcessedNodes := 0

		//string represents ServerIdentity
		nodesInTree := make(map[string]*onet.TreeNode)
		roster := make([]*network.ServerIdentity, 0)

		var rootID *network.ServerIdentity

		// create root node
		//log.LLvl1("in create rings root name is", rootName)
		rootID = all.NameToServerIdentity(rootName)
		//log.LLvl1("root id is", rootID)
		//log.LLvl1("rootid", rootID, "name to int", NodeNameToInt(rootName))

		root := onet.NewTreeNode(NodeNameToInt(rootName), rootID)
		nodesInTree[root.ServerIdentity.String()] = root

		//Creates Roster*

		roster = make([]*network.ServerIdentity, 0)
		for i, n := range Filterr {
			if n.Name == rootName {
				rootIdxInRoster = i
			}

			roster = append(roster, n.ServerIdentity)

		}

		//Remplacing firts element of Roster by the root
		replacement := roster[0]
		roster[0] = rootID
		roster[rootIdxInRoster] = replacement

		// set root children
		parents = CreateAndSetChildren(true, AllowedNodes, file, all, root, nodesInTree, parents)

		nrProcessedNodes++

		nextLevelNodes := make(map[*onet.TreeNode]bool)
		nextLevelNodesAux := make(map[*onet.TreeNode]bool)

		for _, childNode := range root.Children {

			parents = CreateAndSetChildren(true, AllowedNodes, file, all, childNode, nodesInTree, parents)

			nrProcessedNodes++

			for _, child := range childNode.Children {
				nextLevelNodesAux[child] = true
			}
		}

		nextLevelNodes = nextLevelNodesAux
		for nrProcessedNodes <= len(all.All) {
			if len(root.Children) == 0 {
				break
			}
			for childNode := range nextLevelNodes {

				parents = CreateAndSetChildren(true, AllowedNodes, file, all, childNode, nodesInTree, parents)
				nrProcessedNodes++
				for _, child := range childNode.Children {

					nextLevelNodesAux[child] = true
				}
			}
			nextLevelNodes = nextLevelNodesAux
		}

		//Computes Final Roster

		// deterministic roster order
		// the roster order does not affect the locality graph, because the locality graph is built using Parents
		if len(roster) > 1 {

			// put in rosterAux all elements of roster except the root, which is at index 0
			rosterAux := make([]*network.ServerIdentity, 0)
			for i, x := range roster {
				if i == 0 {
					continue
				}
				rosterAux = append(rosterAux, x)
			}

			// sort rosterAux
			sort.Sort(ByServerIdentityAlphabetical(rosterAux))

			// add the sorted elements in roster after the root
			for i := range roster {
				if i == 0 {
					continue
				}
				roster[i] = rosterAux[i-1]
			}
		}

		finalRoster := onet.NewRoster(roster)

		if len(finalRoster.List) == prevRosterLen {
			countt++
			log.Lvl2("it's a repeat tree, skip it!")
			if countt == len(radiuses) {
				break
			}
			continue
		}

		prevRosterLen = len(finalRoster.List)

		//log.LLvl1(finalRoster)

		//Hashing
		h := sha256.New()
		for _, r := range roster {
			r.Public.MarshalTo(h)
		}

		url := network.NamespaceURL + "tree/" + finalRoster.ID.String() + hex.EncodeToString(h.Sum(nil))

		//Creates Tree
		t := &onet.Tree{
			Roster: finalRoster,
			Root:   root,
			ID:     onet.TreeID(uuid.NewV5(uuid.NamespaceURL, url)),
		}

		//Creates List of Nodes In Tree
		list := make([]*onet.TreeNode, 0)

		for _, v := range nodesInTree {

			list = append(list, v)

		}

		file.Close()

		if len(list) == 1 {

			log.Lvl3("Only Root in tree number:", countt)
		}

		Lists = append(Lists, list)
		Trees = append(Trees, t)
		Parents = append(Parents, parents)
		TreeRadiuses = append(TreeRadiuses, radiuses[countt])

		countt++
		if countt == len(radiuses) {
			break
		}
	}

	return Trees, Lists, Parents, TreeRadiuses
}
