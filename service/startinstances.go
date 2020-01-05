package service

import (
	"sync"

	"go.dedis.ch/onet/v3"
	"go.dedis.ch/onet/v3/log"
)

// Check that *StartIPFSProtocol implements onet.ProtocolInstance
var _ onet.ProtocolInstance = (*StartInstancesProtocol)(nil)

// NewStartInstancesProtocol initialises the structure for use in one round
func NewStartInstancesProtocol(n *onet.TreeNodeInstance, getServ FnService) (
	onet.ProtocolInstance, error) {
	t := &StartInstancesProtocol{
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
func (p *StartInstancesProtocol) Start() error {
	log.Lvl1("Starting IPFS and IPFS-Cluster instances")
	return p.SendTo(p.TreeNode(), &StartInstancesAnnounce{})
}

// Dispatch implements the main logic of the protocol. The function is only
// called once. The protocol is considered finished when Dispatch returns and
// Done is called.
func (p *StartInstancesProtocol) Dispatch() error {
	defer p.Done()

	s := p.GetService()
	ipfs := IPFSInformation{Name: s.Name, IP: s.MyIP}

	if !p.IsLeaf() {
		// send request to children
		p.SendToChildren(&StartInstancesAnnounce{})

		nodeChan := make(chan NodeInfo, 1)
		go func(c chan NodeInfo) {
			// starting ARAs where node is leader
			c <- NodeInfo{IPFS: ipfs, Clusters: s.startLocalInstances()}
		}(nodeChan)

		// wait for children replies
		ipfsReplies := <-p.repliesChan
		node := <-nodeChan

		p.Nodes = make(map[string]*NodeInfo)
		p.Nodes[s.Name] = &node
		for i := 0; i < len(ipfsReplies); i++ {
			p.Nodes[ipfsReplies[i].Node.IPFS.Name] = ipfsReplies[i].Node
		}

		/*
			if !p.IsRoot() {
				return p.SendToParent(&StartInstancesReply{Node: &node})
			}*/ // root

		p.Ready <- true

		return nil
	}
	// node is a leaf of the global tree

	// starting ARAs where node is leader
	node := NodeInfo{IPFS: ipfs, Clusters: s.startLocalInstances()}
	// return the info to global tree root
	return p.SendToParent(&StartInstancesReply{Node: &node})
}

// startInstances start all ipfs and ipfs cluter instances where each node is
// the root
func (s *Service) startLocalInstances() []ClusterInfo {
	list := make([]ClusterInfo, 0)
	listMutex := sync.Mutex{}
	wg := sync.WaitGroup{}
	// iterate over all ARA trees where the local node is the root
	for _, tree := range s.BinaryTree[s.Name] {
		wg.Add(1)
		go func(t *onet.Tree) {
			pi, err := s.CreateProtocol(StartARAName, t)
			checkErr(err)

			// start the ARA
			pi.Start()
			<-pi.(*StartARAProtocol).Ready

			// append the newly started cluster information to the local list
			// of clusters
			listMutex.Lock()
			list = append(list, pi.(*StartARAProtocol).Info)
			listMutex.Unlock()
			wg.Done()
		}(tree)
	}
	wg.Wait()
	return list
}
