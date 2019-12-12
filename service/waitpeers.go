package service

import (
	"fmt"
	"time"

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
	log.Lvl3(p.ServerIdentity(), "Starting WaitpeersProtocol")
	return p.SendTo(p.TreeNode(), &WaitpeersAnnounce{"ready!"})
}

// Dispatch implements the main logic of the protocol. The function is only
// called once. The protocol is considered finished when Dispatch returns and
// Done is called.
func (p *WaitpeersProtocol) Dispatch() error {
	defer p.Done()

	fmt.Println("Entering protocol")
	ann := <-p.announceChan

	if p.IsLeaf() {
		if true {
			log.Lvl1("Oopsi")
			time.Sleep(3 * time.Second)
		}
		log.Lvl1("Leaf is done")
		return p.SendToParent(&WaitpeersReply{true})
	}
	p.SendToChildren(&ann.WaitpeersAnnounce)
	log.Lvl1("Sent to children")
	<-p.repliesChan
	log.Lvl1("Got a reply")
	if !p.IsRoot() {
		log.Lvl1("Middle is done")
		return p.SendToParent(&WaitpeersReply{true})
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
