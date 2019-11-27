package service

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"time"

	"go.dedis.ch/onet/v3/log"
	"go.dedis.ch/onet/v3/network"
)

// StartIPFS starts an IPFS instance for the given service
func (s *Service) StartIPFS() {
	// get config path
	// e.g $GOPATH/src/github.com/dedis/student_19_cruxIPFS/simulation/build
	pwd, err := os.Getwd()
	checkErr(err)
	s.ConfigPath = filepath.Join(pwd, ConfigsFolder)

	if !LocalSim {
		// if not local, everyone create its own configs directory
		checkErr(CreateEmptyDir(s.ConfigPath))
	} else if s.Name == Node0 {
		// only node 0, create the empty folder
		checkErr(CreateEmptyDir(s.ConfigPath))
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
	s.MyIPFSPath = filepath.Join(s.ConfigPath, s.Name, IPFSFolder)
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

// ExecReqIPFSInfo sends own IPFS instance information to the node asking for it
func (s *Service) ExecReqIPFSInfo(env *network.Envelope) error {
	return nil
}

// ExecReplyIPFSInfo recieving IPFS information of another IPFS instance, and
// stores it in its own table to be able to
func (s *Service) ExecReplyIPFSInfo(env *network.Envelope) error {
	return nil
}

func (s *Service) ExecReqBootstrapCluster(env *network.Envelope) error {
	// create cluster dir

	// bootstrap peer
	return nil
}

func (s *Service) ExecReplyBootstrapCluster(env *network.Envelope) error {
	// ack cluster is ready
	return nil
}

// ManageClusters create and start all desired clusters
func (s *Service) ManageClusters() {
	// iterate over all nodes
	for _, n := range s.Nodes.All {
		if s.Name == n.Name {
			// what if a node is alone in the cluster it is a root ???
			if len(n.Cluster) > 0 {
				cluster := make([]string, len(n.Cluster))
				i := 0
				for s := range n.Cluster {
					cluster[i] = s
					i++
				}
				s.LaunchCluster(cluster)
			}
			break
		}
	}
}

// LaunchCluster launches a cluster, leader first, and then bootstrap other
// peers until the cluster is complete
func (s *Service) LaunchCluster(nodes []string) {
	// create cluster dir
	clusterPath := filepath.Join(s.ConfigPath, ClusterFolderPrefix+s.Name)
	checkErr(CreateEmptyDir(clusterPath))

	// start cluster leader
	secret, p, err := s.SetupClusterLeader(clusterPath, DefaultReplMin,
		DefaultReplMax)
	checkErr(err)
	bootstrap := p.IP + strconv.Itoa(p.ClusterPort)

	// send init request
	// iterate over all nodes
	for _, node := range s.Nodes.All {
		// iterate over all nodes in cluster
		for _, n := range nodes {
			// if name is a match
			if node.Name == n {
				// send bootstrap request
				// create bootstrap request
				req := ReqBootstrapCluster{
					SenderName: s.Name,
					Bootstrap:  bootstrap,
					Secret:     secret,
				}
				// send request
				e := s.SendRaw(node.ServerIdentity, &req)
				if e != nil {
					panic(e)
				}
				break
			}
		}
	}

}
