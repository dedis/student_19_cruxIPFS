package service

/*
The service.go defines what to do for each API-call. This part of the service
runs on the node.
*/

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
	"sync"

	cruxIPFS "github.com/dedis/student_19_cruxIPFS"
	"github.com/dedis/student_19_cruxIPFS/ARAgen"
	"github.com/dedis/student_19_cruxIPFS/gentree"
	"go.dedis.ch/onet/v3"
	"go.dedis.ch/onet/v3/log"
	"go.dedis.ch/onet/v3/network"
	"go.dedis.ch/onet/v3/simul/monitor"
)

// Used for tests
var templateID onet.ServiceID

var execReqPingsMsgID network.MessageTypeID
var execReplyPingsMsgID network.MessageTypeID

var execReqIPFSInfoMsgID network.MessageTypeID
var execReplyIPFSInfoMsgID network.MessageTypeID

var execReqBootstrapClusterMsgID network.MessageTypeID
var execReplyBootstrapClusterMsgID network.MessageTypeID

func init() {
	var err error
	templateID, err = onet.RegisterNewService(cruxIPFS.ServiceName, newService)
	log.ErrFatal(err)

	execReqPingsMsgID = network.RegisterMessage(&cruxIPFS.ReqPings{})
	execReplyPingsMsgID = network.RegisterMessage(&cruxIPFS.ReplyPings{})

	execReqIPFSInfoMsgID = network.RegisterMessage(ReqIPFSInfo{})
	execReplyIPFSInfoMsgID = network.RegisterMessage(ReplyIPFSInfo{})

	execReqBootstrapClusterMsgID = network.RegisterMessage(ReqBootstrapCluster{})
	execReplyBootstrapClusterMsgID = network.RegisterMessage(ReplyBootstrapCluster{})

	network.RegisterMessage(&storage{})
}

// InitRequest init the tree
func (s *Service) InitRequest(req *cruxIPFS.InitRequest) (*cruxIPFS.InitResponse, error) {
	//log.Lvl1("here", s.ServerIdentity().String())

	s.Setup(req)

	return &cruxIPFS.InitResponse{}, nil
}

// Setup IPFS cluster ARAs
func (s *Service) Setup(req *cruxIPFS.InitRequest) {

	// copied from nyle
	log.Lvl3("Setup service")

	s.Nodes.All = req.Nodes
	s.Nodes.ServerIdentityToName = make(map[network.ServerIdentityID]string)
	for k, v := range req.ServerIdentityToName {
		s.Nodes.ServerIdentityToName[k.ID] = v
	}

	for _, myNode := range s.Nodes.All {

		myNode.ADist = make([]float64, 0)
		myNode.PDist = make([]string, 0)
		myNode.OptimalCluster = make(map[string]bool)
		myNode.OptimalBunch = make(map[string]bool)
		myNode.Cluster = make(map[string]bool)
		myNode.Bunch = make(map[string]bool)
		myNode.Rings = make([]string, 0)

	}
	// order nodesin s.Nodes in the order of index
	nodes := make([]*gentree.LocalityNode, len(s.Nodes.All))
	for _, n := range s.Nodes.All {
		nodes[gentree.NodeNameToInt(n.Name)] = n
	}
	s.Nodes.All = nodes
	s.Nodes.ClusterBunchDistances = make(map[*gentree.LocalityNode]map[*gentree.LocalityNode]float64)
	s.Nodes.Links = make(map[*gentree.LocalityNode]map[*gentree.LocalityNode]map[*gentree.LocalityNode]bool)
	s.GraphTree = make(map[string][]ARAgen.GraphTree)
	s.BinaryTree = make(map[string][]*onet.Tree)

	// allocate distances
	for _, node := range s.Nodes.All {
		s.Nodes.ClusterBunchDistances[node] = make(map[*gentree.LocalityNode]float64)
		s.Nodes.Links[node] = make(map[*gentree.LocalityNode]map[*gentree.LocalityNode]bool)
		for _, node2 := range s.Nodes.All {
			s.Nodes.ClusterBunchDistances[node][node2] = math.MaxFloat64
			s.Nodes.Links[node][node2] = make(map[*gentree.LocalityNode]bool)

			if node == node2 {
				s.Nodes.ClusterBunchDistances[node][node2] = 0
			}

			//log.LLvl1("init map", node.Name, node2.Name)
		}
	}

	//s.CosiWg = make(map[int]*sync.WaitGroup)
	//s.NodeWg = &sync.WaitGroup{}
	s.ClusterWg = &sync.WaitGroup{}
	s.PortMutex = &sync.Mutex{}
	s.metrics = make(map[string]*monitor.TimeMeasure)

	s.OwnPings = make(map[string]float64)
	s.PingDistances = make(map[string]map[string]float64)

	myip := strings.Split(s.ServerIdentity().String(), "/")
	myip = strings.Split(myip[len(myip)-1], ":")
	s.MyIP = myip[0]

	s.Name = s.Nodes.GetServerIdentityToName(s.ServerIdentity())

	if s.Name == "" {
		return
	}

	// wait after we created ARAs
	if LocalSim {
		s.NodeWg.Add(1)
	}

	//log.LLvl1("called init service on", s.Nodes.GetServerIdentityToName(s.ServerIdentity()), s.ServerIdentity())

	s.getPings(false)

	AuxNodes, dist2, ARATreeStruct, ARAOnetTrees := ARAgen.GenARAs(s.Nodes,
		s.Nodes.GetServerIdentityToName(s.ServerIdentity()), s.PingDistances, 3)

	s.Distances = dist2
	s.Nodes = AuxNodes
	s.GraphTree = ARATreeStruct
	s.BinaryTree = ARAOnetTrees

	if s.Name == Node0 {
		// print ping distances
		for _, n := range s.Nodes.All {
			str := n.Name + "\n"
			str += "Cluster: " + fmt.Sprintln(n.Cluster)
			str += "Bunch: " + fmt.Sprintln(n.Bunch)
			log.Lvl1(str)
		}
	}

	s.StartIPFS()
	s.ManageClusters()
	if s.Name == Node0 {
		log.Lvl1("All cluster instances successfully started")
	}
}

// NewProtocol is called on all nodes of a Tree (except the root, since it is
// the one starting the protocol) so it's the Service that will be called to
// generate the PI on all others node.
// If you use CreateProtocolOnet, this will not be called, as the Onet will
// instantiate the protocol on its own. If you need more control at the
// instantiation of the protocol, use CreateProtocolService, and you can
// give some extra-configuration to your protocol in here.
func (s *Service) NewProtocol(tn *onet.TreeNodeInstance, conf *onet.GenericConfig) (onet.ProtocolInstance, error) {
	log.Lvl3("Not templated yet")
	return nil, nil
}

// saves all data.
func (s *Service) save() {
	s.storage.Lock()
	defer s.storage.Unlock()
	err := s.Save(storageID, s.storage)
	if err != nil {
		log.Error("Couldn't save data:", err)
	}
}

// Tries to load the configuration and updates the data in the service
// if it finds a valid config-file.
func (s *Service) tryLoad() error {
	s.storage = &storage{}
	msg, err := s.Load(storageID)
	if err != nil {
		return err
	}
	if msg == nil {
		return nil
	}
	var ok bool
	s.storage, ok = msg.(*storage)
	if !ok {
		return errors.New("Data of wrong type")
	}
	return nil
}

// newService receives the context that holds information about the node it's
// running on. Saving and loading can be done using the context. The data will
// be stored in memory for tests and simulations, and on disk for real deployments.
func newService(c *onet.Context) (onet.Service, error) {
	//log.LLvl1("new service")

	s := &Service{
		ServiceProcessor: onet.NewServiceProcessor(c),
	}
	log.ErrFatal(s.RegisterHandlers(s.InitRequest))

	s.RegisterProcessorFunc(execReqPingsMsgID, s.ExecReqPings)
	s.RegisterProcessorFunc(execReplyPingsMsgID, s.ExecReplyPings)

	s.RegisterProcessorFunc(execReqIPFSInfoMsgID, s.ExecReqIPFSInfo)
	s.RegisterProcessorFunc(execReplyIPFSInfoMsgID, s.ExecReplyIPFSInfo)

	s.RegisterProcessorFunc(execReqBootstrapClusterMsgID, s.ExecReqBootstrapCluster)
	s.RegisterProcessorFunc(execReplyBootstrapClusterMsgID, s.ExecReplyBootstrapCluster)

	if err := s.tryLoad(); err != nil {
		log.Error(err)
		return nil, err
	}

	return s, nil
}

func (s *Service) getNodeID() int {
	n, err := strconv.Atoi(s.Nodes.GetServerIdentityToName(
		s.ServerIdentity())[len(NodeName):])
	checkErr(err)
	return n
}
