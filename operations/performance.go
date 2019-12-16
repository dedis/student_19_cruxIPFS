package operations

import (
	"math/rand"
	"strconv"

	"github.com/dedis/student_19_cruxIPFS/service"
	"go.dedis.ch/onet/v3/log"
)

var nodesN = 11

// Test0 perf test
func Test0() {
	for i := 0; i < 100; i++ {
		f := "file" + strconv.Itoa(i) + ".txt"
		r0 := rand.Intn(nodesN)
		nodeW := service.NodeName + strconv.Itoa(r0)
		r1 := rand.Intn(nodesN - 1)
		if r1 >= r0 {
			r1++
		}
		nodeR := service.NodeName + strconv.Itoa(r1)
		NewFile(f)
		n, t1 := Write(nodeW, f)
		t2 := Read(nodeR, n)
		log.Lvl1(nodeW, "write time:", t1, nodeR, "read time:", t2)
	}
	log.Lvl1("\nDone!")
}
