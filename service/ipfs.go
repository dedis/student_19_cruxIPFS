package service

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"time"

	"go.dedis.ch/onet/v3/log"
)

// StartIPFS starts an IPFS instance for the given service
func (s *Service) StartIPFS() {

	// get config path
	// e.g $GOPATH/src/github.com/dedis/student_19_cruxIPFS/simulation/build
	pwd, err := os.Getwd()
	checkErr(err)
	//pwd := "/users/guissou"
	configPath := filepath.Join(pwd, ConfigsFolder)

	// set the port range allocated to s
	s.MinPort = BaseHostPort + s.getNodeID()*MaxPortNumberPerHost
	s.MaxPort = s.MinPort + MaxPortNumberPerHost

	s.ConfigPath = filepath.Join(configPath, s.Name)

	// create own config home folder and ipfs config folder
	s.MyIPFSPath = filepath.Join(s.ConfigPath, IPFSFolder)
	checkErr(CreateEmptyDir(s.MyIPFSPath))

	ou, err := exec.Command("ipfs", "-c"+s.MyIPFSPath, "init").Output()
	if err != nil {
		log.Lvl1(string(ou))
		log.Error(err)
	}

	// edit the ip in the config file
	s.EditIPFSConfig()

	// start ipfs daemon
	go func() {
		exec.Command("ipfs", "-c"+s.MyIPFSPath, "daemon").Run()
		fmt.Println("ipfs at ip", s.Name, "crashed")
	}()
	// wait until it has started
	time.Sleep(IPFSStartupTime)
}

// EditIPFSConfig edit the ipfs configuration file (mainly the ip)
func (s *Service) EditIPFSConfig() {
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
		Name:        s.Name,
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

	cmd += GetEnvVar("CLUSTER_REPLICATIONFACTORMIN", "3")
	cmd += GetEnvVar("CLUSTER_REPLICATIONFACTORMAX", "3")

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

	// generate random secret
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		return "", nil, errors.New("could not generate secret")
	}
	secret := hex.EncodeToString(key)

	path = path + "-" + secret
	if CreateEmptyDir(path) != nil {
		return "", nil, err
	}

	s.PortMutex.Lock()
	ints, err := GetNextAvailablePorts(s.MinPort, s.MaxPort, ClusterPortNumber)
	if err != nil {
		return "", nil, err
	}

	// set the ports that the cluster will use
	ports := ClusterInstance{
		HostName:      s.Name,
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
	if ClusterConsensusMode == "crdt" {
		cmd += " --consensus crdt"
	}
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
		log.Lvl1(s.Name + " cluster leader crashed")
	}()

	// wait for the daemon to be launched
	time.Sleep(ClusterStartupTime)
	s.PortMutex.Unlock()

	addr := IPVersion + s.MyIP + TransportProtocol + strconv.Itoa(ports.RestAPIPort)
	log.Lvl2("Started ipfs-cluster leader at " + addr)

	return secret, &ports, nil
}

// SetupClusterSlave setup a cluster slave instance
func (s *Service) SetupClusterSlave(path, bootstrap, secret string,
	replmin, replmax int) (*ClusterInstance, error) {

	err := CreateEmptyDir(path)
	if err != nil {
		return nil, err
	}

	s.PortMutex.Lock()

	ints, err := GetNextAvailablePorts(s.MinPort, s.MaxPort, ClusterPortNumber)
	if err != nil {
		log.Lvl1(err)
		return nil, err
	}

	// set the ports that the cluster will use
	ports := ClusterInstance{
		HostName:      s.Name,
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
	if ClusterConsensusMode == "crdt" {
		cmd += " --consensus crdt"
	}
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
	//log.Lvl1("ipfs-cluster-service -c " + path + " daemon --bootstrap " + bootstrap)
	log.Lvl2("Started ipfs-cluster slave at " + addr)

	return &ports, nil
}
