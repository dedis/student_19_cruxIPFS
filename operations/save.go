package operations

import (
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/dedis/student_19_cruxIPFS/service"
	"github.com/ipfs/ipfs-cluster/api/rest/client"
	ma "github.com/multiformats/go-multiaddr"

	"go.dedis.ch/onet/v3/log"
)

// SaveState save ipfs and ipfs-cluster instances state
func SaveState(filename string, instances map[string]*service.NodeInfo) {

	// identifiers for saving state

	str := ""
	str += "\n##### IPFS instances #####\n\n"

	ids := []string{"IP", "Swarm Port", "API Port", "Gateway Port",
		"Cluster number", "Clusters"}
	ids = alignIds(ids, 2)
	cids := []string{"Cluster leader", "Secret", "Size", "Hosts"}
	cids = alignIds(cids, 4)
	pids := []string{"IP", "RestAPI Port", "IPFS Proxy Port", "Cluster Port"}
	pids = alignIds(pids, 8)

	for _, i := range instances {
		ii := i.IPFS
		str += ii.Name + ":"
		str += ids[0] + ii.IP
		str += ids[1] + strconv.Itoa(ii.SwarmPort)
		str += ids[2] + strconv.Itoa(ii.APIPort)
		str += ids[3] + strconv.Itoa(ii.GatewayPort)
		str += ids[4] + strconv.Itoa(len(i.Clusters))
		str += ids[5]

		for _, c := range i.Clusters {
			str += cids[0] + c.Leader
			str += cids[1] + c.Secret
			str += cids[2] + strconv.Itoa(c.Size)
			str += cids[3]
			for _, ci := range c.Instances {
				str += "\n      " + ci.HostName
				str += pids[0] + ci.IP
				str += pids[1] + strconv.Itoa(ci.RestAPIPort)
				str += pids[2] + strconv.Itoa(ci.IPFSProxyPort)
				str += pids[3] + strconv.Itoa(ci.ClusterPort)
				str += "\n"
			}
		}
		str += "\n"
	}
	log.Lvl1(str)
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
		ids[i] = "\n" + strings.Repeat(" ", indent) + t +
			strings.Repeat(" ", max-len(t)) + " : "
	}
	return ids
}

// LoadClusterInstances load ipfs-cluster instances for each node
func LoadClusterInstances(filename string) map[string]*Node {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	str := string(bytes)
	nodes := make(map[string]*Node)
	lines := strings.Split(str, "\n")
	nLines := len(lines)
	for l := 0; l < nLines; l++ {
		line := lines[l]
		// cluster instance described next
		if strings.Contains(line, service.NodeName) &&
			!strings.Contains(line, ":") {

			size := 0
			for m := l; m >= 0; m-- {
				if strings.Contains(lines[m], "Size") {
					split := strings.Split(lines[m], " ")
					size, err = strconv.Atoi(split[len(split)-1])
					checkErr(err)
					break
				}
			}
			if size <= 2 {
				continue
			}

			secret := ""
			for m := l; m >= 0; m-- {
				if strings.Contains(lines[m], "Secret") {
					split := strings.Split(lines[m], " ")
					secret = split[len(split)-1]
					checkErr(err)
					break
				}
			}
			if secret == "" {
				panic("secret not found")
			}

			split := strings.Split(line, " ")
			name := split[len(split)-1]

			split = strings.Split(lines[l+1], " ")
			ip := split[len(split)-1]

			split = strings.Split(lines[l+2], " ")
			api := ip + split[len(split)-1]
			apiAddr, err := ma.NewMultiaddr(api)
			if err != nil {
				panic(err)
			}
			split = strings.Split(lines[l+3], " ")
			proxy := ip + split[len(split)-1]
			proxyAddr, err := ma.NewMultiaddr(proxy)
			if err != nil {
				panic(err)
			}

			conf := client.Config{
				APIAddr:   apiAddr,
				ProxyAddr: proxyAddr,
			}
			c, err := client.NewDefaultClient(&conf)
			if err != nil {
				panic(err)
			}

			if _, ok := nodes[name]; ok {
				nodes[name].Clients = append(nodes[name].Clients, c)
				nodes[name].Secrets = append(nodes[name].Secrets, secret)
			} else {
				// node don't exist, create it and append address
				cli := make([]client.Client, 1)
				cli[0] = c
				secrets := make([]string, 1)
				secrets[0] = secret
				nodes[name] = &Node{
					Name:    name,
					Clients: cli,
					Secrets: secrets,
				}

			}
		}
	}
	return nodes
}
