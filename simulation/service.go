package main

import (
	"fmt"
	"io/ioutil"
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

const cruxified = true

func init() {
	onet.SimulationRegister(cruxIPFS.ServiceName, NewIPFSSimulation)
}

// NewIPFSSimulation returns the new simulation, where all fields are
// initialised using the config-file
func NewIPFSSimulation(config string) (onet.Simulation, error) {
	es := &IPFSSimulation{}
	_, err := toml.Decode(config, es)
	if err != nil {
		return nil, err
	}
	return es, nil
}

func (s *IPFSSimulation) Setup(dir string, hosts []string) (
	*onet.SimulationConfig, error) {

	app.Copy(dir, filepath.Join(DATAFOLDER, NODEPATHREMOTE))
	app.Copy(dir, "prescript.sh")
	app.Copy(dir, "nodes.txt")
	app.Copy(dir, "install/ipfs")
	app.Copy(dir, "install/ipfs-cluster-service")

	b, err := ioutil.ReadFile("../detergen/details.txt")
	if err != nil {
		log.Error(err)
	}
	log.Lvl1("\ndetails: simulation, mode: " +
		service.ClusterConsensusMode + ", " + string(b))

	sc := &onet.SimulationConfig{}
	s.CreateRoster(sc, hosts, 2000)
	err = s.CreateTree(sc)
	if err != nil {
		return nil, err
	}
	return sc, nil
}

func (s *IPFSSimulation) Node(config *onet.SimulationConfig) error {

	index, _ := config.Roster.Search(config.Server.ServerIdentity.ID)
	if index < 0 {
		log.Fatal("Didn't find this node in roster")
	}
	log.Lvl3("Initializing node-index", index)

	s.ReadNodeInfo(false, *config)

	mymap := s.initializeMaps(config, true)

	myService := config.GetService(cruxIPFS.ServiceName).(*service.Service)

	serviceReq := &cruxIPFS.InitRequest{
		Nodes:                s.Nodes.All,
		ServerIdentityToName: mymap,
		OnetTree:             config.Tree,
		Roster:               config.Roster,
		Cruxified:            cruxified,
	}

	myService.InitRequest(serviceReq)

	if cruxified {
		for _, trees := range myService.BinaryTree {
			for _, tree := range trees {
				config.Overlay.RegisterTree(tree)
			}
		}
	} else {
		bt := make(map[string][]*onet.Tree)
		bt[service.Node0] = []*onet.Tree{config.Tree}
		myService.BinaryTree = bt
	}

	return s.SimulationBFTree.Node(config)
}

func (s *IPFSSimulation) Run(config *onet.SimulationConfig) error {
	myService := config.GetService(cruxIPFS.ServiceName).(*service.Service)

	pi, err := myService.CreateProtocol(service.StartIPFSName, config.Tree)
	if err != nil {
		fmt.Println(err)
	}
	pi.Start()

	<-pi.(*service.StartIPFSProtocol).Ready
	operations.SaveState(cruxIPFS.SaveFile, pi.(*service.StartIPFSProtocol).Nodes)

	// wait for some time for clusters to converge
	time.Sleep(20 * time.Second)
	operations.Test2(100, len(myService.Nodes.All))
	return nil
}

// ReadNodeInfo read node information
func (s *IPFSSimulation) ReadNodeInfo(isLocalTest bool, config onet.SimulationConfig) {
	_, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	s.ReadNodesFromFile("nodes.txt", config)
	/*
		if isLocalTest {
			//log.Lvl1("NODEPATHLOCAL:", NODEPATHLOCAL)
			s.ReadNodesFromFile(NODEPATHLOCAL, config)
		} else {
			//log.Lvl1("NODEPATHREMOTE:", "nodes_local_11.txt")
			s.ReadNodesFromFile("nodes_local_11.txt", config)
		}
	*/
}

// ReadNodesFromFile read nodes information from a text file
func (s *IPFSSimulation) ReadNodesFromFile(filename string, config onet.SimulationConfig) {
	s.Nodes.All = make([]*gentree.LocalityNode, 0)

	readLine := cruxIPFS.ReadFileLineByLine(filename)

	for i := 0; i < len(config.Roster.List); i++ {
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
		name, _, _, _, levelstr := tokens[0], coords[0], coords[1], tokens[2], tokens[3]

		level, err := strconv.Atoi(levelstr)

		if err != nil {
			log.Lvl1("Error", err)

		}

		myNode := gentree.CreateNode(name, level)
		s.Nodes.All = append(s.Nodes.All, myNode)
	}
}

func (s *IPFSSimulation) initializeMaps(config *onet.SimulationConfig, isLocalTest bool) map[*network.ServerIdentity]string {

	s.Nodes.ServerIdentityToName = make(map[network.ServerIdentityID]string)
	ServerIdentityToName := make(map[*network.ServerIdentity]string)

	if isLocalTest {
		for i := range s.Nodes.All {
			treeNode := config.Tree.List()[i]
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
		}
	}

	return ServerIdentityToName
}
