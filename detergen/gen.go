package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

var lastPart = "$ns rtproto Manual\n\n$ns run"

func readNodes() (map[string]map[string]float64, int, int,
	map[string]map[string]bool, map[string]map[string]string) {

	nodes := make(map[string]map[string]float64)
	directions := make(map[string]map[string]bool)

	readLine, _ := ReadFileLineByLine("deter.ns")
	nrNodes := 0
	nrRouters := 0

	ipLinks := make(map[string]map[string]string)

	for true {
		line := readLine()
		//fmt.Println(line)
		if line == "" {
			//fmt.Println("end")
			continue
		}

		if strings.HasPrefix(line, "$ns run") {
			break
		}

		if strings.HasPrefix(line, "#") {
			continue
		}

		if strings.HasPrefix(line, "set n_nodes") {
			nrNodes, _ = strconv.Atoi(strings.Split(line, " ")[2])

			for i := 0; i < nrNodes; i++ {
				name := "site_" + strconv.Itoa(i)
				nodes[name] = make(map[string]float64)
				directions[name] = make(map[string]bool)
				ipLinks[name] = make(map[string]string)
			}
		}

		if strings.HasPrefix(line, "set n_nodes") {
			nrNodes, _ = strconv.Atoi(strings.Split(line, " ")[2])

			for i := 0; i < nrNodes; i++ {
				name := "site_" + strconv.Itoa(i)
				nodes[name] = make(map[string]float64)
				directions[name] = make(map[string]bool)
				ipLinks[name] = make(map[string]string)
			}
		}

		if strings.HasPrefix(line, "set n_routers") {
			nrRouters, _ = strconv.Atoi(strings.Split(line, " ")[2])

			for i := 0; i < nrRouters; i++ {

				name := "router_" + strconv.Itoa(i)
				nodes[name] = make(map[string]float64)
				directions[name] = make(map[string]bool)
				ipLinks[name] = make(map[string]string)

				for j := 0; j < nrNodes; j++ {
					//namej := "node_" + strconv.Itoa(j)
					//nodes[name][namej] = math.MaxFloat64
				}
			}

			//fmt.Println(nrNodes)
		}

		if strings.HasPrefix(line, "set link") {
			words := strings.Split(line, " ")
			node1, node2 := words[4], words[5]
			// $site(0)
			node1idx, _ := strconv.Atoi(strings.Split(strings.Split(node1,
				"(")[1], ")")[0])
			name1 := strings.Split(strings.Split(node1, "(")[0], "$")[1] +
				"_" + strconv.Itoa(node1idx)

			// $site(1)
			//fmt.Println(node1, node2)
			node2idx, _ := strconv.Atoi(strings.Split(strings.Split(node2,
				"(")[1], ")")[0])
			name2 := strings.Split(strings.Split(node2, "(")[0], "$")[1] +
				"_" + strconv.Itoa(node2idx)

			//print(name1, name2)

			directions[name1][name2] = true
			directions[name2][name1] = true

			distStr := words[7]
			//fmt.Println(distStr, distStr[:len(distStr)-2])
			dist, err := strconv.ParseFloat(distStr[:len(distStr)-2], 64)
			if err != nil {
				fmt.Println("Problem when parsing pings", distStr,
					distStr[:len(distStr)-2])
			}

			nodes[name1][name2] = dist
			nodes[name2][name1] = dist

			lineNode1 := readLine()
			words2 := strings.Split(lineNode1, " ")
			node1iface, ip := words2[1], words2[3]
			node1ifaceIdx, _ := strconv.Atoi(strings.Split(strings.Split(
				node1iface, "(")[1], ")")[0])
			if node1ifaceIdx == node1idx {
				ipLinks[name1][name2] = ip
			}
			if node1ifaceIdx == node2idx {
				ipLinks[name2][name1] = ip
			}

			lineNode2 := readLine()
			words3 := strings.Split(lineNode2, " ")
			node2iface, ip := words3[1], words3[3]
			node2ifaceIdx, _ := strconv.Atoi(strings.Split(strings.Split(
				node2iface, "(")[1], ")")[0])
			if node2ifaceIdx == node1idx {
				ipLinks[name1][name2] = ip
			}
			if node2ifaceIdx == node2idx {
				ipLinks[name2][name1] = ip
			}

		}

	}

	//fmt.Println("here", nodes)

	return nodes, nrNodes, nrRouters, directions, ipLinks
}

func readPings() map[string]map[string]float64 {

	pings := make(map[string]map[string]float64)

	// read from file lines of fomrm "ping node_19 node_7 = 32.317"
	readLine, _ := ReadFileLineByLine("pings.txt")

	for true {
		line := readLine()
		if line == "" {
			break
		}

		if strings.HasPrefix(line, "#") {
			continue
		}

		tokens := strings.Split(line, " ")
		src := tokens[1]
		dst := tokens[2]
		pingTime, err := strconv.ParseFloat(tokens[4], 64)
		if err != nil {
			fmt.Println("Problem when parsing pings")
		}

		if _, ok := pings[src]; !ok {
			pings[src] = make(map[string]float64)
		}

		pings[src][dst] = math.Round(pingTime*100) / 100
	}

	return pings
}

func floydWarshall(N int, R int, dist map[string]map[string]float64) (
	map[string]map[string]float64, map[string]map[string]string) {

	shortest := make(map[string]map[string]float64)
	next := make(map[string]map[string]string)

	for i := 0; i < N; i++ {
		name := "site_" + strconv.Itoa(i)
		shortest[name] = make(map[string]float64)
		next[name] = make(map[string]string)
	}

	for i := 0; i < R; i++ {
		name := "router_" + strconv.Itoa(i)
		shortest[name] = make(map[string]float64)
		next[name] = make(map[string]string)
	}

	for i := 0; i < N; i++ {
		name := "site_" + strconv.Itoa(i)
		shortest[name] = make(map[string]float64)
		for j := 0; j < N; j++ {
			namej := "site_" + strconv.Itoa(j)
			shortest[name][namej] = math.MaxFloat64
		}
		for j := 0; j < R; j++ {
			namej := "router_" + strconv.Itoa(j)
			shortest[name][namej] = math.MaxFloat64
		}
	}

	for i := 0; i < R; i++ {
		name := "router_" + strconv.Itoa(i)
		shortest[name] = make(map[string]float64)
		for j := 0; j < N; j++ {
			namej := "site_" + strconv.Itoa(j)
			shortest[name][namej] = math.MaxFloat64
		}
		for j := 0; j < R; j++ {
			namej := "router_" + strconv.Itoa(j)
			shortest[name][namej] = math.MaxFloat64
		}
	}

	for x, m := range dist {
		for y, d := range m {
			shortest[x][y] = d
			shortest[x][x] = 0
			next[x][y] = y
			next[x][x] = ""
		}
	}

	for k := 0; k < N; k++ {
		namek := "site_" + strconv.Itoa(k)
		for i := 0; i < N; i++ {
			namei := "site_" + strconv.Itoa(i)
			for j := 0; j < N; j++ {
				namej := "site_" + strconv.Itoa(j)
				if shortest[namei][namej] >
					shortest[namei][namek]+shortest[namek][namej] {

					shortest[namei][namej] =
						shortest[namei][namek] + shortest[namek][namej]
					next[namei][namej] = next[namei][namek]
				}
			}
			for j := 0; j < R; j++ {
				namej := "router_" + strconv.Itoa(j)
				if shortest[namei][namej] >
					shortest[namei][namek]+shortest[namek][namej] {

					shortest[namei][namej] =
						shortest[namei][namek] + shortest[namek][namej]
					next[namei][namej] = next[namei][namek]
				}
			}
		}
		for i := 0; i < R; i++ {
			namei := "router_" + strconv.Itoa(i)
			for j := 0; j < N; j++ {
				namej := "site_" + strconv.Itoa(j)
				if shortest[namei][namej] >
					shortest[namei][namek]+shortest[namek][namej] {

					shortest[namei][namej] =
						shortest[namei][namek] + shortest[namek][namej]
					next[namei][namej] = next[namei][namek]
				}
			}
			for j := 0; j < R; j++ {
				namej := "router_" + strconv.Itoa(j)
				if shortest[namei][namej] >
					shortest[namei][namek]+shortest[namek][namej] {

					shortest[namei][namej] =
						shortest[namei][namek] + shortest[namek][namej]
					next[namei][namej] = next[namei][namek]
				}
			}
		}
	}

	for k := 0; k < R; k++ {
		namek := "router_" + strconv.Itoa(k)
		for i := 0; i < N; i++ {
			namei := "site_" + strconv.Itoa(i)
			for j := 0; j < N; j++ {
				namej := "site_" + strconv.Itoa(j)
				if shortest[namei][namej] >
					shortest[namei][namek]+shortest[namek][namej] {

					shortest[namei][namej] =
						shortest[namei][namek] + shortest[namek][namej]
					next[namei][namej] = next[namei][namek]
				}
			}
			for j := 0; j < R; j++ {
				namej := "router_" + strconv.Itoa(j)
				if shortest[namei][namej] >
					shortest[namei][namek]+shortest[namek][namej] {

					shortest[namei][namej] =
						shortest[namei][namek] + shortest[namek][namej]
					next[namei][namej] = next[namei][namek]
				}
			}
		}
		for i := 0; i < R; i++ {
			namei := "router_" + strconv.Itoa(i)
			for j := 0; j < N; j++ {
				namej := "site_" + strconv.Itoa(j)
				if shortest[namei][namej] >
					shortest[namei][namek]+shortest[namek][namej] {

					shortest[namei][namej] =
						shortest[namei][namek] + shortest[namek][namej]
					next[namei][namej] = next[namei][namek]
				}
			}
			for j := 0; j < R; j++ {
				namej := "router_" + strconv.Itoa(j)
				if shortest[namei][namej] >
					shortest[namei][namek]+shortest[namek][namej] {

					shortest[namei][namej] =
						shortest[namei][namek] + shortest[namek][namej]
					next[namei][namej] = next[namei][namek]
				}
			}
		}
	}
	return shortest, next
}

func findPath(u string, v string, next map[string]map[string]string) []string {

	path := make([]string, 0)

	if next[u][v] == "" {
		return path
	}

	path = append(path, u)

	for {
		u = next[u][v]
		path = append(path, u)

		if u == v {
			break
		}
	}

	return path
}

func printShortestNsPaths(N int, R int, shortest map[string]map[string]float64,
	next map[string]map[string]string,
	directions map[string]map[string]bool) map[string]string {

	deternsBytes, err := ioutil.ReadFile("deter.ns")
	if err != nil {
		panic(err)
	}
	deterns := string(deternsBytes)
	os.Remove("deter.ns")
	startDeterns := strings.Split(deterns, lastPart)[0]

	file8, err := os.Create("../data/gen/deter.ns")
	if err != nil {
		fmt.Println(err)
	}
	w8 := bufio.NewWriter(file8)

	w8.WriteString(startDeterns)

	firstLink := make(map[string]string)

	isFirstFor := make(map[string]map[string]bool)

	existing := make(map[string]map[string]map[string]string)
	for i := 0; i < N; i++ {
		namei := "site_" + strconv.Itoa(i)
		existing[namei] = make(map[string]map[string]string)
		firstLink[namei] = ""
		isFirstFor[namei] = make(map[string]bool)
		for j := 0; j < N; j++ {
			namej := "site_" + strconv.Itoa(j)
			existing[namei][namej] = make(map[string]string)
		}
		for j := 0; j < R; j++ {
			namej := "router_" + strconv.Itoa(j)
			existing[namei][namej] = make(map[string]string)
		}
	}

	for i := 0; i < R; i++ {
		namei := "router_" + strconv.Itoa(i)
		existing[namei] = make(map[string]map[string]string)
		firstLink[namei] = ""
		isFirstFor[namei] = make(map[string]bool)
		for j := 0; j < N; j++ {
			namej := "site_" + strconv.Itoa(j)
			existing[namei][namej] = make(map[string]string)
		}
		for j := 0; j < R; j++ {
			namej := "router_" + strconv.Itoa(j)
			existing[namei][namej] = make(map[string]string)
		}
	}

	// router router and router site
	for i := 0; i < R; i++ {
		namei := "router_" + strconv.Itoa(i)
		for j := 0; j < R; j++ {

			if i == j {
				continue
			}

			namej := "router_" + strconv.Itoa(j)

			if shortest[namei][namej] > 10000 {
				panic("no route!")
			}

			path := findPath(namei, namej, next)

			//fmt.Println(namei, namej, path, shortest[namei][namej])

			nameNextHop := path[1]
			nameLastHop := path[len(path)-2]
			idxNextHop := strings.Split(nameNextHop, "_")[1]
			idxLastHop := strings.Split(nameLastHop, "_")[1]

			str := ""

			// this is the shortest path to that point

			// if the path there is not defined yet, define it now
			if len(path) > 2 {

				if existing[namei][namej][nameLastHop] == "" {
					str = "$router(" + strconv.Itoa(i) +
						") add-route [$ns link $router(" + strconv.Itoa(j) +
						") $router(" + idxLastHop + ")] $router(" + idxNextHop +
						")\n"
					existing[namei][namej][nameLastHop] = nameNextHop
					existing[namei][nameLastHop][namej] = nameNextHop
					w8.WriteString(str)
				}
			}

			if existing["site_"+strconv.Itoa(i)]["router_"+strconv.
				Itoa(j)][nameLastHop] == "" {

				str = "$site(" + strconv.Itoa(i) +
					") add-route [$ns link $router(" + strconv.Itoa(j) +
					") $router(" + idxLastHop + ")] $router(" +
					strconv.Itoa(i) + ")\n"
				existing["site_"+strconv.Itoa(i)]["router_"+strconv.
					Itoa(j)][nameLastHop] = "router_" + strconv.Itoa(i)
				existing["site_"+strconv.Itoa(i)][nameLastHop]["router_"+
					strconv.Itoa(j)] = "router_" + strconv.Itoa(i)
				w8.WriteString(str)
			}

			if existing["router_"+strconv.Itoa(i)]["site_"+
				strconv.Itoa(j)]["router_"+strconv.Itoa(j)] == "" {

				str = "$router(" + strconv.Itoa(i) +
					") add-route [$ns link $site(" + strconv.Itoa(j) +
					") $router(" + strconv.Itoa(j) + ")] $router(" +
					idxNextHop + ")\n"
				existing["router_"+strconv.Itoa(i)]["site_"+strconv.
					Itoa(j)]["router_"+strconv.Itoa(j)] = nameNextHop
				existing["router_"+strconv.Itoa(i)]["router_"+strconv.
					Itoa(j)]["site_"+strconv.Itoa(j)] = nameNextHop
				w8.WriteString(str)
			}

			if existing["site_"+strconv.Itoa(i)]["site_"+
				strconv.Itoa(j)]["router_"+strconv.Itoa(j)] == "" {

				str = "$site(" + strconv.Itoa(i) +
					") add-route [$ns link $site(" + strconv.Itoa(j) +
					") $router(" + strconv.Itoa(j) + ")] $router(" +
					strconv.Itoa(i) + ")\n"
				existing["site_"+strconv.Itoa(i)]["site_"+strconv.
					Itoa(j)]["router_"+strconv.Itoa(j)] = "router_" +
					strconv.Itoa(i)
				existing["site_"+strconv.Itoa(i)]["router_"+strconv.
					Itoa(j)]["site_"+strconv.Itoa(j)] = "router_" +
					strconv.Itoa(i)
				w8.WriteString(str)
			}

			// path to main interface has to be shortest

			// is the main interface defined?
			if firstLink[namej] == "" {
				firstLink[namej] = nameLastHop
			}

			// look at all outgoing links of namej
			for outgoingNode, exists := range directions[namej] {
				if exists && outgoingNode != nameLastHop {
					idxOutgoingNode := strings.Split(outgoingNode, "_")[1]
					if firstLink[namej] != outgoingNode &&
						existing[namei][namej][outgoingNode] == "" {

						if strconv.Itoa(j) != idxOutgoingNode {
							str = "$router(" + strconv.Itoa(i) +
								") add-route [$ns link $router(" +
								strconv.Itoa(j) + ") $router(" +
								idxOutgoingNode + ")] $router(" + idxNextHop +
								")\n"
							existing[namei][namej][outgoingNode] = nameNextHop
							existing[namei][outgoingNode][namej] = nameNextHop
							w8.WriteString(str)

							if existing["site_"+strconv.Itoa(i)]["router_"+
								strconv.Itoa(j)][outgoingNode] == "" {

								str = "$site(" + strconv.Itoa(i) +
									") add-route [$ns link $router(" + strconv.
									Itoa(j) + ") $router(" + idxOutgoingNode +
									")] $router(" + strconv.Itoa(i) + ")\n"
								existing["site_"+strconv.Itoa(i)]["router_"+
									strconv.Itoa(j)][outgoingNode] =
									"router_" + strconv.Itoa(i)
								existing["site_"+strconv.
									Itoa(i)][outgoingNode]["router_"+strconv.
									Itoa(j)] = "router_" + strconv.Itoa(i)
								w8.WriteString(str)
							}

							if existing["router_"+strconv.Itoa(i)]["site_"+
								strconv.Itoa(j)]["router_"+strconv.Itoa(j)] ==
								"" {

								str = "$router(" + strconv.Itoa(i) +
									") add-route [$ns link $site(" +
									strconv.Itoa(j) + ") $router(" +
									strconv.Itoa(j) + ")] $router(" +
									idxNextHop + ")\n"
								existing["router_"+strconv.Itoa(i)]["site_"+
									strconv.Itoa(j)]["router_"+
									strconv.Itoa(j)] = nameNextHop
								existing["router_"+strconv.Itoa(i)]["router_"+
									strconv.Itoa(j)]["site_"+strconv.Itoa(j)] =
									nameNextHop
								w8.WriteString(str)
							}

							if existing["site_"+strconv.Itoa(i)]["site_"+
								strconv.Itoa(j)]["router_"+strconv.Itoa(j)] ==
								"" {

								str = "$site(" + strconv.Itoa(i) +
									") add-route [$ns link $site(" +
									strconv.Itoa(j) + ") $router(" +
									strconv.Itoa(j) + ")] $router(" +
									strconv.Itoa(i) + ")\n"
								existing["site_"+strconv.Itoa(i)]["site_"+
									strconv.Itoa(j)]["router_"+
									strconv.Itoa(j)] = "router_" +
									strconv.Itoa(i)
								existing["site_"+strconv.Itoa(i)]["router_"+
									strconv.Itoa(j)]["site_"+strconv.Itoa(j)] =
									"router_" + strconv.Itoa(i)
								w8.WriteString(str)
							}
						}
					} else {
						// compare which distance is better: though here or
						// through the primary
						if strconv.Itoa(j) != idxOutgoingNode {
							if shortest[namei][namej] <
								shortest[namei][firstLink[namej]]+
									shortest[firstLink[namej]][namej] &&
								existing[namei][namej][outgoingNode] == "" {

								str = "$router(" + strconv.Itoa(i) +
									") add-route [$ns link $router(" +
									strconv.Itoa(j) + ") $router(" +
									idxOutgoingNode + ")] $router(" +
									idxNextHop + ")\n"
								existing[namei][namej][outgoingNode] =
									nameNextHop
								existing[namei][outgoingNode][namej] =
									nameNextHop
								w8.WriteString(str)

								if existing["site_"+strconv.Itoa(i)]["router_"+
									strconv.Itoa(j)][outgoingNode] == "" {

									str = "$site(" + strconv.Itoa(i) +
										") add-route [$ns link $router(" +
										strconv.Itoa(j) + ") $router(" +
										idxOutgoingNode + ")] $router(" +
										strconv.Itoa(i) + ")\n"
									existing["site_"+strconv.Itoa(i)]["router_"+
										strconv.Itoa(j)][outgoingNode] =
										"router_" + strconv.Itoa(i)
									existing["site_"+strconv.
										Itoa(i)][outgoingNode]["router_"+
										strconv.Itoa(j)] = "router_" +
										strconv.Itoa(i)
									w8.WriteString(str)
								}

								if existing["router_"+strconv.Itoa(i)]["site_"+
									strconv.Itoa(j)]["router_"+
									strconv.Itoa(j)] == "" {

									str = "$router(" + strconv.Itoa(i) +
										") add-route [$ns link $site(" +
										strconv.Itoa(j) + ") $router(" +
										strconv.Itoa(j) + ")] $router(" +
										idxNextHop + ")\n"
									existing["router_"+strconv.Itoa(i)]["site_"+
										strconv.Itoa(j)]["router_"+
										strconv.Itoa(j)] = nameNextHop
									existing["router_"+strconv.
										Itoa(i)]["router_"+
										strconv.Itoa(j)]["site_"+
										strconv.Itoa(j)] = nameNextHop
									w8.WriteString(str)
								}

								if existing["site_"+strconv.Itoa(i)]["site_"+
									strconv.Itoa(j)]["router_"+
									strconv.Itoa(j)] == "" {

									str = "$site(" + strconv.Itoa(i) +
										") add-route [$ns link $site(" +
										strconv.Itoa(j) + ") $router(" +
										strconv.Itoa(j) + ")] $router(" +
										strconv.Itoa(i) + ")\n"
									existing["site_"+strconv.Itoa(i)]["site_"+
										strconv.Itoa(j)]["router_"+
										strconv.Itoa(j)] = "router_" +
										strconv.Itoa(i)
									existing["site_"+strconv.Itoa(i)]["router_"+
										strconv.Itoa(j)]["site_"+
										strconv.Itoa(j)] = "router_" +
										strconv.Itoa(i)
									w8.WriteString(str)
								}
							} else {
								// it'll be filled in on the other side
							}
						}

					}
				}
			}
		}
	}

	//fmt.Println("1 2", shortest["site_0"]["site_1"],
	//shortest["site_0"]["site_2"])

	w8.WriteString("\n" + lastPart)

	w8.Flush()
	file8.Close()

	return firstLink
}

// ReadFileLineByLine ReadFileLineByLine
func ReadFileLineByLine(configFilePath string) (func() string, error) {
	f, err := os.Open(configFilePath)
	//defer close(f)

	if err != nil {
		return func() string { return "" }, err
	}
	checkErr(err)
	reader := bufio.NewReader(f)
	//defer close(reader)
	var line string
	return func() string {
		if err == io.EOF {
			return ""
		}
		line, err = reader.ReadString('\n')
		checkErr(err)
		line = strings.Split(line, "\n")[0]
		return line
	}, nil
}

func checkErr(e error) {
	if e != nil && e != io.EOF {
		panic(e)
	}
}

func shortestNSRoutes() {
	dist, N, R, directions, _ := readNodes()

	//file7, _ := os.Create("shortest.txt")
	//w7 := bufio.NewWriter(file7)
	shortest, next := floydWarshall(N, R, dist)

	/*
		for n1, m := range shortest {
			for n2 := range m {
				if !strings.Contains(n1, "router") &&
					!strings.Contains(n2, "router") {

					//w7.WriteString("ping " + n1 + " " + n2 + " = " +
					//fmt.Sprintf("%.2f", d) + "\n")
				}
			}
		}

		//w7.Flush()
		//file7.Close()

		//fmt.Println(shortest)
	*/

	printShortestNsPaths(N, R, shortest, next, directions)

	/*
		fmt.Println(ipLInks)

		for i := 0; i < N; i++ {
			namei := "node_" + strconv.Itoa(i)
			fmt.Println(namei, ipLInks[namei][firstLinks[namei]])
		}
	*/

}

// DeterNode structure
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

// maxIfaces MAX_IFACES
var maxIfaces = 4

func dist(a DeterNode, b DeterNode) float64 {
	return math.Sqrt(math.Pow(float64(a.X-b.X), 2.0) +
		math.Pow(float64(a.Y-b.Y), 2.0))
}

func genAndPrintRndRouters(N int, R int, SpaceMax int, K int, zeroY bool,
	trueLatency bool) {

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
		routers[i].Y = RndSrc.Float64() * float64(SpaceMax)
		//routers[i].Y = RndSrc.Float64() * float64(2)

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

	file2, _ := os.Create("deter.ns")
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
	//w2.WriteString("\ttb-set-hardware $site($i) {MicroCloud}\n")dl380g3
	//w2.WriteString("\ttb-set-hardware $site($i) {bpc2133}\n")
	w2.WriteString("\ttb-set-hardware $site($i) {dl380g3}\n")

	w2.WriteString("\ttb-set-node-os $site($i) Ubuntu1404-64-STD\n")

	w2.WriteString("}\n")

	// routers
	w2.WriteString("for {set i 0} {$i < $n_routers} {incr i} {\n")

	w2.WriteString("\tset router($i) [$ns node]\n")
	//w2.WriteString("\ttb-set-hardware $router($i) {MicroCloud}\n")
	w2.WriteString("\ttb-set-hardware $router($i) {bpc2133}\n")
	//w2.WriteString("\ttb-set-hardware $router($i) {dl380g3}\n")
	w2.WriteString("\ttb-set-node-os $router($i) Ubuntu1404-64-STD\n")

	w2.WriteString("}\n\n")

	// make sure the network is connected

	linkNr := 0

	for i := 0; i < R; i++ {
		attempts := 0
		for {

			//log.LLvl1("node", i, "attempts", attempts)
			if len(routers[i].IP) == maxIfaces || attempts == 100 {
				break
			}
			// generate a random node to connect to
			peerIdx := RndSrc.Intn(R)
			peerName := routers[peerIdx].Name

			for {
				if (peerIdx != i && !routers[i].Links[peerName] &&
					len(routers[peerIdx].IP) < maxIfaces) || attempts == 100 {

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
			dist := dist(routers[i], routers[peerIdx])

			// round distance to closest integer
			if dist < 1 {
				dist = 0
			} else if dist < 2 {
				dist = 2
			}

			//fmt.Println("generated", dist)

			w2.WriteString("set link" + strconv.Itoa(linkNr) +
				" [$ns duplex-link $router(" + idx1 + ") $router(" + idx2 +
				") 100Mb " + fmt.Sprintf("%f", dist) + "ms DropTail]\n")
			w2.WriteString("tb-set-ip-link $router(" + idx1 + ") $link" +
				strconv.Itoa(linkNr) + " 10.1." + strconv.Itoa(linkNr+1) +
				".2\n")
			w2.WriteString("tb-set-ip-link $router(" + idx2 + ") $link" +
				strconv.Itoa(linkNr) + " 10.1." + strconv.Itoa(linkNr+1) +
				".3\n")

			routers[i].Links[peerName] = true
			routers[peerIdx].Links[routers[i].Name] = true
			routers[i].Dist[peerName] = dist
			routers[peerIdx].Dist[routers[i].Name] = dist

			routers[i].IP = append(routers[i].IP, "10.1."+
				strconv.Itoa(linkNr+1)+".2")
			routers[peerIdx].IP = append(routers[peerIdx].IP, "10.1."+
				strconv.Itoa(linkNr+1)+".3")

			//log.LLvl1(i, peerIdx, dist)

			linkNr++

		}
	}

	// connect two nodes to their corresponding router through a link of latency
	// 0, respectively random latencies

	for i := 0; i < N; i++ {
		routerIdx := strconv.Itoa(i)
		nodeIdx := strconv.Itoa(i)

		latency := 0

		w2.WriteString("set link" + strconv.Itoa(linkNr) +
			" [$ns duplex-link $router(" + routerIdx + ") $site(" + nodeIdx +
			") 100Mb " + fmt.Sprintf("%d", latency) + "ms DropTail]\n")
		w2.WriteString("tb-set-ip-link $router(" + routerIdx + ") $link" +
			strconv.Itoa(linkNr) + " 10.1." + strconv.Itoa(linkNr+1) + ".2\n")
		w2.WriteString("tb-set-ip-link $site(" + nodeIdx + ") $link" +
			strconv.Itoa(linkNr) + " 10.1." + strconv.Itoa(linkNr+1) + ".3\n")

		nodes[i].HostIP = "10.1." + strconv.Itoa(linkNr+1) + ".3"

		linkNr++
	}

	w2.WriteString("\n")
	// OSPF routing
	w2.WriteString(lastPart)

	w2.Flush()
	/*

		file, _ := os.Create("out.txt")
		defer file.Close()
		w := bufio.NewWriter(file)

		// print nodes in the out experiment file
		for i := 0; i < N; i++ {
			ips := ""
			for _, ip := range nodes[i].IP {
				ips += ip + ","
			}

			w.WriteString(nodes[i].Name + " " + nodes[i].HostIP + " " +
				strconv.Itoa(nodes[i].Level) + "\n")
		}

		w.Flush()
	*/

	file3, _ := os.Create("../data/nodes.txt")
	defer file3.Close()
	w3 := bufio.NewWriter(file3)

	for i := 0; i < N; i++ {
		w3.WriteString(nodes[i].Name + " " + fmt.Sprintf("%f", routers[i].X) +
			" " + fmt.Sprintf("%f", routers[i].Y) + " " + nodes[i].HostIP +
			" " + strconv.Itoa(nodes[i].Level) + "\n")
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
	shortestNSRoutes()
	str := fmt.Sprintf("K=%d\nN=%d\nR=%d\nD=%d\n", *K, *N, *R, *SpaceMax)
	ioutil.WriteFile("../data/gen/details.txt", []byte(str), 0777)

}
