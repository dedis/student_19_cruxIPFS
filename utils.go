package template

import (
	"strings"

	"go.dedis.ch/onet/v3/network"
)

// ServerIdentityToIPString return the ip address of a server identity
func ServerIdentityToIPString(identity *network.ServerIdentity) string {
	// identity.String() - tcp://127.0.0.8:35091
	s := strings.Split(identity.String(), "/")
	// s[2] - 127.0.0.8:35091
	ip := strings.Split(s[2], ":")
	// ip[0] - 127.0.0.8
	return ip[0]
}
