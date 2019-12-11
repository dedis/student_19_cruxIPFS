package service

import (
	"path/filepath"
	"strconv"

	"go.dedis.ch/onet/v3"
	"go.dedis.ch/onet/v3/log"
)

// Check that *StartIPFSProtocol implements onet.ProtocolInstance
var _ onet.ProtocolInstance = (*ClusterBootstrapProtocol)(nil)

// NewLaunchClustersProtocol initialises the structure for use in one round
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
	log.Lvl3(p.ServerIdentity(), "Starting ClusterBootstrapProtocol")
	return p.SendTo(p.TreeNode(), &ClusterBootstrapAnnounce{})
}

// Dispatch implements the main logic of the protocol. The function is only
// called once. The protocol is considered finished when Dispatch returns and
// Done is called.
func (p *ClusterBootstrapProtocol) Dispatch() error {
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
		bootstrap := info.IP + strconv.Itoa(info.ClusterPort)
		return p.SendToChildren(&ClusterBootstrapAnnounce{
			SenderName: s.Name,
			Bootstrap:  bootstrap,
			Secret:     secret,
		})
	}

	if p.IsLeaf() {
		return p.SendToParent(&ClusterBootstrapReply{true})
	}
	p.SendToChildren(&ann.ClusterBootstrapAnnounce)
	<-p.repliesChan
	if !p.IsRoot() {
		return p.SendToParent(&ClusterBootstrapReply{true})
	}
	p.Ready <- true
	log.Lvl1("Root is done")
	return nil
}
