package service

import (
	"strconv"
	"strings"

	template "github.com/dedis/student_19_cruxIPFS"
)

// IPFSConfigName name of the config file generated by IPFS
const IPFSConfigName string = "config"

// DefaultTransportProtocol default transport protocol
const DefaultTransportProtocol string = "/tcp/"

// DefaultIPFSSwarmPort default swarm port for ipfs
const DefaultIPFSSwarmPort string = DefaultTransportProtocol + "4001"

// DefaultIPFSAPIPort default API port for ipfs
const DefaultIPFSAPIPort string = DefaultTransportProtocol + "5001"

// DefaultIPFSGatewayPort default gateway port for ipfs
const DefaultIPFSGatewayPort string = DefaultTransportProtocol + "8080"

// EditIPFSConfig edit the ipfs configuration file
func EditIPFSConfig(ports *template.IPFSPorts, path string) error {
	// load the config
	filepath := path + "/" + IPFSConfigName
	conf, err := ReadConfig(filepath)
	if err != nil {
		return err
	}

	// replace the ports
	conf = strings.ReplaceAll(conf, DefaultIPFSSwarmPort,
		DefaultTransportProtocol+strconv.Itoa(ports.Swarm))
	conf = strings.ReplaceAll(conf, DefaultIPFSAPIPort,
		DefaultTransportProtocol+strconv.Itoa(ports.API))
	conf = strings.ReplaceAll(conf, DefaultIPFSGatewayPort,
		DefaultTransportProtocol+strconv.Itoa(ports.Gateway))

	// write the modified config
	err = WriteConfig(filepath, conf)
	if err != nil {
		return err
	}

	return nil
}
