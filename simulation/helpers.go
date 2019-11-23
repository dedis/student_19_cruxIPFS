package main

import (
	"strconv"

	"github.com/dedis/student_19_cruxIPFS/gentree"
)

// CreateNode with the given parameters
func CreateNode(Name string, x float64, y float64, IP string, level int) *gentree.LocalityNode {
	var myNode gentree.LocalityNode

	myNode.X = x
	myNode.Y = y
	myNode.Name = Name
	myNode.IP = make(map[string]bool)
	myNode.IP[IP] = true
	myNode.Level = level
	myNode.ADist = make([]float64, 0)
	myNode.PDist = make([]string, 0)
	myNode.Cluster = make(map[string]bool)
	myNode.Bunch = make(map[string]bool)
	myNode.Rings = make([]string, 0)
	return &myNode
}

// SetNodePaths set the node paths for remote and local node files
func SetNodePaths(n int) {
	NODEPATHREMOTE = NODEPATHNAME + strconv.Itoa(n) + ".txt"
	NODEPATHLOCAL = "../" + NODEPATHREMOTE
}
