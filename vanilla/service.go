package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/BurntSushi/toml"

	cruxIPFS "github.com/dedis/student_19_cruxIPFS"
	"github.com/dedis/student_19_cruxIPFS/gentree"
	"github.com/dedis/student_19_cruxIPFS/operations"
	"github.com/dedis/student_19_cruxIPFS/service"

	"go.dedis.ch/onet/v3"
	"go.dedis.ch/onet/v3/app"
	"go.dedis.ch/onet/v3/log"
	"go.dedis.ch/onet/v3/network"
)

/*
 * Defines the simulation for the service-template
 */

func init() {
	onet.SimulationRegister(cruxIPFS.ServiceName, NewSimulationService)
}

// NewSimulationService returns the new simulation, where all fields are
// initialised using the config-file
func NewSimulationService(config string) (onet.Simulation, error) {
	es := &IPFSSimulation{}
	_, err := toml.Decode(config, es)
	if err != nil {
		return nil, err
	}
	return es, nil
}

// Setup creates the tree used for that simulation
func (s *IPFSSimulation) Setup(dir string, hosts []string) (
	*onet.SimulationConfig, error) {
	log.Lvl1("Starting Setup()")

	SetNodePaths(len(hosts))

	app.Copy(dir, filepath.Join(DATAFOLDER, NODEPATHREMOTE))
	app.Copy(dir, "prescript.sh")
	app.Copy(dir, "local_nodes.txt")
	app.Copy(dir, "install/ipfs")
	app.Copy(dir, "install/ipfs-cluster-service")

	sc := &onet.SimulationConfig{}
	s.CreateRoster(sc, hosts, 2000)
	err := s.CreateTree(sc)
	if err != nil {
		return nil, err
	}
	return sc, nil
}

// Node can be used to initialize each node before it will be run
// by the server. Here we call the 'Node'-method of the
// SimulationBFTree structure which will load the roster- and the
// tree-structure to speed up the first round.
func (s *IPFSSimulation) Node(config *onet.SimulationConfig) error {
	log.Lvl1("Starting Node()")

	index, _ := config.Roster.Search(config.Server.ServerIdentity.ID)
	if index < 0 {
		log.Fatal("Didn't find this node in roster")
	}
	log.Lvl3("Initializing node-index", index)

	//config.Overlay.RegisterTree()

	s.ReadNodeInfo(false)

	mymap := s.initializeMaps(config, true)

	myService := config.GetService(cruxIPFS.ServiceName).(*service.Service)

	serviceReq := &cruxIPFS.InitRequest{
		Nodes:                s.Nodes.All,
		ServerIdentityToName: mymap,
		OnetTree:             config.Tree,
		Roster:               config.Roster,
		Cruxified:            false,
	}

	myService.InitRequest(serviceReq)

	bt := make(map[string][]*onet.Tree)
	bt[service.Node0] = []*onet.Tree{config.Tree}
	myService.BinaryTree = bt

	return s.SimulationBFTree.Node(config)
}

// Run is used on the destination machines and runs a number of
// rounds
func (s *IPFSSimulation) Run(config *onet.SimulationConfig) error {
	log.Lvl1("Starting Run()")

	myService := config.GetService(cruxIPFS.ServiceName).(*service.Service)

	pi, err := myService.CreateProtocol(service.StartIPFSName, config.Tree)
	if err != nil {
		fmt.Println(err)
	}
	pi.Start()

	<-pi.(*service.StartIPFSProtocol).Ready
	operations.SaveState(cruxIPFS.SaveFile, pi.(*service.StartIPFSProtocol).Nodes)

	time.Sleep(20 * time.Second)

	operations.Test0()
	return nil
}

// ReadNodeInfo read node information
func (s *IPFSSimulation) ReadNodeInfo(isLocalTest bool) {
	_, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	if isLocalTest {
		//log.Lvl1("NODEPATHLOCAL:", NODEPATHLOCAL)
		s.ReadNodesFromFile(NODEPATHLOCAL)
	} else {
		//log.Lvl1("NODEPATHREMOTE:", "nodes_local_11.txt")
		s.ReadNodesFromFile("nodes_local_11.txt")
	}
}

// ReadNodesFromFile read nodes information from a text file
func (s *IPFSSimulation) ReadNodesFromFile(filename string) {
	s.Nodes.All = make([]*gentree.LocalityNode, 0)

	readLine := cruxIPFS.ReadFileLineByLine(filename)

	for true {
		line := readLine()
		//fmt.Println(line)
		if line == "" {
			//fmt.Println("end")
			break
		}

		if strings.HasPrefix(line, "#") {
			continue
		}

		tokens := strings.Split(line, " ")
		coords := strings.Split(tokens[1], ",")
		name, xstr, ystr, IP, levelstr := tokens[0], coords[0], coords[1], tokens[2], tokens[3]

		x, _ := strconv.ParseFloat(xstr, 64)
		y, _ := strconv.ParseFloat(ystr, 64)
		level, err := strconv.Atoi(levelstr)

		if err != nil {
			log.Lvl1("Error", err)

		}

		//	log.Lvl1("reqd node level", name, level_str, "lvl", level)

		myNode := cruxIPFS.CreateNode(name, x, y, IP, level)
		s.Nodes.All = append(s.Nodes.All, myNode)
	}
}

func (s *IPFSSimulation) initializeMaps(config *onet.SimulationConfig, isLocalTest bool) map[*network.ServerIdentity]string {

	s.Nodes.ServerIdentityToName = make(map[network.ServerIdentityID]string)
	ServerIdentityToName := make(map[*network.ServerIdentity]string)

	if isLocalTest {
		for i := range s.Nodes.All {
			treeNode := config.Tree.List()[i]
			// quick fix
			//treeNode := config.Tree.List()[i+1]
			s.Nodes.All[i].ServerIdentity = treeNode.ServerIdentity
			s.Nodes.ServerIdentityToName[treeNode.ServerIdentity.ID] = s.Nodes.All[i].Name
			ServerIdentityToName[treeNode.ServerIdentity] = s.Nodes.All[i].Name
			//log.Lvl1("associating", treeNode.ServerIdentity.String(), "to", s.Nodes.All[i].Name)
		}
	} else {
		for _, treeNode := range config.Tree.List() {
			serverIP := treeNode.ServerIdentity.Address.Host()
			node := s.Nodes.GetByIP(serverIP)
			node.ServerIdentity = treeNode.ServerIdentity
			s.Nodes.ServerIdentityToName[treeNode.ServerIdentity.ID] = node.Name
			ServerIdentityToName[treeNode.ServerIdentity] = node.Name
			//log.Lvl1("associating", treeNode.ServerIdentity.String(), "to", node.Name)
		}
	}

	return ServerIdentityToName
}
