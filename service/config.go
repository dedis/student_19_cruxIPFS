package service

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"os/exec"
	"strconv"
	"time"

	template "github.com/dedis/student_19_cruxIPFS"
)

// IPVersion default ip version
const IPVersion string = "/ip4/"

// TransportProtocol default transport protocol
const TransportProtocol string = "/tcp/"

// EditIPFSConfig edit the ipfs configuration file (mainly the ip)
func EditIPFSConfig(path, ip string) {
	addr := IPVersion + ip + TransportProtocol

	// /ip4/127.0.0.1/tcp/5001
	API := MakeJSONElem(addr + strconv.Itoa(template.DefaultIPFSAPIPort))
	// /ip4/127.0.0.1/tcp/8080
	Gateway := MakeJSONElem(addr +
		strconv.Itoa(template.DefaultIPFSGatewayPort))

	// [\"/ip4/0.0.0.0/tcp/5001\", \"/ip6/::/tcp/5001\"]
	//swarmList := []string{addr + SwarmPort}
	Swarm := MakeJSONArray([]string{addr +
		strconv.Itoa(template.DefaultIPFSSwarmPort)})
	EditIPFSField(path, "Addresses.API", API)
	EditIPFSField(path, "Addresses.Gateway", Gateway)
	EditIPFSField(path, "Addresses.Swarm", Swarm)
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
func SetClusterLeaderConfig(path, ip, peername string,
	replmin, replmax int, ports *template.ClusterPorts) (
	string, string, error) {

	// generate random secret
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		return "", "", errors.New("could not generate secret")
	}
	secret := hex.EncodeToString(key)

	vars := GetClusterVariables(path, ip, peername, secret,
		replmin, replmax, ports)
	return vars, secret, nil
}

// GetClusterVariables get the cluster variables
func GetClusterVariables(path, ip, secret, peername string,
	replmin, replmax int, ports *template.ClusterPorts) string {

	apiIPFSAddr := IPVersion + ip +
		TransportProtocol + strconv.Itoa(ports.IPFSAPI) // 5001
	restAPIAddr := IPVersion + ip +
		TransportProtocol + strconv.Itoa(ports.RestAPI) // 9094
	IPFSProxyAddr := IPVersion + ip +
		TransportProtocol + strconv.Itoa(ports.IPFSProxy) // 9095
	clusterAddr := IPVersion + ip +
		TransportProtocol + strconv.Itoa(ports.Cluster) // 9096

	cmd := ""

	// edit peername
	//cmd += GetEnvVar("CLUSTER_PEERNAME", peername)

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
func SetupClusterLeader(configPath, nodeID, ip string,
	replmin, replmax int) (string, string, error) {

	// generate random secret
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		return "", "", errors.New("could not generate secret")
	}
	secret := hex.EncodeToString(key)

	// path for config files
	path := configPath + "/cluster_" + secret
	err = CreateEmptyDir(path)
	if err != nil {
		return "", "", err
	}

	ints, err := GetNextAvailablePorts(14000, 15000, 3)
	if err != nil {
		return "", "", err
	}

	// set the ports that the cluster will use
	ports := template.ClusterPorts{
		IPFSAPI:   5001,
		RestAPI:   (*ints)[0],
		IPFSProxy: (*ints)[1],
		Cluster:   (*ints)[2],
	}

	bootstrap := IPVersion + ip + TransportProtocol + strconv.Itoa((*ints)[2])

	// get the environment variables to set cluster configs
	vars := GetClusterVariables(path, ip, secret, nodeID,
		replmin, replmax, &ports)

	// init command to be run
	cmd := vars + "ipfs-cluster-service -c " + path + " init"
	o, err := exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		fmt.Println(string(o))
		fmt.Println(err)
		return "", "", err
	}

	// start cluster daemon
	cmd = "ipfs-cluster-service -c " + path + " daemon"
	go func() {
		exec.Command("bash", "-c", cmd).Run()
		fmt.Println(ip + " cluster crashed")
	}()

	// wait for the daemon to be launched
	time.Sleep(2 * time.Second)

	return secret, bootstrap, nil
}

// SetupClusterSlave setup a cluster slave instance
func SetupClusterSlave(configPath, nodeID, ip, bootstrap, secret string,
	replmin, replmax int) error {

	// create the config directory, identified by the secret of the cluster
	path := configPath + "/cluster_" + nodeID + "_" + secret
	err := CreateEmptyDir(path)
	if err != nil {
		return err
	}

	ints, err := GetNextAvailablePorts(14000, 15000, 3)
	if err != nil {
		return err
	}

	// set the ports that the cluster will use
	ports := template.ClusterPorts{
		IPFSAPI:   template.DefaultIPFSAPIPort,
		RestAPI:   (*ints)[0],
		IPFSProxy: (*ints)[1],
		Cluster:   (*ints)[2],
	}

	// get the environment variables to set cluster configs
	vars := GetClusterVariables(path, ip, secret, nodeID,
		replmin, replmax, &ports)

	// init command to be run
	cmd := vars + "ipfs-cluster-service -c " + path + " init"
	err = exec.Command("bash", "-c", cmd).Run()
	if err != nil {
		return err
	}

	// start cluster daemon
	cmd = "ipfs-cluster-service -c " + path + " daemon --bootstrap " + bootstrap
	go exec.Command("bash", "-c", cmd).Run()

	// wait for the daemon to be launched
	time.Sleep(2 * time.Second)

	return nil
}

// Protocol to start all clusters in an ARA
func Protocol(configPath, nodeID, ip string, replmin, replmax int) error {

	// for all ARAs (trees ?) where nodeID is the leader: do

	// setup the leader of the cluster
	secret, bootstrap, err := SetupClusterLeader(configPath, "master", ip,
		replmin, replmax)
	if err != nil {
		return err
	}
	// for all nodes in this ARA
	err = SetupClusterSlave(configPath, "slave1", ip, bootstrap, secret,
		replmin, replmax)
	if err != nil {
		return err
	}
	err = SetupClusterSlave(configPath, "slave2", ip, bootstrap, secret,
		replmin, replmax)
	if err != nil {
		return err
	}

	return nil
}
