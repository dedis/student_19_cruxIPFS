package service

import "time"

const (
	// LocalSim simulation is local
	LocalSim = true

	//DefaultIPFSAPIPort DefaultIPFSAPIPort
	DefaultIPFSAPIPort = 5001
	// DefaultIPFSGatewayPort DefaultIPFSGatewayPort
	DefaultIPFSGatewayPort = 8080
	// DefaultIPFSSwarmPort DefaultIPFSSwarmPort
	DefaultIPFSSwarmPort = 4001

	// BaseHostPort first port allocated to a node
	BaseHostPort = 14000

	// IPVersion default ip version
	IPVersion = "/ip4/"
	// TransportProtocol default transport protocol
	TransportProtocol = "/tcp/"

	// MaxPortNumberPerHost max number of ports that a host can use
	MaxPortNumberPerHost = 100
	// IPFSPortNumber number of ports used by an IPFS instance
	IPFSPortNumber = 3
	// ClusterPortNumber number of ports used by an ipfs cluster instance
	ClusterPortNumber = 3

	// DefaultReplMin ipfs cluster minimal replication factor
	DefaultReplMin = 2
	// DefaultReplMax ipfs cluster maximal replication factor
	DefaultReplMax = 3

	// IPFSStartupTime IPFSStartupTime
	IPFSStartupTime = 13 * time.Second
	// ClusterStartupTime ClusterStartupTime
	ClusterStartupTime = 2 * time.Second

	// ConfigsFolder folder name
	ConfigsFolder = "configs"
	// IPFSFolder ipfs config folder name
	IPFSFolder = "ipfs"
	// ClusterFolderPrefix prefix of cluster configs folder name
	ClusterFolderPrefix = "cluster-"

	// NodeName name of a node instance
	NodeName = "node_"
	// Node0 name of the first node
	Node0                = NodeName + "0"
	ClusterConsensusMode = "crdt"

	// WaitpeersName name of WaitPeers protocol
	WaitpeersName = "WaitPeers"
	// StartIPFSName name of StartIPFS protocol
	StartIPFSName = "StartIPFS"
	// ClusterBootstrapName name of ClusterBootstrap protocol
	ClusterBootstrapName = "ClusterBootstrap"

	// PingsFile File with stored pings
	PingsFile = "../pings.txt"
)
