package service

/*
The service.go defines what to do for each API-call. This part of the service
runs on the node.
*/

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"errors"

	"github.com/dedis/paper_crux/dsn_exp/gentree"
	template "github.com/dedis/student_19_cruxIPFS"
	"github.com/dedis/student_19_cruxIPFS/protocol"
	"go.dedis.ch/onet"
	"go.dedis.ch/onet/log"
	"go.dedis.ch/onet/network"
)

// Used for tests
var templateID onet.ServiceID

// storageID reflects the data we're storing - we could store more
// than one structure.
var storageID = []byte("main")

// Name of the service
var Name = "IPFS"

func init() {
	var err error
	templateID, err = onet.RegisterNewService(Name, newService)
	log.ErrFatal(err)
	network.RegisterMessage(&storage{})
}

// InitRequest handles initialisation requests
func (s *Service) InitRequest(req *InitRequest) (*InitResponse, error) {
	//log.Lvl1("here", s.ServerIdentity().String())
	s.Setup(req)

	return &InitResponse{}, nil
}

// Setup the service
func (s *Service) Setup(req *InitRequest) {
	// set the nodes
	s.Nodes.All = req.Nodes

	// fill in the server identity to name
	s.Nodes.ServerIdentityToName = make(map[network.ServerIdentityID]string)
	for k, v := range req.ServerIdentityToName {
		s.Nodes.ServerIdentityToName[k.ID] = v
	}

	nodes := make([]*gentree.LocalityNode, len(s.Nodes.All))
	for _, n := range s.Nodes.All {
		nodes[gentree.NodeNameToInt(n.Name)] = n
		//log.Info(s.ServerIdentity(), fmt.Sprintf("%+v", nodes[gentree.NodeNameToInt(n.Name)]))
	}
	s.Nodes.All = nodes

	s.OwnPings = make(map[string]float64)
	s.PingDistances = make(map[string]map[string]float64)
	s.NrPingAnswers = 0

	// compute the ping from hosts
	s.ComputePing()
	fmt.Println(s.OwnPings)
}

// ComputePing computes the ping distance used as metric between nodes
func (s *Service) ComputePing() {
	log.Lvl1("Computing ping", len(s.Nodes.All))
	myName := s.Nodes.GetServerIdentityToName(s.ServerIdentity())
	for _, node := range s.Nodes.All {

		s1 := node.ServerIdentity

		fmt.Println(s1)

		//if node.ServerIdentity.String() != s.ServerIdentity().String() {
		if node.ServerIdentity.Address.String() != s.ServerIdentity().String() {

			log.Lvl2(myName, "meas ping to ", s.Nodes.GetServerIdentityToName(node.ServerIdentity))

			for {

				peerName := node.Name
				pingCmdStr := "ping -W 150 -q -c 3 -i 1 " + node.ServerIdentity.Address.Host() + " | tail -n 1"
				pingCmd := exec.Command("sh", "-c", pingCmdStr)
				pingOutput, err := pingCmd.Output()
				if err != nil {
					log.Fatal("couldn't measure ping")
				}

				if strings.Contains(string(pingOutput), "pipe") || len(strings.TrimSpace(string(pingOutput))) == 0 {
					log.Lvl1(s.Nodes.GetServerIdentityToName(s.ServerIdentity()), "retrying for", peerName, node.ServerIdentity.Address.Host(), node.ServerIdentity.String())
					log.LLvl1("retry")
					continue
				}

				processPingCmdStr := "echo " + string(pingOutput) + " | cut -d ' ' -f 4 | cut -d '/' -f 1-2 | tr '/' ' '"
				processPingCmd := exec.Command("sh", "-c", processPingCmdStr)
				processPingOutput, _ := processPingCmd.Output()
				if string(pingOutput) == "" {
					//log.Lvl1("empty ping ", myName + " " + peerName)
				} else {
					//log.Lvl1("%%%%%%%%%%%%% ping ", myName + " " + peerName, "output ", string(pingOutput), "processed output ", string(processPingOutput))
				}

				//log.Lvl1("%%%%%%%%%%%%% ping ", s.Nodes.GetServerIdentityToName(s.ServerIdentity()) + " " + peerName, "output ", string(pingOutput), "processed output ", string(processPingOutput))

				strPingOut := string(processPingOutput)

				pingRes := strings.Split(strPingOut, "/")
				//log.LLvl1("pingRes", pingRes)

				avgPing, err := strconv.ParseFloat(pingRes[5], 64)
				if err != nil {
					log.Fatal("Problem when parsing pings")
				}

				s.OwnPings[peerName] = float64(avgPing / 2.0)

				break
			}

		}
	}
}

// Clock starts a template-protocol and returns the run-time.
func (s *Service) Clock(req *template.Clock) (*template.ClockReply, error) {
	s.storage.Lock()
	s.storage.Count++
	s.storage.Unlock()
	s.save()
	tree := req.Roster.GenerateNaryTreeWithRoot(2, s.ServerIdentity())
	if tree == nil {
		return nil, errors.New("couldn't create tree")
	}
	pi, err := s.CreateProtocol(protocol.Name, tree)
	if err != nil {
		return nil, err
	}
	start := time.Now()
	pi.Start()
	resp := &template.ClockReply{
		Children: <-pi.(*protocol.TemplateProtocol).ChildCount,
	}
	resp.Time = time.Now().Sub(start).Seconds()
	return resp, nil
}

// Count returns the number of instantiations of the protocol.
func (s *Service) Count(req *template.Count) (*template.CountReply, error) {
	s.storage.Lock()
	defer s.storage.Unlock()
	return &template.CountReply{Count: s.storage.Count}, nil
}

// GenSecret generate a secret key for ipfs cluster
func (s *Service) GenSecret(req *template.GenSecret) (*template.GenSecretReply,
	error) {

	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		return nil, errors.New("could not generate secret")
	}

	reply := template.GenSecretReply{Secret: hex.EncodeToString(key)}
	return &reply, nil
}

// StartIPFS launch an ipfs instance, and returns the used ports
// it cleans the config folder, init ipfs in there, get unused ports, edit the
// config with the available ports and given id, starts ipfs daemon and returns
// the ports used by ipfs in this order [swarm, api, gateway]
func (s *Service) StartIPFS(req *template.StartIPFS) (*template.StartIPFSReply,
	error) {

	path := req.ConfigPath + "/ipfs"

	// create the empty directory that will store ipfs configs
	err := CreateEmptyDir(path)
	if err != nil {
		return nil, err
	}

	// init ipfs in the desired folder
	exec.Command("ipfs", "-c"+path, "init").Run()

	// edit the ip in the config file
	EditIPFSConfig(path, req.IP)

	// start the ipfs daemon
	// we need to fork the process
	go exec.Command("ipfs", "-c"+path, "daemon").Run()

	// sleep with the daemon launches
	time.Sleep(13 * time.Second)

	return &template.StartIPFSReply{}, nil
}

// StartCluster start a cluster instance
func (s *Service) StartCluster(req *template.StartCluster) (
	*template.StartClusterReply, error) {

	peername := "clusterX"
	replmin := 3
	replmax := 5

	err := Protocol(req.ConfigPath, peername, req.IP, replmin, replmax)
	if err != nil {
		return nil, err
	}

	return &template.StartClusterReply{}, nil
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
	s := &Service{
		ServiceProcessor: onet.NewServiceProcessor(c),
	}
	if err := s.RegisterHandlers(s.Clock, s.Count,
		s.GenSecret, s.StartIPFS, s.StartCluster); err != nil {
		return nil, errors.New("Couldn't register messages")
	}
	if err := s.tryLoad(); err != nil {
		log.Error(err)
		return nil, err
	}
	return s, nil
}
