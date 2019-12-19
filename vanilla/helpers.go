package main

import (
	"path/filepath"
	"strconv"

	"go.dedis.ch/onet/v3/log"
)

// SetNodePaths set the node paths for remote and local node files
func SetNodePaths(n int) {
	NODEPATHREMOTE = NODEPATHNAME + strconv.Itoa(n) + ".txt"
	log.Lvl1("NODEPATHREMOTE:", NODEPATHREMOTE)
	NODEPATHLOCAL = filepath.Join("..", DATAFOLDER, NODEPATHREMOTE)
	log.Lvl1("NODEPATHLOCAL:", NODEPATHLOCAL)
}
