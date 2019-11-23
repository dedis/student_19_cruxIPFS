package template

/*
This holds the messages used to communicate with the service over the network.
*/

import (
	"github.com/dedis/student_19_cruxIPFS/gentree"
	"go.dedis.ch/onet/v3/network"
)

// We need to register all messages so the network knows how to handle them.
func init() {
	network.RegisterMessages()
}

const (
	// ErrorParse indicates an error while parsing the protobuf-file.
	ErrorParse = iota + 4000
)

// InitRequest packet
type InitRequest struct {
	Nodes                []*gentree.LocalityNode
	ServerIdentityToName map[*network.ServerIdentity]string
}

// InitResponse packet
type InitResponse struct {
}

type ReqPings struct {
	SenderName string
}

type ReplyPings struct {
	Pings      string
	SenderName string
}
