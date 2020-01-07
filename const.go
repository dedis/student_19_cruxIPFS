package cruxIPFS

import (
	"go.dedis.ch/onet/v3/network"
)

const (
	// SaveFile file where instances are saved
	SaveFile = "../save.txt"

	// DataFolder name of the folder containing data
	DataFolder = "data"
	// ScriptsFolder name of the folder containing the scripts
	ScriptsFolder = "scripts"

	// ErrorParse indicates an error while parsing the protobuf-file.
	ErrorParse = iota + 4000
)

// We need to register all messages so the network knows how to handle them.
func init() {
	network.RegisterMessages()
}
