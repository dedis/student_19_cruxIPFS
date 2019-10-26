package template

/*
This holds the messages used to communicate with the service over the network.
*/

import (
	"go.dedis.ch/onet/v3"
	"go.dedis.ch/onet/v3/network"
)

// We need to register all messages so the network knows how to handle them.
func init() {
	network.RegisterMessages(
		Count{}, CountReply{},
		Clock{}, ClockReply{},
		GenSecret{}, GenSecretReply{},
		StartIPFS{}, StartIPFSReply{},
	)
}

const (
	// ErrorParse indicates an error while parsing the protobuf-file.
	ErrorParse = iota + 4000
)

// Clock will run the tepmlate-protocol on the roster and return
// the time spent doing so.
type Clock struct {
	Roster *onet.Roster
}

// ClockReply returns the time spent for the protocol-run.
type ClockReply struct {
	Time     float64
	Children int
}

// Count will return how many times the protocol has been run.
type Count struct {
}

// CountReply returns the number of protocol-runs
type CountReply struct {
	Count int
}

// GenSecret will return a new shared secret
type GenSecret struct {
}

// GenSecretReply reply of gensecret containing the secret as a string
type GenSecretReply struct {
	Secret string
}

// StartIPFS ipfs start packet
type StartIPFS struct {
	ConfigPath string
	NodeID     string
	PortMin    int
	PortMax    int
}

// StartIPFSReply ss
type StartIPFSReply struct {
	Ports []int
}
