package service

import (
	"fmt"
	"path/filepath"
	"strconv"

	"go.dedis.ch/onet/v3"
	"go.dedis.ch/onet/v3/log"
)

// Check that *StartIPFSProtocol implements onet.ProtocolInstance
var _ onet.ProtocolInstance = (*StartAllProtocol)(nil)

// NewStartAllProtocol initialises the structure for use in one round
func NewStartAllProtocol(n *onet.TreeNodeInstance, getServ FnService) (
	onet.ProtocolInstance, error) {
	t := &StartAllProtocol{
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
func (p *StartAllProtocol) Start() error {
	log.Lvl2("Starting an ARA with root", p.GetService().Name)
	return p.SendTo(p.TreeNode(), &StartAllAnnounce{})
}

// Dispatch implements the main logic of the protocol. The function is only
// called once. The protocol is considered finished when Dispatch returns and
// Done is called.
func (p *StartAllProtocol) Dispatch() error {
	defer p.Done()

	ann := <-p.announceChan
	s := p.GetService()

	apiIPFSAddr := IPVersion + s.MyIPFS[0].IP +
		TransportProtocol + strconv.Itoa(s.MyIPFS[0].APIPort) // 5001

	if p.IsRoot() {
		// generate secret
		secret := genSecret()

		// starting IPFS and cluster instance
		instance := s.StartIPFSAndCluster(s.Name, secret, "")

		// adding information to cluster info
		p.Info = ClusterInfo{
			Leader:    s.Name,
			Secret:    secret,
			Instances: make([]ClusterInstance, 1),
		}
		p.Info.Instances[0] = *instance

		if len(p.TreeNodeInstance.Children()) > 0 {
			bootstrap := instance.IP + strconv.Itoa(instance.ClusterPort)
			p.SendToChildren(&ClusterBootstrapAnnounce{
				SenderName: s.Name,
				Bootstrap:  bootstrap,
				Secret:     secret,
			})
			replies := <-p.repliesChan
			for i := 0; i < len(replies); i++ {
				p.Info.Instances = append(p.Info.Instances,
					*replies[i].Cluster...)
			}
		}
		p.Info.Size = len(p.Info.Instances)
		p.Ready <- true
		return nil
	} else if !p.IsLeaf() {
		p.SendToChildren(&ann.ClusterBootstrapAnnounce)
		replies := <-p.repliesChan

		clusterPath := filepath.Join(s.ConfigPath,
			ClusterFolderPrefix+ann.SenderName+"-"+ann.Secret)

		// bootstrap peer
		cluster, err := s.SetupClusterSlave(clusterPath, ann.Bootstrap,
			ann.Secret, apiIPFSAddr, DefaultReplMin, DefaultReplMax)
		if err != nil {
			fmt.Println("Error slave:", err)
		}
		instances := make([]ClusterInstance, 0)
		for _, r := range replies {
			instances = append(instances, *r.Cluster...)
		}
		instances = append(instances, *cluster)
		return p.SendToParent(&ClusterBootstrapReply{Cluster: &instances})
	}
	// leaf
	clusterPath := filepath.Join(s.ConfigPath,
		ClusterFolderPrefix+ann.SenderName+"-"+ann.Secret)

	// bootstrap peer
	cluster, err := s.SetupClusterSlave(clusterPath, ann.Bootstrap, ann.Secret,
		apiIPFSAddr, DefaultReplMin, DefaultReplMax)
	if err != nil {
		fmt.Println("Error slave:", err)
	}
	return p.SendToParent(&ClusterBootstrapReply{
		Cluster: &[]ClusterInstance{*cluster}})
}
