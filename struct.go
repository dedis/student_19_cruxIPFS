package template

/*
This holds the messages used to communicate with the service over the network.
*/

import (
	"go.dedis.ch/onet/v3"
	"go.dedis.ch/onet/v3/network"
)

// We need to register all messages so the network knows how to handle them.
func init() {
	network.RegisterMessages(
		Count{}, CountReply{},
		Clock{}, ClockReply{},
		GenSecret{}, GenSecretReply{},
		StartIPFS{}, StartIPFSReply{},
		StartCluster{}, StartClusterReply{},
	)
}

const (
	// ErrorParse indicates an error while parsing the protobuf-file.
	ErrorParse = iota + 4000
)

// Clock will run the tepmlate-protocol on the roster and return
// the time spent doing so.
type Clock struct {
	Roster *onet.Roster
}

// ClockReply returns the time spent for the protocol-run.
type ClockReply struct {
	Time     float64
	Children int
}

// Count will return how many times the protocol has been run.
type Count struct {
}

// CountReply returns the number of protocol-runs
type CountReply struct {
	Count int
}

// GenSecret will return a new shared secret
type GenSecret struct {
}

// GenSecretReply reply of gensecret containing the secret as a string
type GenSecretReply struct {
	Secret string
}

// StartIPFS ipfs start packet
type StartIPFS struct {
	ConfigPath string
	NodeID     int
	PortMin    int
	PortMax    int
	IP         string
}

// IPFSPortN number of ports IPFS is using
const IPFSPortN int = 3

// IPFSSwarmID id of the Swarm port
const IPFSSwarmID int = 0

// IPFSAPIID id of the API port
const IPFSAPIID int = 1

// IPFSGatewayID if of the Gateway port
const IPFSGatewayID int = 2

// IPFSPorts structure containing all ports used by IPFS
type IPFSPorts struct {
	Swarm   int
	API     int
	Gateway int
}

// DefaultConfigPath default path to config files
const DefaultConfigPath string = "/home/guillaume/ipfs_test/myfolder"

// StartIPFSReply ss
type StartIPFSReply struct {
	Ports *IPFSPorts
}

// ClusterPortN number of ports needed by the cluster
const ClusterPortN int = 3

// StartCluster packet that is sent to start a cluster instance
type StartCluster struct {
	ConfigPath  string
	NodeID      int
	ClusterID   int
	PortMin     int
	PortMax     int
	IPFSAPIPort int
	Peername    string
	Secret      string
}

// ClusterPorts ports that are used by the cluster
type ClusterPorts struct {
	IPFSAPI   int
	RestAPI   int
	IPFSProxy int
	Cluster   int
}

// StartClusterReply reply sent once that the cluster instance has started
type StartClusterReply struct {
	Ports *ClusterPorts
}
