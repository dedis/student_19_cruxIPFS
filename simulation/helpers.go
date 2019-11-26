package main

import (
	"path/filepath"
	"strconv"
)

// SetNodePaths set the node paths for remote and local node files
func SetNodePaths(n int) {
	NODEPATHREMOTE = NODEPATHNAME + strconv.Itoa(n) + ".txt"
	NODEPATHLOCAL = filepath.Join("..", DATAFOLDER, NODEPATHREMOTE)
}
