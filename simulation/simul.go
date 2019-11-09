package main

import (
	// Service needs to be imported here to be instantiated.
	_ "github.com/dedis/student_19_cruxIPFS/service"
	"go.dedis.ch/onet/simul"
)

func main() {
	simul.Start()
}
