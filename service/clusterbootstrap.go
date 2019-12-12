package service

import (
	"fmt"
	"path/filepath"
	"strconv"

	"go.dedis.ch/onet/v3"
	"go.dedis.ch/onet/v3/log"
)

// Check that *StartIPFSProtocol implements onet.ProtocolInstance
var _ onet.ProtocolInstance = (*ClusterBootstrapProtocol)(nil)

// NewClusterBootstrapProtocol initialises the structure for use in one round
func NewClusterBootstrapProtocol(n *onet.TreeNodeInstance, getServ ServiceFn) (
	onet.ProtocolInstance, error) {
	t := &ClusterBootstrapProtocol{
		TreeNodeInstance: n,
		Ready:            make(chan bool),
		GetService:       getServ,
	}
	if err := n.RegisterChannels(&t.announceChan, &t.repliesChan); err != nil {
		return nil, err
	}
	return t, nil
}

// Start sends the Announce-message to all children
func (p *ClusterBootstrapProtocol) Start() error {
	log.Lvl1(p.GetService().Name, "starting an ARA")
	return p.SendTo(p.TreeNode(), &ClusterBootstrapAnnounce{})
}

// Dispatch implements the main logic of the protocol. The function is only
// called once. The protocol is considered finished when Dispatch returns and
// Done is called.
func (p *ClusterBootstrapProtocol) Dispatch() error {
	// the tree should have a branching factor of n-1, with n being the number
	// of nodes
	defer p.Done()

	ann := <-p.announceChan
	s := p.GetService()

	if p.IsRoot() {
		// create cluster dir
		clusterPath := filepath.Join(s.ConfigPath, ClusterFolderPrefix+s.Name)

		// start cluster leader
		secret, info, err := s.SetupClusterLeader(clusterPath, DefaultReplMin,
			DefaultReplMax)
		checkErr(err)

		p.Info = ClusterInfo{
			Leader: s.Name,
			Secret: secret,
			Size:   len(p.TreeNodeInstance.Children()) + 1,
			Instances: make([]ClusterInstance,
				len(p.TreeNodeInstance.Children())+1),
		}
		p.Info.Instances[0] = *info

		if len(p.TreeNodeInstance.Children()) > 0 {
			bootstrap := info.IP + strconv.Itoa(info.ClusterPort)
			p.SendToChildren(&ClusterBootstrapAnnounce{
				SenderName: s.Name,
				Bootstrap:  bootstrap,
				Secret:     secret,
			})
			replies := <-p.repliesChan
			for i := 0; i < len(replies); i++ {
				p.Info.Instances[i+1] = *replies[i].Cluster
			}
		}
		p.Ready <- true
		return nil
	}
	// leaf
	clusterPath := filepath.Join(s.ConfigPath,
		ClusterFolderPrefix+ann.SenderName+"-"+ann.Secret)

	// bootstrap peer
	cluster, err := s.SetupClusterSlave(clusterPath, ann.Bootstrap, ann.Secret,
		DefaultReplMin, DefaultReplMax)
	if err != nil {
		fmt.Println("Error slave:", err)
	}
	return p.SendToParent(&ClusterBootstrapReply{Cluster: cluster})
}
