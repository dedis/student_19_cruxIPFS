package main

import (
	"fmt"
	"strconv"

	"github.com/BurntSushi/toml"
	template "github.com/dedis/student_19_cruxIPFS"
	"go.dedis.ch/onet/app"
	"go.dedis.ch/onet/v3"
	"go.dedis.ch/onet/v3/log"
	"go.dedis.ch/onet/v3/simul/monitor"
)

var mySC onet.SimulationConfig

/*
 * Defines the simulation for the service-template
 */

func init() {
	onet.SimulationRegister("MyTestSimulation", NewSimulationService)
}

// SimulationService only holds the BFTree simulation
type SimulationService struct {
	onet.SimulationBFTree
}

// NewSimulationService returns the new simulation, where all fields are
// initialised using the config-file
func NewSimulationService(config string) (onet.Simulation, error) {
	es := &SimulationService{}
	_, err := toml.Decode(config, es)
	if err != nil {
		return nil, err
	}
	return es, nil
}

// Setup creates the tree used for that simulation
func (s *SimulationService) Setup(dir string, hosts []string) (
	*onet.SimulationConfig, error) {
	app.Copy(dir, "../clean.sh")

	sc := &onet.SimulationConfig{}
	s.CreateRoster(sc, hosts, 2000)

	err := s.CreateTree(sc)
	if err != nil {
		return nil, err
	}
	mySC = *sc
	fmt.Println(sc.Roster)
	fmt.Println(hosts)
	for _, h := range sc.Roster.List {
		fmt.Println(h)
	}
	//StartIPFSDaemon(sc, 2)

	return sc, nil
}

// Node can be used to initialize each node before it will be run
// by the server. Here we call the 'Node'-method of the
// SimulationBFTree structure which will load the roster- and the
// tree-structure to speed up the first round.
func (s *SimulationService) Node(config *onet.SimulationConfig) error {
	index, _ := config.Roster.Search(config.Server.ServerIdentity.ID)
	if index < 0 {
		log.Fatal("Didn't find this node in roster")
	}
	log.Lvl3("Initializing node-index", index)

	fmt.Println(config.Roster)
	for _, h := range config.Roster.List {
		fmt.Println(h)
	}

	//StartIPFSDaemon(config, index)

	return s.SimulationBFTree.Node(config)
}

// Run is used on the destination machines and runs a number of
// rounds
func (s *SimulationService) Run(config *onet.SimulationConfig) error {
	size := config.Tree.Size()
	log.Lvl2("Size is:", size, "rounds:", s.Rounds)
	c := template.NewClient()
	for round := 0; round < s.Rounds; round++ {
		log.Lvl1("Starting round", round)

		round := monitor.NewTimeMeasure("round")
		resp, err := c.Clock(config.Roster)
		log.ErrFatal(err)

		if resp.Time <= 0 {
			log.Fatal("0 time elapsed")
		}
		round.Record()
	}
	return nil
}

// StartIPFSDaemon select ports for each peer and start ipfs on those ports
func StartIPFSDaemon(sc *onet.SimulationConfig, index int) {

	c := template.NewClient()
	// mySC is the SimulConfig created in Setup()
	// I guess I shouldn't use it, but it runs
	identity := sc.Roster.Get(index)

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

	/*
		req2 := template.StartCluster{
			ConfigPath: configPath,
			IP:         ip,
		}

		c2 := template.NewClient()
		reply2 := &template.StartClusterReply{}
		err = c2.SendProtobuf(identity, &req2, reply2)
		if err != nil {
			fmt.Println(err)
		}*/
}
