package service

import (
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
	log.Lvl3(p.ServerIdentity(), "Starting WaitpeersProtocol")
	return p.SendTo(p.TreeNode(), &StartIPFSAnnounce{"ready!"})
}

// Dispatch implements the main logic of the protocol. The function is only
// called once. The protocol is considered finished when Dispatch returns and
// Done is called.
func (p *StartIPFSProtocol) Dispatch() error {
	defer p.Done()

	ann := <-p.announceChan

	service := p.GetService()
	service.getNodeID()

	if p.IsLeaf() {
		return p.SendToParent(&StartIPFSReply{true})
	}
	p.SendToChildren(&ann.StartIPFSAnnounce)
	<-p.repliesChan
	if !p.IsRoot() {
		return p.SendToParent(&StartIPFSReply{true})
	}
	p.Ready <- true
	log.Lvl1("Root is done")

	/*
		if !p.IsLeaf() {
			err := p.SendToChildren(&ann.WaitpeersAnnounce)
			log.Lvl1(p.ServerIdentity(), "sent request")
			checkErr(err)
		}
		if !p.IsRoot() {
			log.Lvl1(p.ServerIdentity(), "sending response")
			time.Sleep(2 * time.Second)
			err := p.SendToParent(&WaitpeersReply{true})
			checkErr(err)
			<-p.repliesChan
			log.Lvl1("Child is done")
			p.Ready <- true
			p.Done()
		} else {
			<-p.repliesChan
			p.SendToChildren(&WaitpeersReply{true})
			log.Lvl1("Root is done")
			p.Ready <- true
			p.Done()
		}
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
