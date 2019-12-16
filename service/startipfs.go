package service

import (
	"sync"

	"go.dedis.ch/onet/v3"
	"go.dedis.ch/onet/v3/log"
)

// Check that *StartIPFSProtocol implements onet.ProtocolInstance
var _ onet.ProtocolInstance = (*StartIPFSProtocol)(nil)

// NewStartIPFSProtocol initialises the structure for use in one round
func NewStartIPFSProtocol(n *onet.TreeNodeInstance, getServ FnService) (
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
	s := p.GetService()

	if !p.IsLeaf() {
		// send request to children
		p.SendToChildren(&StartIPFSAnnounce{Message: "IPFS"})

		// waitgroup to wait own ipfs instance
		wg := sync.WaitGroup{}
		wg.Add(1)
		go func() {
			// start ipfs
			s.StartIPFS()
			wg.Done()
		}()

		// wait for children replies
		ipfsReplies := <-p.repliesChan
		wg.Wait()

		p.Nodes = make(map[string]*NodeInfo)
		p.Nodes[s.Name] = &NodeInfo{
			IPFS:     s.MyIPFS,
			Clusters: make([]ClusterInfo, 0),
		}
		for i := 0; i < len(ipfsReplies); i++ {
			p.Nodes[ipfsReplies[i].IPFS.Name] = &NodeInfo{
				IPFS:     *ipfsReplies[i].IPFS,
				Clusters: make([]ClusterInfo, 0),
			}
		}

		if !p.IsRoot() {
			return p.SendToParent(&StartIPFSReply{IPFS: &s.MyIPFS})
		} // root
		log.Lvl1("All IPFS instances started successfully")

		// start cluster instances on the root
		p.SendToChildren(&StartIPFSAnnounce{Message: "Clusters"})

		wg.Add(1)
		go func() {
			p.Nodes[s.Name].Clusters = append(p.Nodes[s.Name].Clusters,
				s.startClusters()...)
			wg.Done()
		}()
		// wait for children replies
		replies := <-p.repliesChan
		wg.Wait()
		for _, r := range replies {
			for _, c := range *(r.Clusters) {
				p.Nodes[c.Leader].Clusters =
					append(p.Nodes[c.Leader].Clusters, c)
			}
		}
		p.Ready <- true

		return nil

	}
	// node is a leaf
	if ann.Message == "IPFS" {
		// start ipfs
		s.StartIPFS()
		// send ok to parent when it's done
		p.SendToParent(&StartIPFSReply{IPFS: &s.MyIPFS})
		ann = <-p.announceChan
		if ann.Message == "Clusters" {
			info := s.startClusters()
			return p.SendToParent(&StartIPFSReply{Clusters: &info})
		}
	}
	return nil
}

func (s *Service) startClusters() []ClusterInfo {
	list := make([]ClusterInfo, 0)
	listMutex := sync.Mutex{}
	wg := sync.WaitGroup{}
	// iterate over all ARA trees where node is the root
	for _, tree := range s.BinaryTree[s.Name] {
		wg.Add(1)
		go func(t *onet.Tree) {
			pi, err := s.CreateProtocol(ClusterBootstrapName, t)
			checkErr(err)
			pi.Start()
			<-pi.(*ClusterBootstrapProtocol).Ready
			listMutex.Lock()
			list = append(list, pi.(*ClusterBootstrapProtocol).Info)
			listMutex.Unlock()
			wg.Done()
		}(tree)
	}
	wg.Wait()
	return list
}
