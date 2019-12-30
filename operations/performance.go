package operations

import (
	"fmt"
	"io/ioutil"
	"math"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/dedis/student_19_cruxIPFS/service"
	"go.dedis.ch/onet/v3/log"
)

// Test0 perf test takes random nodes (non reproductible), output max and min
// time when performing an operation (write + read)
func Test0(nodesN int) {
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
		str += fmt.Sprintln("min"+id, min.Milliseconds())
		str += fmt.Sprintln("max"+id, max.Milliseconds())
		log.Lvl1(str)
	}
	log.Lvl1("\nDone!")
}

// genSequence generate a random sequence of nOps pairs among nodesN nodes
func genSequence(nOps, nodesN int) []byte {
	log.Lvl1("Generating new sequence")
	str := ""
	for i := 0; i < nOps; i++ {
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

// Test1 try to load an existing sequence, if sequence does not match generate a
// new one. Performance tests, yield min and max operation time (write + read)
func Test1(nOps, nodesN int) {
	ops, err := ioutil.ReadFile(sequenceName)
	if err != nil {
		ops = genSequence(nOps, nodesN)
	}

	lines := strings.Split(string(ops), "\n")
	if len(lines)-1 != nOps {
		ops = genSequence(nOps, nodesN)
		lines = strings.Split(string(ops), "\n")
	}
	for i, l := range lines {
		nodes := strings.Split(l, " ")
		nodeW := nodes[0]
		nodeR := nodes[1]
		f := "file" + strconv.Itoa(i) + ".txt"

		NewFile(f)
		cid, res0 := Write(nodeW, f)
		res1 := Read(nodeR, cid)

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
		str += fmt.Sprintln("min"+id, min.Milliseconds())
		str += fmt.Sprintln("max"+id, max.Milliseconds())
		log.Lvl1(str)
	}
	log.Lvl1("\nDone!")
}

func Test2(nOps, nodesN int) {
	ops, err := ioutil.ReadFile(sequenceName)
	if err != nil {
		log.Lvl1("Error in reading operation sequence file")
		ops = genSequence(nOps, nodesN)
	}

	lines := strings.Split(string(ops), "\n")
	if len(lines)-1 != nOps {
		log.Lvl1("Operation number do not match target")
		log.Lvl1("Got:", len(lines)-1, "Target:", nOps)
		ops = genSequence(nOps, nodesN)
		lines = strings.Split(string(ops), "\n")
		log.Lvl1("After generation")
		log.Lvl1("Got:", len(lines)-1, "Target:", nOps)
	}
	for i, l := range lines {
		nodes := strings.Split(l, " ")
		nodeW := nodes[0]
		nodeR := nodes[1]
		f := "file" + strconv.Itoa(i) + ".txt"

		NewFile(f)
		cid, res0 := Write(nodeW, f)
		res1 := Read(nodeR, cid)

		min := time.Duration(math.MaxInt64)
		max := time.Duration(0)
		str := "\n"
		strread := ""
		strwrite := ""
		for cluster, t0 := range res0 {
			if t1, ok := res1[cluster]; ok {
				strread += fmt.Sprintf("%d ", t0.Milliseconds())
				strwrite += fmt.Sprintf("%d ", t1.Milliseconds())

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
		str += fmt.Sprintln("min"+id, min.Milliseconds())
		str += fmt.Sprintln("max"+id, max.Milliseconds())
		str += strread + "\n"
		str += strwrite + "\n"

		log.Lvl1(str)
	}
	log.Lvl1("\nDone!")
}
