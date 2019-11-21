package ARAgen

import (
	"go.dedis.ch/onet/v3"
)

// GraphTree Represents the actual graph that will be linked to the Binary
// Tree of the Protocol
type GraphTree struct {
	Tree        *onet.Tree
	ListOfNodes []*onet.TreeNode
	Parents     map[*onet.TreeNode][]*onet.TreeNode
	Radius      float64
}
