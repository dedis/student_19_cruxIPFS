package service

import (
	"go.dedis.ch/onet/v3"
	"go.dedis.ch/onet/v3/log"
)

// Check that *StartIPFSProtocol implements onet.ProtocolInstance
var _ onet.ProtocolInstance = (*LaunchClustersProtocol)(nil)

// NewLaunchClustersProtocol initialises the structure for use in one round
func NewLaunchClustersProtocol(n *onet.TreeNodeInstance, getServ ServiceFn) (
	onet.ProtocolInstance, error) {
	t := &LaunchClustersProtocol{
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
func (p *LaunchClustersProtocol) Start() error {
	log.Lvl3(p.ServerIdentity(), "Starting LaunchClustersProtocol")
	return p.SendTo(p.TreeNode(), &LaunchClustersAnnounce{"ready!"})
}

// Dispatch implements the main logic of the protocol. The function is only
// called once. The protocol is considered finished when Dispatch returns and
// Done is called.
func (p *LaunchClustersProtocol) Dispatch() error {
	defer p.Done()

	ann := <-p.announceChan

	if p.IsLeaf() {
		return p.SendToParent(&LaunchClustersReply{true})
	}
	p.SendToChildren(&ann.LaunchClustersAnnounce)
	<-p.repliesChan
	if !p.IsRoot() {
		return p.SendToParent(&LaunchClustersReply{true})
	}
	p.Ready <- true
	log.Lvl1("Root is done")
	return nil
}
