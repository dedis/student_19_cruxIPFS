package protocol

import (
	"go.dedis.ch/onet/v3"
	"go.dedis.ch/onet/v3/log"
)

// Check that *WaitpeersProtocol implements onet.ProtocolInstance
var _ onet.ProtocolInstance = (*WaitpeersProtocol)(nil)

// NewWaitpeersProtocol initialises the structure for use in one round
func NewWaitpeersProtocol(n *onet.TreeNodeInstance) (onet.ProtocolInstance, error) {
	t := &WaitpeersProtocol{
		TreeNodeInstance: n,
		Ready:            make(chan bool),
	}
	if err := n.RegisterChannels(&t.announceChan, &t.repliesChan); err != nil {
		return nil, err
	}
	return t, nil
}

// Start sends the Announce-message to all children
func (p *WaitpeersProtocol) Start() error {
	log.Lvl3(p.ServerIdentity(), "Starting TemplateProtocol")
	return p.SendTo(p.TreeNode(), &Announce{"cothority rulez!"})
}

// Dispatch implements the main logic of the protocol. The function is only
// called once. The protocol is considered finished when Dispatch returns and
// Done is called.
func (p *WaitpeersProtocol) Dispatch() error {
	defer p.Done()

	/*
		ann := <-p.announceChan
		if p.IsLeaf() {
			return p.SendToParent(&Reply{1})
		}
		p.SendToChildren(&ann.Announce)

		replies := <-p.repliesChan
		n := 1
		for _, c := range replies {
			n += c.ChildrenCount
		}

		if !p.IsRoot() {
			return p.SendToParent(&Reply{n})
		}

		p.ChildCount <- n
	*/
	return nil
}
