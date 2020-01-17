package service

import (
	"bufio"
	"os"
	"sync"

	"github.com/dedis/student_19_cruxIPFS/gentree"

	"go.dedis.ch/onet/v3"
	"go.dedis.ch/onet/v3/network"
	"go.dedis.ch/onet/v3/simul/monitor"
)

// storageID reflects the data we're storing - we could store more
// than one structure.
var storageID = []byte("main")

// Service is our template-service
type Service struct {
	// We need to embed the ServiceProcessor, so that incoming messages
	// are correctly handled.
	*onet.ServiceProcessor

	storage *storage

	Nodes        gentree.LocalityNodes
	LocalityTree *onet.Tree
	Parents      []*onet.TreeNode
	GraphTree    map[string][]gentree.GraphTree
	BinaryTree   map[string][]*onet.Tree
	alive        bool
	Distances    map[*gentree.LocalityNode]map[*gentree.LocalityNode]float64

	PortMutex    *sync.Mutex
	W            *bufio.Writer
	File         *os.File
	metrics      map[string]*monitor.TimeMeasure
	metricsMutex sync.Mutex

	BandwidthRx uint64
	BandwidthTx uint64
	NrMsgRx     uint64
	NrMsgTx     uint64

	NrProtocolsStarted uint64

	OwnPings      map[string]float64
	DonePing      bool
	PingDistances map[string]map[string]float64
	NrPingAnswers int
	PingAnswerMtx sync.Mutex
	PingMapMtx    sync.Mutex

	Name       string // name of the service (node_2)
	MyIP       string // IP address
	ConfigPath string // path to home config folder
	MyIPFSPath string // path to ipfs config folder of that service
	MinPort    int    // port range allocated to this node
	MaxPort    int
	MyIPFS     []IPFSInformation          // own ipfs information
	OtherIPFS  map[string]IPFSInformation // node_x -> IP, ports etc.

	OnetTree      *onet.Tree
	StartIPFSProt onet.ProtocolInstance
}

// FnService function that returns the service of that node
type FnService func() *Service

// storage is used to save our data.
type storage struct {
	Count int
	sync.Mutex
}

// ClusterInfo information about a given cluster
type ClusterInfo struct {
	Leader    string
	Secret    string
	Size      int
	Instances []ClusterInstance
}

// ClusterInstance details of a cluster
type ClusterInstance struct {
	HostName      string
	IP            string
	IPFSAPIAddr   string
	RestAPIPort   int
	IPFSProxyPort int
	ClusterPort   int
}

// IPFSInformation structure containing information about an IPFS instance
type IPFSInformation struct {
	Name        string
	IP          string
	SwarmPort   int
	APIPort     int
	GatewayPort int
}

// InitRequest packet
type InitRequest struct {
	Nodes                []*gentree.LocalityNode
	ServerIdentityToName map[*network.ServerIdentity]string
	OnetTree             *onet.Tree
	Roster               *onet.Roster
	Cruxified            bool
	ComputePings         bool
	Mode                 string
}

// InitResponse packet
type InitResponse struct {
}

// ReqPings request packet for ping service
type ReqPings struct {
	SenderName string
}

// ReplyPings reply packet for ping service
type ReplyPings struct {
	Pings      string
	SenderName string
}

// NodeInfo contains ipfs and ipfs-cluster information for a given node
type NodeInfo struct {
	IPFS     IPFSInformation
	Clusters []ClusterInfo
}

// StartIPFSProtocol structure
type StartIPFSProtocol struct {
	*onet.TreeNodeInstance
	announceChan chan announceWrapperStartIPFS
	repliesChan  chan []replyWrapperStartIPFS
	Ready        chan bool
	GetService   FnService
	Nodes        map[string]*NodeInfo
}

// StartIPFSAnnounce is used to pass a message to all children.
type StartIPFSAnnounce struct {
	Message string
}

// announceWrapperWaitpeers just contains Announce and the data necessary to
// identify and process the message in onet.
type announceWrapperStartIPFS struct {
	*onet.TreeNode
	StartIPFSAnnounce
}

// StartIPFSReply returns true when ready.
type StartIPFSReply struct {
	IPFS     *IPFSInformation
	Clusters *[]ClusterInfo
}

// replyWrapper just contains Reply and the data necessary to identify and
// process the message in onet.
type replyWrapperStartIPFS struct {
	*onet.TreeNode
	StartIPFSReply
}

// ClusterBootstrapProtocol structure
type ClusterBootstrapProtocol struct {
	*onet.TreeNodeInstance
	announceChan chan announceWrapperClusterBootstrap
	repliesChan  chan []replyWrapperClusterBootstrap
	Ready        chan bool
	Info         ClusterInfo
	GetService   FnService
}

// ClusterBootstrapAnnounce is used to pass a message to all children.
type ClusterBootstrapAnnounce struct {
	SenderName string
	Bootstrap  string
	Secret     string
}

// announceWrapperClusterBootstrap just contains Announce and the data necessary
// to identify and process the message in onet.
type announceWrapperClusterBootstrap struct {
	*onet.TreeNode
	ClusterBootstrapAnnounce
}

// ClusterBootstrapReply returns true when ready.
type ClusterBootstrapReply struct {
	Cluster *[]ClusterInstance
}

// replyWrapperClusterBootstrap just contains Reply and the data necessary to
// identify and process the message in onet.
type replyWrapperClusterBootstrap struct {
	*onet.TreeNode
	ClusterBootstrapReply
}

// StartARAProtocol structure
type StartARAProtocol struct {
	*onet.TreeNodeInstance
	announceChan chan announceWrapperStartARA
	repliesChan  chan []replyWrapperStartARA
	Ready        chan bool
	Info         ClusterInfo
	GetService   FnService
}

// StartARAAnnounce is used to pass a message to all children.
type StartARAAnnounce struct {
	SenderName string
	Bootstrap  string
	Secret     string
}

// announceWrapperStartARA just contains Announce and the data necessary
// to identify and process the message in onet.
type announceWrapperStartARA struct {
	*onet.TreeNode
	StartARAAnnounce
}

// StartARAReply returns true when ready.
type StartARAReply struct {
	Cluster *[]ClusterInstance
}

// replyWrapperStartARA just contains Reply and the data necessary to
// identify and process the message in onet.
type replyWrapperStartARA struct {
	*onet.TreeNode
	StartARAReply
}

// StartInstancesProtocol structure
type StartInstancesProtocol struct {
	*onet.TreeNodeInstance
	announceChan chan announceWrapperStartInstances
	repliesChan  chan []replyWrapperStartInstances
	Ready        chan bool
	GetService   FnService
	Nodes        map[string]*NodeInfo
}

// StartInstancesAnnounce is used to pass a message to all children.
type StartInstancesAnnounce struct{}

// announceWrapperStartInstances just contains Announce and the data necessary
// to identify and process the message in onet.
type announceWrapperStartInstances struct {
	*onet.TreeNode
	StartInstancesAnnounce
}

// StartInstancesReply returns true when ready.
type StartInstancesReply struct {
	Node *NodeInfo
}

// replyWrapperStartInstances just contains Reply and the data necessary to
// identify and process the message in onet.
type replyWrapperStartInstances struct {
	*onet.TreeNode
	StartInstancesReply
}
