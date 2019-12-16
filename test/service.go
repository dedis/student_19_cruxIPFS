package main

import (
	"github.com/BurntSushi/toml"
	"go.dedis.ch/onet/v3"
	"go.dedis.ch/onet/v3/log"

	"github.com/dedis/student_19_cruxIPFS/operations"
)

/*
 * Defines the simulation for the service-template
 */

func init() {
	onet.SimulationRegister("TestService", NewSimulationService)
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
func (s *SimulationService) Node(config *onet.SimulationConfig) error {
	/*
		index, _ := config.Roster.Search(config.Server.ServerIdentity.ID)
		if index < 0 {
			log.Fatal("Didn't find this node in roster")
		}
		log.Lvl3("Initializing node-index", index)
	*/
	return s.SimulationBFTree.Node(config)
}

// Run is used on the destination machines and runs a number of
// rounds
func (s *SimulationService) Run(config *onet.SimulationConfig) error {
	f1 := "file1.txt"
	operations.NewFile(f1)
	n, t1 := operations.Write("node_0", f1)
	log.Lvl1(n, "written in", t1)
	t2 := operations.Read("node_1", n)

	log.Lvl1("Write time:", t1, "Read time:", t2)

	f1 = "file2.txt"
	operations.NewFile(f1)
	n, t1 = operations.Write("node_6", f1)
	t2 = operations.Read("node_4", n)

	log.Lvl1("Write time:", t1, "Read time:", t2)
	return nil
}
