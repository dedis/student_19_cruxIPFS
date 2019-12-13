/*
Generates an ns file for a random graph with OSPF routing.

Takes as input the nr of nodes N and size of x,y coordinate space and sets up a graph with N routers and N nodes. generates random cordinates for each router and sets
links between routers according to their coordinates. The use of coordinates avoids generating routes with TIVs
(triangle inequality violations).

Outputs
- the ns file delay.ns
- the node ips and levels in out.txt
- the coordinates of each router in -N.txt

*/

package main

import (
	"bufio"
	"flag"
	"fmt"
	"math"
	"math/rand"
	"os"
	"strconv"
	"time"

	"go.dedis.ch/onet/v3/log"
)

type DeterNode struct {
	Name   string
	Level  int
	Links  map[string]bool
	Dist   map[string]float64
	IP     []string
	HostIP string
	X      float64
	Y      float64
}

//type myNodes []DeterNode

var MAX_IFACES = 4

/*
func (s myNodes) Len() int {
	return len(s)
}
func (s myNodes) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s myNodes) Less(i, j int) bool {
	return math.Pow(s[i].X, 2.0) + math.Pow(s[i].Y, 2.0) < math.Pow(s[j].X, 2.0) + math.Pow(s[j].Y, 2.0)
}
*/

func dist(a DeterNode, b DeterNode) float64 {
	return math.Sqrt(math.Pow(float64(a.X-b.X), 2.0) + math.Pow(float64(a.Y-b.Y), 2.0))
}

func genAndPrintRndRouters(N int, R int, SpaceMax int, K int, zeroY bool, trueLatency bool) {

	var RndSrc *rand.Rand
	RndSrc = rand.New(rand.NewSource(time.Now().UnixNano()))
	routers := make([]DeterNode, R)
	nodes := make([]DeterNode, N)

	for i := 0; i < R; i++ {
		routers[i].Name = "router_" + strconv.Itoa(i)
		routers[i].Links = make(map[string]bool)
		routers[i].Dist = make(map[string]float64)
		routers[i].IP = make([]string, 0)
		routers[i].Level = 0
		routers[i].X = RndSrc.Float64() * float64(SpaceMax)
		routers[i].Y = RndSrc.Float64() * float64(2)

	}

	for i := 0; i < N; i++ {
		nodes[i].Name = "node_" + strconv.Itoa(i)
		nodes[i].Links = make(map[string]bool)
		nodes[i].Dist = make(map[string]float64)
		nodes[i].IP = make([]string, 0)
		nodes[i].Level = 0
	}

	prob := 1.0 / math.Pow(float64(N), 1.0/float64(K))
	for lvl := 0; lvl < K; lvl++ {
		for i := 0; i < N; i++ {
			if nodes[i].Level == lvl-1 {
				rnd := RndSrc.Float64()
				if rnd < prob {
					nodes[i].Level = lvl
				}
			}
		}
	}

	file2, _ := os.Create("delay.ns")
	defer file2.Close()
	w2 := bufio.NewWriter(file2)

	w2.WriteString("set ns [new Simulator]\n")
	w2.WriteString("source tb_compat.tcl\n")
	w2.WriteString("\n")
	w2.WriteString("tb-use-endnodeshaping 1\n")
	w2.WriteString("set n_nodes " + strconv.Itoa(N) + "\n")
	w2.WriteString("set n_routers " + strconv.Itoa(R) + "\n")
	w2.WriteString("\n")

	// nodes
	w2.WriteString("for {set i 0} {$i < $n_nodes} {incr i} {\n")

	w2.WriteString("\tset site($i) [$ns node]\n")
	w2.WriteString("\ttb-set-hardware $site($i) {MicroCloud}\n")
	w2.WriteString("\ttb-set-node-os $site($i) Ubuntu1404-64-STD\n")

	w2.WriteString("}\n")

	// routers
	w2.WriteString("for {set i 0} {$i < $n_routers} {incr i} {\n")

	w2.WriteString("\tset router($i) [$ns node]\n")
	w2.WriteString("\ttb-set-hardware $router($i) {MicroCloud}\n")
	w2.WriteString("\ttb-set-node-os $router($i) Ubuntu1404-64-STD\n")

	w2.WriteString("}\n\n")

	// make sure the network is connected

	linkNr := 0

	for i := 0; i < R; i++ {
		attempts := 0
		for {

			log.LLvl1("node", i, "attempts", attempts)
			if len(routers[i].IP) == MAX_IFACES || attempts == 100 {
				break
			}
			// generate a random node to connect to
			peerIdx := RndSrc.Intn(R)
			peerName := routers[peerIdx].Name

			for {
				if (peerIdx != i && !routers[i].Links[peerName] && len(routers[peerIdx].IP) < MAX_IFACES) || attempts == 100 {
					break
				}
				peerIdx = RndSrc.Intn(R)
				peerName = routers[peerIdx].Name
				attempts++
			}

			if attempts == 100 {
				continue
			}

			// connect the two nodes

			idx1 := strconv.Itoa(i)
			idx2 := strconv.Itoa(peerIdx)

			// generate random latency in 2-10 ms
			//min := 0.0

			dist := dist(routers[i], routers[peerIdx])
			//dist := RndSrc.Float64() * 5.0 + 2.0

			// round distance to closest integer
			if dist < 1 {
				dist = 0
			} else if dist < 2 {
				dist = 2
			}

			fmt.Println("generated", dist)

			//l := ""

			/*
				// check if triangle inequality stuff
				for j := 0 ; j < N ; j++ {
					// do they connect through j? then add that as constraint
					nameJ := nodes[j].Name
					if nodes[i].Links[nameJ] && nodes[peerIdx].Links[nameJ] {
						min = nodes[i].Dist[nameJ] + nodes[peerIdx].Dist[nameJ]
						l += nameJ
					}
				}
				if min != 0 {
					log.LLvl1(i, peerIdx, "have min", min, "through node", l)
					dist = min
				}
			*/

			w2.WriteString("set link" + strconv.Itoa(linkNr) + " [$ns duplex-link $router(" + idx1 + ") $router(" + idx2 + ") 100Mb " + fmt.Sprintf("%f", dist) + "ms DropTail]\n")
			w2.WriteString("tb-set-ip-link $router(" + idx1 + ") $link" + strconv.Itoa(linkNr) + " 10.1." + strconv.Itoa(linkNr+1) + ".2\n")
			w2.WriteString("tb-set-ip-link $router(" + idx2 + ") $link" + strconv.Itoa(linkNr) + " 10.1." + strconv.Itoa(linkNr+1) + ".3\n")

			//w2.WriteString("tb-set-ip-link $site(" + idx1 + ") $link" + strconv.Itoa(linkNr) + " 10.1." + strconv.Itoa(i+1) + "." + strconv.Itoa(len(nodes[i].IP) + 2)+ "\n")
			//w2.WriteString("tb-set-ip-link $site(" + idx2 + ") $link" + strconv.Itoa(linkNr) + " 10.1." + strconv.Itoa(peerIdx+1) + "." + strconv.Itoa(len(nodes[peerIdx].IP) + 2) + "\n")

			routers[i].Links[peerName] = true
			routers[peerIdx].Links[routers[i].Name] = true
			routers[i].Dist[peerName] = dist
			routers[peerIdx].Dist[routers[i].Name] = dist

			routers[i].IP = append(routers[i].IP, "10.1."+strconv.Itoa(linkNr+1)+".2")
			routers[peerIdx].IP = append(routers[peerIdx].IP, "10.1."+strconv.Itoa(linkNr+1)+".3")

			//nodes[i].IP = append(nodes[i].IP, "10.1." + strconv.Itoa(i+1) + "." + strconv.Itoa(len(nodes[i].IP) + 2))
			//nodes[peerIdx].IP = append(nodes[peerIdx].IP, "10.1." + strconv.Itoa(peerIdx+1) + "." + strconv.Itoa(len(nodes[peerIdx].IP) + 2))

			log.LLvl1(i, peerIdx, dist)

			linkNr++

		}
	}

	// connect two nodes to their corresponding router through a link of latency 0, respectively random latencies

	for i := 0; i < N; i++ {
		//routerIdx := strconv.Itoa(i/2)
		routerIdx := strconv.Itoa(i)
		nodeIdx := strconv.Itoa(i)

		latency := 0

		/*
			if i %2 == 1 {
				latency = RndSrc.Intn(1) + 2
			}
		*/

		w2.WriteString("set link" + strconv.Itoa(linkNr) + " [$ns duplex-link $router(" + routerIdx + ") $site(" + nodeIdx + ") 100Mb " + fmt.Sprintf("%d", latency) + "ms DropTail]\n")
		w2.WriteString("tb-set-ip-link $router(" + routerIdx + ") $link" + strconv.Itoa(linkNr) + " 10.1." + strconv.Itoa(linkNr+1) + ".2\n")
		w2.WriteString("tb-set-ip-link $site(" + nodeIdx + ") $link" + strconv.Itoa(linkNr) + " 10.1." + strconv.Itoa(linkNr+1) + ".3\n")

		nodes[i].HostIP = "10.1." + strconv.Itoa(linkNr+1) + ".3"

		linkNr++
	}

	w2.WriteString("\n")
	// OSPF routing
	//w2.WriteString("$ns rtproto Session\n")
	w2.WriteString("$ns rtproto Manual\n")

	w2.WriteString("\n")
	w2.WriteString("$ns run")

	w2.Flush()

	file, _ := os.Create("out.txt")
	defer file.Close()
	w := bufio.NewWriter(file)

	// print nodes in the out experiment file
	for i := 0; i < N; i++ {
		ips := ""
		for _, ip := range nodes[i].IP {
			ips += ip + ","
		}

		//w.WriteString(nodes[i].Name + " " + ips + " " + strconv.Itoa(nodes[i].Level) + "\n")
		w.WriteString(nodes[i].Name + " " + nodes[i].HostIP + " " + strconv.Itoa(nodes[i].Level) + "\n")
	}

	w.Flush()

	file3, _ := os.Create("coords.txt")
	defer file3.Close()
	w3 := bufio.NewWriter(file3)

	for i := 0; i < N; i++ {
		w3.WriteString(nodes[i].Name + " " + fmt.Sprintf("%f", routers[i].X) + " " + fmt.Sprintf("%f", routers[i].Y) + "\n")
	}

	w3.Flush()

}

func main() {

	K := flag.Int("K", 3, "Number of levels.")
	N := flag.Int("N", 10, "Number of nodes.")
	// TODO S doesn't have any function at the moment
	R := flag.Int("R", 10, "Number of routers.")
	SpaceMax := flag.Int("SpaceMax", 15, "Coordinate space size.")

	flag.Parse()

	genAndPrintRndRouters(*N, *R, *SpaceMax, *K, true, true)

}
