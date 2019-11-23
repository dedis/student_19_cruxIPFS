package service

import (
	"bufio"
	"os"
	"sync"

	"github.com/dedis/student_19_cruxIPFS/ARAgen"
	"github.com/dedis/student_19_cruxIPFS/gentree"
	"go.dedis.ch/onet/v3"
	"go.dedis.ch/onet/v3/network"
	"go.dedis.ch/onet/v3/simul/monitor"
)

// storageID reflects the data we're storing - we could store more
// than one structure.
var storageID = []byte("main")

var execReqPingsMsgID network.MessageTypeID
var execReplyPingsMsgID network.MessageTypeID

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

	CosiWg       map[int]*sync.WaitGroup
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
}

// storage is used to save our data.
type storage struct {
	Count int
	sync.Mutex
}
