package service

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"go.dedis.ch/onet/v3/log"
)

// StartIPFS starts an IPFS instance for the given service
func (s *Service) StartIPFS() {
	// get config path
	// e.g $GOPATH/src/github.com/dedis/student_19_cruxIPFS/simulation/build
	pwd, err := os.Getwd()
	checkErr(err)
	configPath := filepath.Join(pwd, ConfigsFolder)

	if !LocalSim {
		// if not local, everyone create its own configs directory
		checkErr(CreateEmptyDir(configPath))
	} else if s.Name == Node0 {
		// only node 0, create the empty folder
		checkErr(CreateEmptyDir(configPath))
	}

	if LocalSim {
		s.NodeWg.Done()

		// all nodes wait until the configs folder is created and all nodes reach
		// this stage
		s.NodeWg.Wait()
		time.Sleep(500 * time.Millisecond)
		s.NodeWg.Add(1)
	}

	// set the port range allocated to s
	s.MinPort = BaseHostPort + s.getNodeID()*MaxPortNumberPerHost
	s.MaxPort = s.MinPort + MaxPortNumberPerHost

	// create own config home folder and ipfs config folder
	s.MyIPFSPath = filepath.Join(configPath, s.Name, IPFSFolder)
	checkErr(CreateEmptyDir(s.MyIPFSPath))

	// init ipfs in the desired folder
	exec.Command("ipfs", "-c"+s.MyIPFSPath, "init").Run()

	// edit the ip in the config file
	s.EditIPFSConfig()

	// start ipfs daemon
	go func() {
		exec.Command("ipfs", "-c"+s.MyIPFSPath, "daemon").Run()
		fmt.Println("ipfs at ip", s.Name, "crashed")
	}()
	// wait until it has started
	time.Sleep(IPFSStartupTime)

	if LocalSim {
		s.NodeWg.Done()
		s.NodeWg.Wait()
	}
	if s.Name == Node0 {
		log.Lvl1("All IPFS instances successfully started")
	}
}

// ManageClusters create and start all desired clusters
func (s *Service) ManageClusters() {

}

// LaunchCluster launches a cluster, leader first, and then bootstrap other
// peers until the cluster is complete
func (s *Service) LaunchCluster() {

}
