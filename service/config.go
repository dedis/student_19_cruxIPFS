package service

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"os/exec"
	"strconv"
	"time"

	"go.dedis.ch/onet/v3/log"
)

// EditIPFSConfig edit the ipfs configuration file (mainly the ip)
func EditIPFSConfig(s *Service) {
	addr := IPVersion + s.MyIP + TransportProtocol

	// select available ports
	ports, err := GetNextAvailablePorts(s.MinPort, s.MaxPort, IPFSPortNumber)
	checkErr(err)

	// [\"/ip4/0.0.0.0/tcp/5001\", \"/ip6/::/tcp/5001\"]
	//swarmList := []string{addr + SwarmPort}
	Swarm := MakeJSONArray([]string{addr +
		strconv.Itoa((*ports)[0])})

	// /ip4/127.0.0.1/tcp/5001
	API := MakeJSONElem(addr + strconv.Itoa((*ports)[1]))
	// /ip4/127.0.0.1/tcp/8080
	Gateway := MakeJSONElem(addr + strconv.Itoa((*ports)[2]))

	EditIPFSField(s.MyIPFSPath, "Addresses.API", API)
	EditIPFSField(s.MyIPFSPath, "Addresses.Gateway", Gateway)
	EditIPFSField(s.MyIPFSPath, "Addresses.Swarm", Swarm)

	// filling my IPFS info
	s.MyIPFS = IPFSInformation{
		IP:          s.MyIP,
		SwarmPort:   (*ports)[0],
		APIPort:     (*ports)[1],
		GatewayPort: (*ports)[2],
	}

}

// EditIPFSField with the native IPFS config command
func EditIPFSField(path, field, value string) {
	cmd := "ipfs -c " + path + " config --json " + field + " " + value
	o, err := exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		fmt.Println(cmd)
		fmt.Println(string(o))
		fmt.Println(err)
	}
}

// SetClusterLeaderConfig set the configs for the leader of a cluster
func (s *Service) SetClusterLeaderConfig(path string,
	replmin, replmax int, ports ClusterInstance) (
	string, string, error) {

	// generate random secret
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		return "", "", errors.New("could not generate secret")
	}
	secret := hex.EncodeToString(key)

	vars := GetClusterVariables(path, s.MyIP, s.Name, secret,
		replmin, replmax, ports)
	return vars, secret, nil
}

// GetClusterVariables get the cluster variables
func GetClusterVariables(path, ip, secret, peername string,
	replmin, replmax int, ports ClusterInstance) string {

	apiIPFSAddr := IPVersion + ip +
		TransportProtocol + strconv.Itoa(ports.IPFSAPIPort) // 5001
	restAPIAddr := IPVersion + ip +
		TransportProtocol + strconv.Itoa(ports.RestAPIPort) // 9094
	IPFSProxyAddr := IPVersion + ip +
		TransportProtocol + strconv.Itoa(ports.IPFSProxyPort) // 9095
	clusterAddr := IPVersion + ip +
		TransportProtocol + strconv.Itoa(ports.ClusterPort) // 9096

	cmd := ""

	// edit peername
	cmd += GetEnvVar("CLUSTER_PEERNAME", peername)

	// edit the secret
	cmd += GetEnvVar("CLUSTER_SECRET", secret)

	// replace IPFS API port (5001)
	cmd += GetEnvVar("CLUSTER_IPFSPROXY_NODEMULTIADDRESS", apiIPFSAddr) // 5001
	cmd += GetEnvVar("CLUSTER_IPFSHTTP_NODEMULTIADDRESS", apiIPFSAddr)  // 5001

	// replace Cluster ports (9094, 9095, 9096)
	cmd += GetEnvVar("CLUSTER_RESTAPI_HTTPLISTENMULTIADDRESS", restAPIAddr) // 9094
	cmd += GetEnvVar("CLUSTER_IPFSPROXY_LISTENMULTIADDRESS", IPFSProxyAddr) // 9095
	cmd += GetEnvVar("CLUSTER_LISTENMULTIADDRESS", clusterAddr)             // 9096

	// replace replication factor
	cmd += GetEnvVar("CLUSTER_REPLICATIONFACTORMIN", strconv.Itoa(replmin))
	cmd += GetEnvVar("CLUSTER_REPLICATIONFACTORMAX", strconv.Itoa(replmax))

	return cmd
}

// GetEnvVar get the environnment variable for the given field and value
func GetEnvVar(field, value string) string {
	// `CLUSTER_FIELD="value" `
	return field + "=\"" + value + "\" "
}

// MakeJSONElem make a JSON single element
func MakeJSONElem(elem string) string {
	// \"elem\"
	return "\\\"" + elem + "\\\""
}

// MakeJSONArray make a json array from the given elements
func MakeJSONArray(elements []string) string {
	// "[
	str := "\"["
	for _, e := range elements {
		// \"elem\"
		str += "\\\"" + e + "\\\""
	}
	// str + ]"
	return str + "]\""
}

// SetupClusterLeader setup a cluster instance for the ARA leader
func (s *Service) SetupClusterLeader(path string,
	replmin, replmax int) (string, *ClusterInstance, error) {

	err := CreateEmptyDir(path)
	if err != nil {
		return "", nil, err
	}

	// generate random secret
	key := make([]byte, 32)
	_, err = rand.Read(key)
	if err != nil {
		return "", nil, errors.New("could not generate secret")
	}
	secret := hex.EncodeToString(key)

	/*
		// path for config files
		path := configPath + "/cluster_" + secret
		err = CreateEmptyDir(path)
		if err != nil {
			return "", nil, err
		}
	*/

	ints, err := GetNextAvailablePorts(s.MinPort, s.MaxPort, ClusterPortNumber)
	if err != nil {
		return "", nil, err
	}

	// set the ports that the cluster will use
	ports := ClusterInstance{
		IP:            IPVersion + s.MyIP + TransportProtocol,
		IPFSAPIPort:   s.MyIPFS.APIPort,
		RestAPIPort:   (*ints)[0],
		IPFSProxyPort: (*ints)[1],
		ClusterPort:   (*ints)[2],
	}

	// get the environment variables to set cluster configs
	vars := GetClusterVariables(path, s.MyIP, secret, s.Name, replmin, replmax,
		ports)

	// init command to be run
	cmd := vars + "ipfs-cluster-service -c " + path + " init"
	o, err := exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		fmt.Println(cmd)
		fmt.Println(string(o))
		fmt.Println(err)
		return "", nil, err
	}

	// start cluster daemon
	cmd = "ipfs-cluster-service -c " + path + " daemon"
	go func() {
		exec.Command("bash", "-c", cmd).Run()
		fmt.Println(s.Name + " cluster leader crashed")
	}()

	addr := IPVersion + s.MyIP + TransportProtocol + strconv.Itoa(ports.RestAPIPort)
	log.Lvl1("Started ipfs-cluster leader at " + addr)

	// wait for the daemon to be launched
	time.Sleep(ClusterStartupTime)

	return secret, &ports, nil
}

// SetupClusterSlave setup a cluster slave instance
func (s *Service) SetupClusterSlave(path, bootstrap, secret string,
	replmin, replmax int) (*ClusterInstance, error) {

	err := CreateEmptyDir(path)
	if err != nil {
		return nil, err
	}

	/*
		// create the config directory, identified by the secret of the cluster
		path := configPath + "/cluster_" + nodeID + "_" + secret
		err := CreateEmptyDir(path)
		if err != nil {
			return nil, err
		}
	*/
	s.PortMutex.Lock()
	/*
		nBig, err := rand.Int(rand.Reader, big.NewInt(MaxPortNumberPerHost-ClusterPortNumber))
		if err != nil {
			panic(err)
		}
		rand := int(nBig.Int64())
	*/

	ints, err := GetNextAvailablePorts(s.MinPort, s.MaxPort, ClusterPortNumber)
	if err != nil {
		log.Lvl1(err)
		return nil, err
	}

	// set the ports that the cluster will use
	ports := ClusterInstance{
		IP:            IPVersion + s.MyIP + TransportProtocol,
		IPFSAPIPort:   s.MyIPFS.APIPort,
		RestAPIPort:   (*ints)[0],
		IPFSProxyPort: (*ints)[1],
		ClusterPort:   (*ints)[2],
	}

	// get the environment variables to set cluster configs
	vars := GetClusterVariables(path, s.MyIP, secret, s.Name,
		replmin, replmax, ports)

	// init command to be run
	cmd := vars + "ipfs-cluster-service -c " + path + " init"
	err = exec.Command("bash", "-c", cmd).Run()
	if err != nil {
		log.Lvl1(err)
		return nil, err
	}

	// start cluster daemon
	cmd = "ipfs-cluster-service -c " + path + " daemon --bootstrap " + bootstrap
	go func() {
		exec.Command("bash", "-c", cmd).Run()
		log.Lvl1("slave " + s.Name + " crashed")
	}()

	// wait for the daemon to be launched
	time.Sleep(ClusterStartupTime)
	s.PortMutex.Unlock()

	addr := ports.IP + strconv.Itoa(ports.RestAPIPort)
	log.Lvl1("Started ipfs-cluster slave at " + addr)

	return &ports, nil
}

/*
// Protocol to start all clusters in an ARA
func Protocol(configPath, nodeID, ip string, replmin, replmax int) error {

	// for all ARAs (trees ?) where nodeID is the leader: do

	// setup the leader of the cluster
	secret, p, err := SetupClusterLeader(configPath, "master", ip,
		replmin, replmax)
	bootstrap := p.IP + strconv.Itoa(p.ClusterPort)
	if err != nil {
		return err
	}
	// for all nodes in this ARA
	_, err = SetupClusterSlave(configPath, "slave1", ip, bootstrap, secret,
		replmin, replmax)
	if err != nil {
		return err
	}
	_, err = SetupClusterSlave(configPath, "slave2", ip, bootstrap, secret,
		replmin, replmax)
	if err != nil {
		return err
	}

	return nil
}
*/
