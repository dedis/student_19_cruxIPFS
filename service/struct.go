package service

import (
	"bufio"
	"os"
	"sync"

	"github.com/dedis/student_19_cruxIPFS/ARAgen"
	"github.com/dedis/student_19_cruxIPFS/gentree"

	"go.dedis.ch/onet/v3"
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
	GraphTree    map[string][]ARAgen.GraphTree
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
	MyIPFS     IPFSInformation            // own ipfs information
	OtherIPFS  map[string]IPFSInformation // node_x -> IP, ports etc.

	OnetTree      *onet.Tree
	StartIPFSProt onet.ProtocolInstance
}

type ServiceFn func() *Service

// storage is used to save our data.
type storage struct {
	Count int
	sync.Mutex
}

// ClusterInstance details of a cluster
type ClusterInstance struct {
	IP            string
	IPFSAPIPort   int
	RestAPIPort   int
	IPFSProxyPort int
	ClusterPort   int
}

// IPFSInformation structure containing information about an IPFS instance
type IPFSInformation struct {
	IP          string
	SwarmPort   int
	APIPort     int
	GatewayPort int
}

// ReqIPFSInfo IPFS information request packet
type ReqIPFSInfo struct {
	SenderName string
}

// ReplyIPFSInfo IPFS information reply packet
type ReplyIPFSInfo struct {
	SenderName string
	IPFSInfo   IPFSInformation
}

type ReqBootstrapCluster struct {
	SenderName string
	Bootstrap  string
	Secret     string
}

type ReplyBootstrapCluster struct {
	SenderName string
}

// Name can be used from other packages to refer to this protocol.
const Name = "Template"

// Announce is used to pass a message to all children.
type Announce struct {
	Message string
}

// announceWrapper just contains Announce and the data necessary to identify
// and process the message in onet.
type announceWrapper struct {
	*onet.TreeNode
	Announce
}

// Reply returns the count of all children.
type Reply struct {
	ChildrenCount int
}

// replyWrapper just contains Reply and the data necessary to identify and
// process the message in onet.
type replyWrapper struct {
	*onet.TreeNode
	Reply
}

// WaitpeersProtocol structure
type WaitpeersProtocol struct {
	*onet.TreeNodeInstance
	announceChan chan announceWrapperWaitpeers
	repliesChan  chan []replyWrapperWaitpeers
	Ready        chan bool
}

// WaitpeersAnnounce is used to pass a message to all children.
type WaitpeersAnnounce struct {
	Message string
}

// announceWrapperWaitpeers just contains Announce and the data necessary to
// identify and process the message in onet.
type announceWrapperWaitpeers struct {
	*onet.TreeNode
	WaitpeersAnnounce
}

// WaitpeersReply returns true when ready.
type WaitpeersReply struct {
	Ready bool
}

// replyWrapper just contains Reply and the data necessary to identify and
// process the message in onet.
type replyWrapperWaitpeers struct {
	*onet.TreeNode
	WaitpeersReply
}

// StartIPFSProtocol structure
type StartIPFSProtocol struct {
	*onet.TreeNodeInstance
	announceChan chan announceWrapperStartIPFS
	repliesChan  chan []replyWrapperStartIPFS
	Ready        chan bool
	GetService   ServiceFn
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
	Ready bool
}

// replyWrapper just contains Reply and the data necessary to identify and
// process the message in onet.
type replyWrapperStartIPFS struct {
	*onet.TreeNode
	StartIPFSReply
}

// LaunchClustersProtocol structure
type LaunchClustersProtocol struct {
	*onet.TreeNodeInstance
	announceChan chan announceWrapperLaunchClusters
	repliesChan  chan []replyWrapperLaunchClusters
	Ready        chan bool
	GetService   ServiceFn
}

// LaunchClustersAnnounce is used to pass a message to all children.
type LaunchClustersAnnounce struct {
	Message string
}

// announceWrapperLaunchClusters just contains Announce and the data necessary
// to identify and process the message in onet.
type announceWrapperLaunchClusters struct {
	*onet.TreeNode
	LaunchClustersAnnounce
}

// LaunchClustersReply returns true when ready.
type LaunchClustersReply struct {
	Ready bool
}

// replyWrapperLaunchClusters just contains Reply and the data necessary to
// identify and process the message in onet.
type replyWrapperLaunchClusters struct {
	*onet.TreeNode
	LaunchClustersReply
}

// ClusterBootstrapProtocol structure
type ClusterBootstrapProtocol struct {
	*onet.TreeNodeInstance
	announceChan chan announceWrapperClusterBootstrap
	repliesChan  chan []replyWrapperClusterBootstrap
	Ready        chan bool
	GetService   ServiceFn
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
	Ready bool
}

// replyWrapperClusterBootstrap just contains Reply and the data necessary to
// identify and process the message in onet.
type replyWrapperClusterBootstrap struct {
	*onet.TreeNode
	ClusterBootstrapReply
}
