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

	NodeWg *sync.WaitGroup
	//CosiWg       map[int]*sync.WaitGroup
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
}

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
