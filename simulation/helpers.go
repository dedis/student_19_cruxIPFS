package main

import (
	"io/ioutil"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/dedis/student_19_cruxIPFS/service"
)

// SetNodePaths set the node paths for remote and local node files
func SetNodePaths(n int) {
	NODEPATHREMOTE = NODEPATHNAME + strconv.Itoa(n) + ".txt"
	NODEPATHLOCAL = filepath.Join("..", DATAFOLDER, NODEPATHREMOTE)
}

func saveState(filename string, ipfs []service.IPFSInformation,
	clusters []service.ClusterInfo) {

	// identifiers for saving state
	ids := []string{"IP", "Swarm Port", "API Port", "Gateway Port"}
	ids = alignIds(ids, 2)

	str := ""
	str += "\n##### IPFS instances #####\n\n"
	for _, ii := range ipfs {
		str += ii.Name + ":"
		str += ids[0] + ii.IP
		str += ids[1] + strconv.Itoa(ii.SwarmPort)
		str += ids[2] + strconv.Itoa(ii.APIPort)
		str += ids[3] + strconv.Itoa(ii.GatewayPort)
		str += "\n\n"
	}

	str += "\n##### Clusters #####\n"

	cids := []string{"Cluster leader", "Secret", "Size", "Hosts"}
	pids := []string{"IP", "RestAPI Port", "IPFS Proxy Port", "Cluster Port"}
	cids = alignIds(cids, 0)
	pids = alignIds(pids, 2)

	for _, c := range clusters {
		str += cids[0] + c.Leader
		str += cids[1] + c.Secret
		str += cids[2] + strconv.Itoa(c.Size)
		str += cids[3]
		for _, ci := range c.Instances {
			str += "\n  " + ci.HostName
			str += pids[0] + ci.IP
			str += pids[1] + strconv.Itoa(ci.RestAPIPort)
			str += pids[2] + strconv.Itoa(ci.IPFSProxyPort)
			str += pids[3] + strconv.Itoa(ci.ClusterPort)
			str += "\n"
		}
	}

	ioutil.WriteFile(filename, []byte(str), 0)
}

func alignIds(ids []string, indent int) []string {
	max := 0
	for _, t := range ids {
		if len(t) > max {
			max = len(t)
		}
	}
	for i, t := range ids {
		ids[i] = "\n"+strings.Repeat(" ", indent) + t +
		 strings.Repeat(" ", max-len(t)) + " : "
	}
	return ids
}

func loadState() {

}
