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
}

// Test2 measure nOps write + read operation between pairs of nodes among a set
// of nodesN nodes. It first tries to read the operation sequence (to reproduce)
// the same sequence as a previous experiment, if does not exist or invalid
// format, generate a new sequence
func Test2(nOps, nodesN int) {
	log.Lvl1("Starting Test2")
	ops, err := ioutil.ReadFile(sequenceName)
	if err != nil {
		log.Lvl1("Error in reading operation sequence file")
		ops = genSequence(nOps, nodesN)
	}

	lines := strings.Split(string(ops), "\n")
	if len(lines)-1 != nOps {
		ops = genSequence(nOps, nodesN)
		lines = strings.Split(string(ops), "\n")
	} else {
		for _, l := range lines[:len(lines)-1] {
			k, err := strconv.Atoi(
				strings.Split(l, " ")[0][len(service.NodeName):])
			checkErr(err)
			if k >= nodesN {
				ops = genSequence(nOps, nodesN)
				lines = strings.Split(string(ops), "\n")
				break
			}
		}
	}
	for i, l := range lines[:len(lines)-1] {
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
}
