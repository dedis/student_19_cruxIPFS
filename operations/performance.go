package operations

import (
	"fmt"
	"io/ioutil"
	"math"
	"math/rand"
	"strconv"
	"strings"
	"time"

	cruxIPFS "github.com/dedis/student_19_cruxIPFS"
	"github.com/dedis/student_19_cruxIPFS/service"
	"go.dedis.ch/onet/v3/log"
)

var nodesN = 13

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
		//str += fmt.Sprintln(id, min.Milliseconds())
		str += fmt.Sprintln("min"+id, min.Milliseconds())
		str += fmt.Sprintln("max"+id, max.Milliseconds())
		log.Lvl1(str)
	}
	log.Lvl1("\nDone!")
}

func genSequence(n int) []byte {
	log.Lvl1("Generating new sequence")
	str := ""
	for i := 0; i < n; i++ {
		r0 := rand.Intn(nodesN)
		nodeW := service.NodeName + strconv.Itoa(r0)
		r1 := rand.Intn(nodesN - 1)
		if r1 >= r0 {
			r1++
		}
		nodeR := service.NodeName + strconv.Itoa(r1)
		str += nodeW + " " + nodeR + "\n"
	}
	ioutil.WriteFile(sequenceName, []byte(str), 0777)
	return []byte(str)
}

func Test1(n int) {
	ops, err := ioutil.ReadFile(sequenceName)
	if err != nil {
		ops = genSequence(n)
	}

	lines := strings.Split(string(ops), "\n")
	if len(lines)-1 != n {
		ops = genSequence(n)
		lines = strings.Split(string(ops), "\n")
	}
	for i, l := range lines {
		nodes := strings.Split(l, " ")
		nodeW := nodes[0]
		nodeR := nodes[1]
		f := "file" + strconv.Itoa(i) + ".txt"

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
		//str += fmt.Sprintln(id, min.Milliseconds())
		str += fmt.Sprintln("min"+id, min.Milliseconds())
		str += fmt.Sprintln("max"+id, max.Milliseconds())
		log.Lvl1(str)
	}
	log.Lvl1("\nDone!")
}

func Test2() {
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
