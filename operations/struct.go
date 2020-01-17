package operations

import (
	"github.com/ipfs/ipfs-cluster/api/rest/client"
)

// Node with its ipfs cluster ac
type Node struct {
	Name    string
	Clients []client.Client
	Secrets []string
	Addrs   []string
}
