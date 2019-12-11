package service

import (
	"errors"
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

	fmt.Println("HEEEEY")

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

	// set the port range allocated to s
	fmt.Println("getting node id")
	s.MinPort = BaseHostPort + s.getNodeID()*MaxPortNumberPerHost
	s.MaxPort = s.MinPort + MaxPortNumberPerHost

	fmt.Println("ports done")
	s.ConfigPath = filepath.Join(configPath, s.Name)

	// create own config home folder and ipfs config folder
	s.MyIPFSPath = filepath.Join(s.ConfigPath, IPFSFolder)
	checkErr(CreateEmptyDir(s.MyIPFSPath))

	// init ipfs in the desired folder
	exec.Command("ipfs", "-c"+s.MyIPFSPath, "init").Run()

	fmt.Println("go to edit")
	// edit the ip in the config file
	EditIPFSConfig(s)

	// start ipfs daemon
	go func() {
		exec.Command("ipfs", "-c"+s.MyIPFSPath, "daemon").Run()
		fmt.Println("ipfs at ip", s.Name, "crashed")
	}()
	// wait until it has started
	time.Sleep(IPFSStartupTime)
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
	req, ok := env.Msg.(*ReqBootstrapCluster)
	if !ok {
		log.Error(s.ServerIdentity(), "failed to cast to ReqPings")
		return errors.New(s.ServerIdentity().String() + " failed to cast to ReqPings")
	}
	clusterPath := filepath.Join(s.ConfigPath, ClusterFolderPrefix+req.SenderName)

	// create cluster dir

	_, err := s.SetupClusterSlave(clusterPath, req.Bootstrap, req.Secret,
		DefaultReplMin, DefaultReplMax)
	if err != nil {
		fmt.Println("Error slave:", err)
	}

	// bootstrap peer

	requesterIdentity := s.Nodes.GetByName(req.SenderName).ServerIdentity
	e := s.SendRaw(requesterIdentity, &ReplyBootstrapCluster{
		SenderName: s.Name})
	if e != nil {
		panic(e)
	}
	return e
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
	if s.Name == Node0 {
		log.Lvl1("All ipfs-cluster-service instances successfully started")
	}
}
