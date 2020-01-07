package service

/*
The service.go defines what to do for each API-call. This part of the service
runs on the node.
*/

import (
	"errors"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
	"sync"

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

func init() {
	var err error
	templateID, err = onet.RegisterNewService(ServiceName, newService)
	log.ErrFatal(err)

	execReqPingsMsgID = network.RegisterMessage(&ReqPings{})
	execReplyPingsMsgID = network.RegisterMessage(&ReplyPings{})

	for _, i := range []interface{}{
		StartIPFSAnnounce{},
		StartIPFSReply{},
		ClusterBootstrapAnnounce{},
		ClusterBootstrapReply{},
		StartARAAnnounce{},
		StartARAReply{},
		StartInstancesAnnounce{},
		StartInstancesReply{},
		&storage{},
	} {
		network.RegisterMessage(i)
	}

}

// InitRequest init the tree
func (s *Service) InitRequest(req *InitRequest) (
	*InitResponse, error) {

	s.setup(req)

	return &InitResponse{}, nil
}

// Setup IPFS cluster ARAs
func (s *Service) setup(req *InitRequest) {

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
	s.Nodes.ClusterBunchDistances =
		make(map[*gentree.LocalityNode]map[*gentree.LocalityNode]float64)
	s.Nodes.Links = make(map[*gentree.
		LocalityNode]map[*gentree.LocalityNode]map[*gentree.LocalityNode]bool)
	s.GraphTree = make(map[string][]gentree.GraphTree)
	s.BinaryTree = make(map[string][]*onet.Tree)

	// allocate distances
	for _, node := range s.Nodes.All {
		s.Nodes.ClusterBunchDistances[node] =
			make(map[*gentree.LocalityNode]float64)
		s.Nodes.Links[node] =
			make(map[*gentree.LocalityNode]map[*gentree.LocalityNode]bool)
		for _, node2 := range s.Nodes.All {
			s.Nodes.ClusterBunchDistances[node][node2] = math.MaxFloat64
			s.Nodes.Links[node][node2] = make(map[*gentree.LocalityNode]bool)

			if node == node2 {
				s.Nodes.ClusterBunchDistances[node][node2] = 0
			}
		}
	}

	s.PortMutex = &sync.Mutex{}
	s.metrics = make(map[string]*monitor.TimeMeasure)

	s.OwnPings = make(map[string]float64)
	s.PingDistances = make(map[string]map[string]float64)

	s.OnetTree = req.OnetTree

	myip := strings.Split(s.ServerIdentity().String(), "/")
	myip = strings.Split(myip[len(myip)-1], ":")
	s.MyIP = myip[0]

	s.Name = s.Nodes.GetServerIdentityToName(s.ServerIdentity())

	_, err := os.Stat(PingsFile)
	s.getPings(err == nil && !req.ComputePings)
	//os.IsNotExist(err))

	//s.getPings(false)
	if s.Nodes.GetServerIdentityToName(s.ServerIdentity()) == Node0 {
		//s.printDistances("Ping distances")
		s.printPings()
	}

	s.setIPFSVariables()

	if !req.Cruxified {
		return
	}

	maxLvl := 0
	for _, n := range s.Nodes.All {
		if n.Level > maxLvl {
			maxLvl = n.Level
		}
	}
	maxLvl++

	AuxNodes, dist2, ARATreeStruct, ARAOnetTrees := gentree.GenARAs(s.Nodes,
		s.Nodes.GetServerIdentityToName(s.ServerIdentity()),
		s.PingDistances, maxLvl)

	s.Distances = dist2
	s.Nodes = AuxNodes
	s.GraphTree = ARATreeStruct
	s.BinaryTree = ARAOnetTrees
}

// NewProtocol is called on all nodes of a Tree (except the root, since it is
// the one starting the protocol) so it's the Service that will be called to
// generate the PI on all others node.
// If you use CreateProtocolOnet, this will not be called, as the Onet will
// instantiate the protocol on its own. If you need more control at the
// instantiation of the protocol, use CreateProtocolService, and you can
// give some extra-configuration to your protocol in here.
func (s *Service) NewProtocol(tn *onet.TreeNodeInstance,
	conf *onet.GenericConfig) (onet.ProtocolInstance, error) {
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

	s := &Service{
		ServiceProcessor: onet.NewServiceProcessor(c),
	}
	log.ErrFatal(s.RegisterHandlers(s.InitRequest))

	s.RegisterProcessorFunc(execReqPingsMsgID, s.ExecReqPings)
	s.RegisterProcessorFunc(execReplyPingsMsgID, s.ExecReplyPings)

	_, err := s.ProtocolRegister(StartIPFSName, func(n *onet.TreeNodeInstance) (
		onet.ProtocolInstance, error) {

		return NewStartIPFSProtocol(n, s.GetService)
	})
	if err != nil {
		log.Error(err)
		return nil, err
	}

	_, err = s.ProtocolRegister(ClusterBootstrapName,
		func(n *onet.TreeNodeInstance) (onet.ProtocolInstance, error) {

			return NewClusterBootstrapProtocol(n, s.GetService)
		})
	if err != nil {
		log.Error(err)
		return nil, err
	}

	_, err = s.ProtocolRegister(StartARAName,
		func(n *onet.TreeNodeInstance) (onet.ProtocolInstance, error) {

			return NewStartARAProtocol(n, s.GetService)
		})
	if err != nil {
		log.Error(err)
		return nil, err
	}

	_, err = s.ProtocolRegister(StartInstancesName,
		func(n *onet.TreeNodeInstance) (onet.ProtocolInstance, error) {

			return NewStartInstancesProtocol(n, s.GetService)
		})
	if err != nil {
		log.Error(err)
		return nil, err
	}

	if err = s.tryLoad(); err != nil {
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

// GetService Returns the Current SERVICE
func (s *Service) GetService() *Service {
	return s
}

// PrintName PrintName
func (s *Service) PrintName() {
	fmt.Println(s.Name)
}
