package service

import (
	"sync"

	"go.dedis.ch/onet/v3"
	"go.dedis.ch/onet/v3/log"
)

// Check that *StartIPFSProtocol implements onet.ProtocolInstance
var _ onet.ProtocolInstance = (*StartIPFSProtocol)(nil)

// NewStartIPFSProtocol initialises the structure for use in one round
func NewStartIPFSProtocol(n *onet.TreeNodeInstance, getServ ServiceFn) (
	onet.ProtocolInstance, error) {
	t := &StartIPFSProtocol{
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
func (p *StartIPFSProtocol) Start() error {
	log.Lvl1("Starting IPFS instances")
	return p.SendTo(p.TreeNode(), &StartIPFSAnnounce{})
}

// Dispatch implements the main logic of the protocol. The function is only
// called once. The protocol is considered finished when Dispatch returns and
// Done is called.
func (p *StartIPFSProtocol) Dispatch() error {
	defer p.Done()

	ann := <-p.announceChan

	if !p.IsLeaf() {
		// send request to children
		p.SendToChildren(&ann.StartIPFSAnnounce)

		// waitgroup to wait own ipfs instance
		wg := sync.WaitGroup{}
		wg.Add(1)
		go func() {
			// start ipfs
			p.GetService().StartIPFS()
			wg.Done()
		}()

		// wait for children replies
		<-p.repliesChan
		wg.Wait()

		if !p.IsRoot() {
			return p.SendToParent(&StartIPFSReply{true})
		} // root
		p.Ready <- true
		log.Lvl1("All IPFS instances started successfully")
		return nil

	}
	// node is a leaf
	// start ipfs
	p.GetService().StartIPFS()
	// send ok to parent when it's done
	return p.SendToParent(&StartIPFSReply{true})
}
