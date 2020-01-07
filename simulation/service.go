package main

import (
	"fmt"
	"os/exec"
	"time"

	"github.com/BurntSushi/toml"

	cruxIPFS "github.com/dedis/student_19_cruxIPFS"
	"github.com/dedis/student_19_cruxIPFS/operations"
	"github.com/dedis/student_19_cruxIPFS/service"

	"go.dedis.ch/onet/v3"
	"go.dedis.ch/onet/v3/app"
	"go.dedis.ch/onet/v3/log"
)

func init() {
	onet.SimulationRegister(service.ServiceName, NewIPFSSimulation)
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

// Setup the IPFSSimulation, copies files to remote host (deterlab), load
// simulation parameters and create roster and config tree
// This function is run on a single node
func (s *IPFSSimulation) Setup(dir string, hosts []string) (
	*onet.SimulationConfig, error) {

	app.Copy(dir, prescriptLocation)
	app.Copy(dir, nodesLocation)
	app.Copy(dir, ipfsLocation)
	app.Copy(dir, ipfsClusterLocation)
	app.Copy(dir, ipfsCtlLocation)
	app.Copy(dir, simdetailsLocation)

	sc := &onet.SimulationConfig{}
	s.CreateRoster(sc, hosts, 2000)
	err := s.CreateTree(sc)
	if err != nil {
		return nil, err
	}
	return sc, nil
}

// Node is run on all nodes, it reads nodes information (mostly the level for
// ARA generation), and initialize the service (computing/loading ping distance,
// generating ARAs, starting ipfs and clusters)
func (s *IPFSSimulation) Node(config *onet.SimulationConfig) error {

	s.ReadNodeInfo(false, *config)

	parseParams()

	mymap := s.initializeMaps(config, true)

	myService := config.GetService(service.ServiceName).(*service.Service)

	serviceReq := &service.InitRequest{
		Nodes:                s.Nodes.All,
		ServerIdentityToName: mymap,
		OnetTree:             config.Tree,
		Roster:               config.Roster,
		Cruxified:            cruxified,
		ComputePings:         computePings,
		Mode:                 mode,
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

func (s *IPFSSimulation) Run1(config *onet.SimulationConfig) error {

	/*
		o, err := exec.Command("bash", "-c", cmd).Output()
		fmt.Println(string(o))
		if err != nil {
			fmt.Println("Error:", err)
		}
	*/
	go func() {
		cmd := "ipfs daemon"
		o, err := exec.Command("bash", "-c", cmd).Output()
		fmt.Println(string(o))
		if err != nil {
			fmt.Println("Error:", err)
		}
	}()

	time.Sleep(10 * time.Second)

	return nil
}

// Run is run on a single node. Execute performance tests and output results to
// stdout, output needs to be parsed by an external script
func (s *IPFSSimulation) Run(config *onet.SimulationConfig) error {

	myService := config.GetService(service.ServiceName).(*service.Service)
	/*
		pi, err := myService.CreateProtocol(service.StartIPFSName, config.Tree)
		if err != nil {
			fmt.Println(err)
		}
		pi.Start()

		<-pi.(*service.StartIPFSProtocol).Ready

		operations.SaveState(cruxIPFS.SaveFile,
			pi.(*service.StartIPFSProtocol).Nodes)
	*/

	pi, err := myService.CreateProtocol(service.StartInstancesName, config.Tree)
	if err != nil {
		fmt.Println(err)
	}
	pi.Start()

	<-pi.(*service.StartInstancesProtocol).Ready

	operations.SaveState(cruxIPFS.SaveFile,
		pi.(*service.StartInstancesProtocol).Nodes)

	// wait for some time for clusters to converge
	time.Sleep(2 * time.Minute)
	operations.Test2(nOps, len(myService.Nodes.All))

	log.Lvl1("Done")
	return nil
}
