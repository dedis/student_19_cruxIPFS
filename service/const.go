package service

import "time"

const (
	// ClusterConsensusMode "raft" or "crdt"
	ClusterConsensusMode = "raft"

	// DefaultReplMin ipfs cluster minimal replication factor
	DefaultReplMin = 2
	// DefaultReplMax ipfs cluster maximal replication factor
	DefaultReplMax = 3

	// BaseHostPort first port allocated to a node
	BaseHostPort = 14000

	// IPVersion default ip version
	IPVersion = "/ip4/"
	// TransportProtocol default transport protocol
	TransportProtocol = "/tcp/"

	// MaxPortNumberPerHost max number of ports that a host can use
	MaxPortNumberPerHost = 200
	// IPFSPortNumber number of ports used by an IPFS instance
	IPFSPortNumber = 3
	// ClusterPortNumber number of ports used by an ipfs cluster instance
	ClusterPortNumber = 3

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
	Node0 = NodeName + "0"

	// ServiceName name of the service
	ServiceName = "IPFS"

	// StartIPFSName name of StartIPFS protocol
	StartIPFSName = "StartIPFS"
	// ClusterBootstrapName name of ClusterBootstrap protocol
	ClusterBootstrapName = "ClusterBootstrap"
	// StartARAName name of the StartARA protocol
	StartARAName = "StartARA"
	// StartInstancesName name of the StartInstancesName protocol
	StartInstancesName = "StartInstances"

	// PingsFile File with stored pings
	PingsFile = "../pings.txt"
)
