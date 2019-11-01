package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"

	"github.com/BurntSushi/toml"
	template "github.com/dedis/student_19_cruxIPFS"
	"go.dedis.ch/onet/app"
	"go.dedis.ch/onet/v3"
	"go.dedis.ch/onet/v3/log"
	"go.dedis.ch/onet/v3/simul/monitor"
)

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

	fmt.Println(hosts)

	sc := &onet.SimulationConfig{}
	s.CreateRoster(sc, hosts, 2000)

	fmt.Println(sc.Roster.List)
	err := s.CreateTree(sc)
	if err != nil {
		return nil, err
	}

	app.Copy(dir, "../clean.sh")

	// clean all ipfs processes
	o, err := exec.Command("../clean.sh").Output()
	if err != nil {
		fmt.Println(err)
	}
	// remove the config dir
	err = os.RemoveAll(template.DefaultConfigPath)
	if err != nil {
		fmt.Println(err)
	}
	// create an empty config dir
	err = os.MkdirAll(template.DefaultConfigPath, 0777)
	if err != nil {
		fmt.Println(err)
	}
	log.Lvl1(string(o))
	/*
		c := template.NewClient()
		reply := &template.DoSetupReply{}
		req := template.DoSetup{Path: "hey"}
		id := sc.Server.ServerIdentity
		fmt.Println(c, reply, req, id)

		/*
			err = c.SendProtobuf(sc.Server.ServerIdentity, &req, reply)
			if err != nil {
				fmt.Println(err)
			}*/
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

	// should be set to a const
	configPath := "/home/guillaume/ipfs_test/myfolder/Node" +
		strconv.Itoa(index)
	// get node id
	identity := config.Roster.Get(index)

	c := template.NewClient()
	reply := &template.StartIPFSReply{}
	// create start ipfs request
	req := template.StartIPFS{
		ConfigPath: configPath,
		NodeID:     index,
		PortMin:    14000,
		PortMax:    15000,
	}
	err := c.SendProtobuf(identity, &req, reply)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Swarm:", reply.Ports.Swarm, "API:", reply.Ports.API,
		"Gateway", reply.Ports.Gateway)

	// !!!!!! We should have the same secret for all hosts in the same cluster
	resp, err := c.GenSecret(config.Roster)
	log.ErrFatal(err)

	replyC := &template.StartClusterReply{}
	reqC := template.StartCluster{
		ConfigPath: configPath,
		NodeID:     index,
		ClusterID:  0,
		PortMax:    15000,
		PortMin:    14000,
		Secret:     resp.Secret,
		Peername:   "Peer" + strconv.Itoa(index) + "_0",
	}
	err = c.SendProtobuf(identity, &reqC, replyC)

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
