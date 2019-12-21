package operations

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"time"

	cruxIPFS "github.com/dedis/student_19_cruxIPFS"
	"github.com/dedis/student_19_cruxIPFS/service"
	"go.dedis.ch/onet/v3/log"
)

var nodesN = 2

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

		//nodeW := "node_0"
		//nodeR := "node_1"
		NewFile(f)
		n, res0 := Write(nodeW, f)
		res1 := Read(nodeR, n)

		min := time.Duration(math.MaxInt64)
		max := time.Duration(0)
		str := "\n"
		for cluster, t0 := range res0 {
			if t1, ok := res1[cluster]; ok {
				sum := t0 + t1
				//str += fmt.Sprintln("write:", t0, "read:", t1, "total:", sum)
				if sum < min {
					min = sum
				}
				if sum > max {
					max = sum
				}
			}
		}
		id := "optime-" + nodeW + "-" + nodeR
		str += fmt.Sprintln(id, min.Milliseconds())
		//str += fmt.Sprintln("min:", id, min.Milliseconds())
		//str += fmt.Sprintln("max:", id, max.Milliseconds())
		log.Lvl1(str)
	}
	log.Lvl1("\nDone!")
}

func Test1() {
	nodes := LoadClusterInstances(cruxIPFS.SaveFile)
	fmt.Println("Round 1")
	for _, n := range nodes {
		fmt.Println(len(n.Clients))
		ListPeers(n.Clients[0])
	}

	time.Sleep(10 * time.Second)
	fmt.Println("Round 2")
	for _, n := range nodes {
		fmt.Println(len(n.Clients))
		ListPeers(n.Clients[0])
	}

}
