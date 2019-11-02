package template

import (
	"strings"

	"go.dedis.ch/onet/v3/network"
)

// ServerIdentityToIPString return the ip address of a server identity
func ServerIdentityToIPString(identity *network.ServerIdentity) string {
	s := strings.Split(identity.String(), "/")
	ip := strings.Split(s[2], ":")
	return ip[0]
}
