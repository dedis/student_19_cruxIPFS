package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/dedis/paper_crux/dsn_exp/gentree"
	template "github.com/dedis/student_19_cruxIPFS"
	"github.com/dedis/student_19_cruxIPFS/service"
	"go.dedis.ch/onet"
	"go.dedis.ch/onet/app"
	"go.dedis.ch/onet/log"
	"go.dedis.ch/onet/network"
)

const simName = "IPFS"

// NrOps ???
const NrOps = 90000

var mySC onet.SimulationConfig

/*
 * Defines the simulation for the service-template
 */

func init() {
	onet.SimulationRegister(simName, NewIPFSSimulation)
}

// IPFSSimulation only holds the BFTree simulation
type IPFSSimulation struct {
	onet.SimulationBFTree

	Nodes gentree.LocalityNodes
}

// NewIPFSSimulation returns the new ipfs simulation, where all fields are
// initialised using the config-file
func NewIPFSSimulation(config string) (onet.Simulation, error) {
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
	app.Copy(dir, "../clean.sh")
	nodePath := filepath.Join(DataFolder, NodeFile)
	app.Copy(dir, nodePath)

	sc := &onet.SimulationConfig{}

	s.CreateRoster(sc, hosts, 2000)

	err := s.CreateTree(sc)
	if err != nil {
		return nil, err
	}
	mySC = *sc

	return sc, nil
}

// Node can be used to initialize each node before it will be run
// by the server. Here we call the 'Node'-method of the
// SimulationBFTree structure which will load the roster- and the
// tree-structure to speed up the first round.
func (s *IPFSSimulation) Node(config *onet.SimulationConfig) error {
	index, _ := config.Roster.Search(config.Server.ServerIdentity.ID)
	if index < 0 {
		log.Fatal("Didn't find this node in roster")
	}
	log.Lvl3("Initializing node-index", index)

	dir, err := os.Getwd()
	if err != nil {
		fmt.Print(err)
		panic(err)
	}

	nodesLocation := filepath.Join(dir, NodeFile)
	fmt.Println("Node location:", nodesLocation)
	s.ReadNodes(nodesLocation)

	myService := config.GetService(service.Name).(*service.Service)

	mymap := s.InitializeMaps(config, true)

	serviceReq := &service.InitRequest{
		Nodes:                s.Nodes.All,
		ServerIdentityToName: mymap,
		NrOps:                NrOps,
		OpIdxStart:           2 * s.Hosts,
		Roster:               config.Roster,
	}
	_, err = myService.InitRequest(serviceReq)
	if err != nil {
		return err
	}

	//StartIPFSDaemon(config, index)

	return s.SimulationBFTree.Node(config)
}

// Run is used on the destination machines and runs a number of
// rounds
func (s *IPFSSimulation) Run(config *onet.SimulationConfig) error {
	size := config.Tree.Size()
	log.Lvl2("Size is:", size, "rounds:", s.Rounds)
	//c := template.NewClient()
	for round := 0; round < s.Rounds; round++ {
		log.Lvl1("Starting round", round)
		/*
			round := monitor.NewTimeMeasure("round")
			resp, err := c.Clock(config.Roster)
			log.ErrFatal(err)

			if resp.Time <= 0 {
				log.Fatal("0 time elapsed")
			}
			round.Record()
		*/
	}
	return nil
}

// StartIPFSDaemon select ports for each peer and start ipfs on those ports
func StartIPFSDaemon(sc *onet.SimulationConfig, index int) {

	c := template.NewClient()
	// mySC is the SimulConfig created in Setup()
	// I guess I shouldn't use it, but it runs
	identity := mySC.Roster.Get(index)

	// ip of the node that will start
	ip := template.ServerIdentityToIPString(identity)
	// path of the config files of that node
	configPath := template.ConfigPath + "/Node" + strconv.Itoa(index)

	// create the ipfs start request
	req := template.StartIPFS{
		ConfigPath: configPath,
		IP:         ip,
	}
	reply := &template.StartIPFSReply{}
	err := c.SendProtobuf(identity, &req, reply)
	if err != nil {
		fmt.Println(err)
	}

	req2 := template.StartCluster{
		ConfigPath: configPath,
		IP:         ip,
	}

	c2 := template.NewClient()
	reply2 := &template.StartClusterReply{}
	err = c2.SendProtobuf(identity, &req2, reply2)
	if err != nil {
		fmt.Println(err)
	}
}

// InitializeMaps ??
func (s *IPFSSimulation) InitializeMaps(config *onet.SimulationConfig, isLocalTest bool) map[*network.ServerIdentity]string {

	s.Nodes.ServerIdentityToName = make(map[network.ServerIdentityID]string)
	ServerIdentityToName := make(map[*network.ServerIdentity]string)

	nextPortsAvailable := make(map[string]int)
	portIncrement := 1000

	// get machines

	for _, node := range config.Tree.List() {
		machineAddr := strings.Split(strings.Split(node.ServerIdentity.Address.String(), "//")[1], ":")[0]
		//log.LLvl1("machineaddr", machineAddr)
		log.Lvl2("node addr", node.ServerIdentity.Address.String())
		nextPortsAvailable[machineAddr] = 20000
	}

	if isLocalTest {

		for _, treeNode := range config.Tree.List() {
			for i := range s.Nodes.All {

				machineAddr := strings.Split(strings.Split(treeNode.ServerIdentity.Address.String(), "//")[1], ":")[0]
				if !s.Nodes.All[i].IP[machineAddr] {
					continue
				}

				if s.Nodes.All[i].ServerIdentity != nil {
					// current node already has stuff assigned to it, get the next free one
					continue
				}

				if treeNode.ServerIdentity != nil && treeNode.ServerIdentity.Address == "" {
					log.Error("nil 132132", s.Nodes.All[i].Name)
				}

				s.Nodes.All[i].ServerIdentity = treeNode.ServerIdentity
				s.Nodes.All[i].ServerIdentity.Address = treeNode.ServerIdentity.Address

				// set reserved ports
				s.Nodes.All[i].AvailablePortsStart = nextPortsAvailable[machineAddr]
				s.Nodes.All[i].AvailablePortsEnd = s.Nodes.All[i].AvailablePortsStart + portIncrement
				// fot all IP addresses of the machine set the ports!

				for k, v := range s.Nodes.All[i].IP {
					if v {
						nextPortsAvailable[k] = s.Nodes.All[i].AvailablePortsEnd
					}
				}

				s.Nodes.All[i].NextPort = s.Nodes.All[i].AvailablePortsStart
				// set names
				s.Nodes.ServerIdentityToName[treeNode.ServerIdentity.ID] = s.Nodes.All[i].Name
				ServerIdentityToName[treeNode.ServerIdentity] = s.Nodes.All[i].Name

				log.Lvl1("associating", treeNode.ServerIdentity.String(), "to", s.Nodes.All[i].Name, "ports", s.Nodes.All[i].AvailablePortsStart, s.Nodes.All[i].AvailablePortsEnd, s.Nodes.All[i].ServerIdentity.Address)

				break
			}

		}
	} else {
		for _, treeNode := range config.Tree.List() {
			serverIP := treeNode.ServerIdentity.Address.Host()
			node := s.Nodes.GetByIP(serverIP)
			node.ServerIdentity = treeNode.ServerIdentity
			s.Nodes.ServerIdentityToName[treeNode.ServerIdentity.ID] = node.Name
			ServerIdentityToName[treeNode.ServerIdentity] = node.Name
			log.Lvl1("associating", treeNode.ServerIdentity.String(), "to", node.Name)
		}
	}

	return ServerIdentityToName
}
