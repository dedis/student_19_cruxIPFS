package main

import (
	"fmt"
	"os"
	"os/exec"
)

func main() {

	os.RemoveAll("/home/guillaume/.ipfs")
	err := exec.Command("bash", "-c", "ipfs init").Run()
	if err != nil {
		fmt.Println(err)
	}

	err = exec.Command("bash", "-c", "ipfs config Addresses.API \"/ip4/127.0.0.2/tcp/5001\"").Run()
	if err != nil {
		fmt.Println(err)
	}

	o, err := exec.Command("bash", "-c", "ipfs config --json Addresses.Swarm \"[\\\"/ip6/::/tcp/12345\\\"]\"").Output()
	fmt.Println(string(o))
	if err != nil {
		fmt.Println(err)
	}

	TransportProtocol := "/tcp/"
	IPVersion := "/ip4/"

	gateway := "9090"
	swarm := "4009"
	api := "5009"
	ip := "127.0.0.9"

	path := "/home/guillaume/.ipfs"
	addr := IPVersion + ip + TransportProtocol

	// \"/ip4/127.0.0.1/tcp/5001\"
	API := MakeJSONElem(addr + api)
	// \"/ip4/127.0.0.1/tcp/8080\"
	Gateway := MakeJSONElem(addr + gateway)

	// [\"/ip4/0.0.0.0/tcp/5001\", \"/ip6/::/tcp/5001\"]
	swarmList := []string{addr + swarm}
	Swarm := MakeJSONArray(swarmList)
	EditIPFSField(path, "Addresses.API", API)
	EditIPFSField(path, "Addresses.Gateway", Gateway)
	EditIPFSField(path, "Addresses.Swarm", Swarm)

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
