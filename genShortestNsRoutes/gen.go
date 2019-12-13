package main

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"os"
	"strconv"
	"strings"
)


func readNodes() (map[string]map[string]float64, int, int, map[string]map[string]bool, map[string]map[string]string) {

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

			fmt.Println(nrNodes)
		}

		if strings.HasPrefix(line, "set link") {
			words := strings.Split(line, " ")
			node1, node2 := words[4], words[5]
			// $site(0)
			node1idx, _:= strconv.Atoi(strings.Split(strings.Split(node1, "(")[1], ")")[0])
			name1 := strings.Split(strings.Split(node1, "(")[0], "$")[1] + "_" + strconv.Itoa(node1idx)

			// $site(1)
			fmt.Println(node1, node2)
			node2idx, _:= strconv.Atoi(strings.Split(strings.Split(node2, "(")[1], ")")[0])
			name2 := strings.Split(strings.Split(node2, "(")[0], "$")[1] + "_" + strconv.Itoa(node2idx)

			print(name1, name2)

			directions[name1][name2] = true
			directions[name2][name1] = true

			distStr := words[7]
			fmt.Println(distStr, distStr[:len(distStr)-2])
			dist, err := strconv.ParseFloat(distStr[:len(distStr)-2], 64)
			if err != nil {
				fmt.Println("Problem when parsing pings", distStr, distStr[:len(distStr)-2])
			}

			nodes[name1][name2] = dist
			nodes[name2][name1] = dist

			lineNode1 := readLine()
			words2 := strings.Split(lineNode1, " ")
			node1iface, ip := words2[1], words2[3]
			node1ifaceIdx, _ := strconv.Atoi(strings.Split(strings.Split(node1iface, "(")[1], ")")[0])
			if node1ifaceIdx == node1idx {
				ipLinks[name1][name2] = ip
			}
			if node1ifaceIdx == node2idx {
				ipLinks[name2][name1] = ip
			}


			lineNode2 := readLine()
			words3 := strings.Split(lineNode2, " ")
			node2iface, ip := words3[1], words3[3]
			node2ifaceIdx, _ := strconv.Atoi(strings.Split(strings.Split(node2iface, "(")[1], ")")[0])
			if node2ifaceIdx == node1idx {
				ipLinks[name1][name2] = ip
			}
			if node2ifaceIdx == node2idx {
				ipLinks[name2][name1] = ip
			}


		}

	}


	fmt.Println("blablabla")

	fmt.Println("here",nodes)

	return nodes, nrNodes, nrRouters, directions, ipLinks
}


func readPings() map[string]map[string]float64{

	pings :=  make(map[string]map[string]float64)

	// read from file lines of fomrm "ping node_19 node_7 = 32.317"
	readLine,_ := ReadFileLineByLine("pings.txt")

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

		pings[src][dst] = math.Round(pingTime*100)/100
	}

	return pings
}

func floydWarshall(N int, R int, dist map[string]map[string]float64) (map[string]map[string]float64, map[string]map[string]string) {


	shortest := make(map[string]map[string]float64)
	next := make(map[string]map[string]string)

	for i := 0 ; i < N ; i++ {
		name := "site_" + strconv.Itoa(i)
		shortest[name] = make(map[string]float64)
		next[name] = make(map[string]string)
	}

	for i := 0 ; i < R ; i++ {
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

	for x,m := range dist {
		for y,d := range m {
			shortest[x][y] = d
			shortest[x][x] = 0
			next[x][y] = y
			next[x][x] = ""
		}
	}


	for k := 0 ; k < N; k++ {
		namek := "site_" + strconv.Itoa(k)
		for i := 0; i < N; i++ {
			namei := "site_" + strconv.Itoa(i)
			for j := 0; j < N; j++ {
				namej := "site_" + strconv.Itoa(j)
				if shortest[namei][namej] > shortest[namei][namek]+shortest[namek][namej] {
					shortest[namei][namej] = shortest[namei][namek] + shortest[namek][namej]
					next[namei][namej] = next[namei][namek]
				}
			}
			for j := 0; j < R; j++ {
				namej := "router_" + strconv.Itoa(j)
				if shortest[namei][namej] > shortest[namei][namek]+shortest[namek][namej] {
					shortest[namei][namej] = shortest[namei][namek] + shortest[namek][namej]
					next[namei][namej] = next[namei][namek]
				}
			}
		}
		for i := 0; i < R; i++ {
			namei := "router_" + strconv.Itoa(i)
			for j := 0; j < N; j++ {
				namej := "site_" + strconv.Itoa(j)
				if shortest[namei][namej] > shortest[namei][namek]+shortest[namek][namej] {
					shortest[namei][namej] = shortest[namei][namek] + shortest[namek][namej]
					next[namei][namej] = next[namei][namek]
				}
			}
			for j := 0; j < R; j++ {
				namej := "router_" + strconv.Itoa(j)
				if shortest[namei][namej] > shortest[namei][namek]+shortest[namek][namej] {
					shortest[namei][namej] = shortest[namei][namek] + shortest[namek][namej]
					next[namei][namej] = next[namei][namek]
				}
			}
		}
	}

	for k := 0 ; k < R; k++ {
		namek := "router_" + strconv.Itoa(k)
		for i := 0 ; i < N; i++ {
			namei := "site_" + strconv.Itoa(i)
			for j := 0; j < N; j++ {
				namej := "site_" + strconv.Itoa(j)
				if shortest[namei][namej] > shortest[namei][namek]+shortest[namek][namej] {
					shortest[namei][namej] = shortest[namei][namek] + shortest[namek][namej]
					next[namei][namej] = next[namei][namek]
				}
			}
			for j := 0; j < R; j++ {
				namej := "router_" + strconv.Itoa(j)
				if shortest[namei][namej] > shortest[namei][namek]+shortest[namek][namej] {
					shortest[namei][namej] = shortest[namei][namek] + shortest[namek][namej]
					next[namei][namej] = next[namei][namek]
				}
			}
		}
		for i := 0 ; i < R; i++ {
			namei := "router_" + strconv.Itoa(i)
			for j:= 0 ; j < N; j++ {
				namej := "site_" + strconv.Itoa(j)
				if shortest[namei][namej] > shortest[namei][namek]+shortest[namek][namej] {
					shortest[namei][namej] = shortest[namei][namek] + shortest[namek][namej]
					next[namei][namej] = next[namei][namek]
				}
			}
			for j:= 0 ; j < R; j++ {
				namej := "router_" + strconv.Itoa(j)
				if shortest[namei][namej] > shortest[namei][namek] + shortest[namek][namej] {
					shortest[namei][namej] = shortest[namei][namek] + shortest[namek][namej]
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
		path = append(path,u)

		if u == v {
			break
		}
	}

	return path
}

func printShortestNsPaths(N int, R int, shortest map[string]map[string]float64, next map[string]map[string]string, directions map[string]map[string]bool) map[string]string {

	file8, _ := os.Create("shortest.ns")
	w8 := bufio.NewWriter(file8)

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

			fmt.Println(namei, namej, path, shortest[namei][namej])

			nameNextHop := path[1]
			nameLastHop := path[len(path)-2]
			idxNextHop := strings.Split(nameNextHop, "_")[1]
			idxLastHop := strings.Split(nameLastHop, "_")[1]

			str := ""

			// this is the shortest path to that point

			// if the path there is not defined yet, define it now
			if len(path) > 2 {

				if existing[namei][namej][nameLastHop] == "" {
					str = "$router(" + strconv.Itoa(i) + ") add-route [$ns link $router(" + strconv.Itoa(j) + ") $router(" + idxLastHop + ")] $router(" + idxNextHop + ")\n"
					existing[namei][namej][nameLastHop] = nameNextHop
					existing[namei][nameLastHop][namej] = nameNextHop
					w8.WriteString(str)
				}
			}

			if existing["site_" + strconv.Itoa(i)]["router_" + strconv.Itoa(j)][nameLastHop] == "" {
				str = "$site(" + strconv.Itoa(i) + ") add-route [$ns link $router(" + strconv.Itoa(j) + ") $router(" + idxLastHop + ")] $router(" + strconv.Itoa(i) + ")\n"
				existing["site_" + strconv.Itoa(i)]["router_" + strconv.Itoa(j)][nameLastHop] = "router_" + strconv.Itoa(i)
				existing["site_" + strconv.Itoa(i)][nameLastHop]["router_" + strconv.Itoa(j)] = "router_" + strconv.Itoa(i)
				w8.WriteString(str)
			}

			if existing["router_" + strconv.Itoa(i)]["site_" + strconv.Itoa(j)]["router_" + strconv.Itoa(j)] == "" {
				str = "$router(" + strconv.Itoa(i) + ") add-route [$ns link $site(" + strconv.Itoa(j) + ") $router(" + strconv.Itoa(j) + ")] $router(" + idxNextHop + ")\n"
				existing["router_" + strconv.Itoa(i)]["site_" + strconv.Itoa(j)]["router_" + strconv.Itoa(j)] = nameNextHop
				existing["router_" + strconv.Itoa(i)]["router_" + strconv.Itoa(j)]["site_" + strconv.Itoa(j)] = nameNextHop
				w8.WriteString(str)
			}

			if existing["site_" + strconv.Itoa(i)]["site_" + strconv.Itoa(j)]["router_" + strconv.Itoa(j)] == "" {
				str = "$site(" + strconv.Itoa(i) + ") add-route [$ns link $site(" + strconv.Itoa(j) + ") $router(" + strconv.Itoa(j) + ")] $router(" + strconv.Itoa(i) + ")\n"
				existing["site_" + strconv.Itoa(i)]["site_" + strconv.Itoa(j)]["router_" + strconv.Itoa(j)] = "router_" + strconv.Itoa(i)
				existing["site_" + strconv.Itoa(i)]["router_" + strconv.Itoa(j)]["site_" + strconv.Itoa(j)] = "router_" + strconv.Itoa(i)
				w8.WriteString(str)
			}


				/*
				if existing["site_" + strconv.Itoa(i*2)]["router_" + strconv.Itoa(j)][nameLastHop] == "" {
					str = "$site(" + strconv.Itoa(i*2) + ") add-route [$ns link $router(" + strconv.Itoa(j) + ") $router(" + idxLastHop + ")] $router(" + strconv.Itoa(i) + ")\n"
					existing["site_" + strconv.Itoa(i*2)]["router_" + strconv.Itoa(j)][nameLastHop] = "router_" + strconv.Itoa(i)
					existing["site_" + strconv.Itoa(i*2)][nameLastHop]["router_" + strconv.Itoa(j)] = "router_" + strconv.Itoa(i)
					w8.WriteString(str)
				}

				if existing["router_" + strconv.Itoa(i)]["site_" + strconv.Itoa(j*2)]["router_" + strconv.Itoa(j)] == "" {
					str = "$router(" + strconv.Itoa(i) + ") add-route [$ns link $site(" + strconv.Itoa(j*2) + ") $router(" + strconv.Itoa(j) + ")] $router(" + idxNextHop + ")\n"
					existing["router_" + strconv.Itoa(i)]["site_" + strconv.Itoa(j*2)]["router_" + strconv.Itoa(j)] = nameNextHop
					existing["router_" + strconv.Itoa(i)]["router_" + strconv.Itoa(j)]["site_" + strconv.Itoa(j*2)] = nameNextHop
					w8.WriteString(str)
				}

				if existing["site_" + strconv.Itoa(i*2)]["site_" + strconv.Itoa(j*2)]["router_" + strconv.Itoa(j)] == "" {
					str = "$site(" + strconv.Itoa(i*2) + ") add-route [$ns link $site(" + strconv.Itoa(j*2) + ") $router(" + strconv.Itoa(j) + ")] $router(" + strconv.Itoa(i) + ")\n"
					existing["site_" + strconv.Itoa(i*2)]["site_" + strconv.Itoa(j*2)]["router_" + strconv.Itoa(j)] = "router_" + strconv.Itoa(i)
					existing["site_" + strconv.Itoa(i*2)]["router_" + strconv.Itoa(j)]["site_" + strconv.Itoa(j*2)] = "router_" + strconv.Itoa(i)
					w8.WriteString(str)
				}

			if existing["site_" + strconv.Itoa(i*2)]["site_" + strconv.Itoa(j*2+1)]["router_" + strconv.Itoa(j)] == "" {
				str = "$site(" + strconv.Itoa(i*2) + ") add-route [$ns link $site(" + strconv.Itoa(j*2+1) + ") $router(" + strconv.Itoa(j) + ")] $router(" + strconv.Itoa(i) + ")\n"
				existing["site_" + strconv.Itoa(i*2)]["site_" + strconv.Itoa(j*2+1)]["router_" + strconv.Itoa(j)] = "router_" + strconv.Itoa(i)
				existing["site_" + strconv.Itoa(i*2)]["router_" + strconv.Itoa(j)]["site_" + strconv.Itoa(j*2+1)] = "router_" + strconv.Itoa(i)
				w8.WriteString(str)
			}

			// +1
			if existing["site_" + strconv.Itoa(i*2+1)]["router_" + strconv.Itoa(j)][nameLastHop] == "" {
				str = "$site(" + strconv.Itoa(i*2+1) + ") add-route [$ns link $router(" + strconv.Itoa(j) + ") $router(" + idxLastHop + ")] $router(" + strconv.Itoa(i) + ")\n"
				existing["site_" + strconv.Itoa(i*2+1)]["router_" + strconv.Itoa(j)][nameLastHop] = "router_" + strconv.Itoa(i)
				existing["site_" + strconv.Itoa(i*2+1)][nameLastHop]["router_" + strconv.Itoa(j)] = "router_" + strconv.Itoa(i)
				w8.WriteString(str)
			}

			if existing["router_" + strconv.Itoa(i)]["site_" + strconv.Itoa(j*2+1)]["router_" + strconv.Itoa(j)] == "" {
				str = "$router(" + strconv.Itoa(i) + ") add-route [$ns link $site(" + strconv.Itoa(j*2+1) + ") $router(" + strconv.Itoa(j) + ")] $router(" + idxNextHop + ")\n"
				existing["router_" + strconv.Itoa(i)]["site_" + strconv.Itoa(j*2+1)]["router_" + strconv.Itoa(j)] = nameNextHop
				existing["router_" + strconv.Itoa(i)]["router_" + strconv.Itoa(j)]["site_" + strconv.Itoa(j*2+1)] = nameNextHop
				w8.WriteString(str)
			}

			if existing["site_" + strconv.Itoa(i*2+1)]["site_" + strconv.Itoa(j*2+1)]["router_" + strconv.Itoa(j)] == "" {
				str = "$site(" + strconv.Itoa(i*2+1) + ") add-route [$ns link $site(" + strconv.Itoa(j*2+1) + ") $router(" + strconv.Itoa(j) + ")] $router(" + strconv.Itoa(i) + ")\n"
				existing["site_" + strconv.Itoa(i*2+1)]["site_" + strconv.Itoa(j*2+1)]["router_" + strconv.Itoa(j)] = "router_" + strconv.Itoa(i)
				existing["site_" + strconv.Itoa(i*2+1)]["router_" + strconv.Itoa(j)]["site_" + strconv.Itoa(j*2+1)] = "router_" + strconv.Itoa(i)
				w8.WriteString(str)
			}

			if existing["site_" + strconv.Itoa(i*2+1)]["site_" + strconv.Itoa(j*2)]["router_" + strconv.Itoa(j)] == "" {
				str = "$site(" + strconv.Itoa(i*2+1) + ") add-route [$ns link $site(" + strconv.Itoa(j*2) + ") $router(" + strconv.Itoa(j) + ")] $router(" + strconv.Itoa(i) + ")\n"
				existing["site_" + strconv.Itoa(i*2+1)]["site_" + strconv.Itoa(j*2)]["router_" + strconv.Itoa(j)] = "router_" + strconv.Itoa(i)
				existing["site_" + strconv.Itoa(i*2+1)]["router_" + strconv.Itoa(j)]["site_" + strconv.Itoa(j*2)] = "router_" + strconv.Itoa(i)
				w8.WriteString(str)
			}
				*/


			// but which of the outgoing links og namej, aka dest, should i also add?

			// path to main interface has to be shortest

			// is the main interface defined?
			if firstLink[namej] == "" {
				firstLink[namej] = nameLastHop
			}

			// look at all outgoing links of namej
			for outgoingNode, exists := range directions[namej] {
				if exists && outgoingNode != nameLastHop {
					idxOutgoingNode := strings.Split(outgoingNode, "_")[1]
					if firstLink[namej] != outgoingNode && existing[namei][namej][outgoingNode] == "" {

						if strconv.Itoa(j) != idxOutgoingNode {
							str = "$router(" + strconv.Itoa(i) + ") add-route [$ns link $router(" + strconv.Itoa(j) + ") $router(" + idxOutgoingNode + ")] $router(" + idxNextHop + ")\n"
							existing[namei][namej][outgoingNode] = nameNextHop
							existing[namei][outgoingNode][namej] = nameNextHop
							w8.WriteString(str)


							if existing["site_" + strconv.Itoa(i)]["router_" + strconv.Itoa(j)][outgoingNode] == "" {
								str = "$site(" + strconv.Itoa(i) + ") add-route [$ns link $router(" + strconv.Itoa(j) + ") $router(" + idxOutgoingNode + ")] $router(" + strconv.Itoa(i) + ")\n"
								existing["site_" + strconv.Itoa(i)]["router_" + strconv.Itoa(j)][outgoingNode] = "router_" + strconv.Itoa(i)
								existing["site_" + strconv.Itoa(i)][outgoingNode]["router_" + strconv.Itoa(j)] = "router_" + strconv.Itoa(i)
								w8.WriteString(str)
							}


							if existing["router_" + strconv.Itoa(i)]["site_" + strconv.Itoa(j)]["router_" + strconv.Itoa(j)] == "" {
								str = "$router(" + strconv.Itoa(i) + ") add-route [$ns link $site(" + strconv.Itoa(j) + ") $router(" + strconv.Itoa(j) + ")] $router(" + idxNextHop + ")\n"
								existing["router_" + strconv.Itoa(i)]["site_" + strconv.Itoa(j)]["router_" + strconv.Itoa(j)] = nameNextHop
								existing["router_" + strconv.Itoa(i)]["router_" + strconv.Itoa(j)]["site_" + strconv.Itoa(j)] = nameNextHop
								w8.WriteString(str)
							}

							if existing["site_" + strconv.Itoa(i)]["site_" + strconv.Itoa(j)]["router_" + strconv.Itoa(j)] == "" {
								str = "$site(" + strconv.Itoa(i) + ") add-route [$ns link $site(" + strconv.Itoa(j) + ") $router(" + strconv.Itoa(j) + ")] $router(" + strconv.Itoa(i) + ")\n"
								existing["site_" + strconv.Itoa(i)]["site_" + strconv.Itoa(j)]["router_" + strconv.Itoa(j)] = "router_" + strconv.Itoa(i)
								existing["site_" + strconv.Itoa(i)]["router_" + strconv.Itoa(j)]["site_" + strconv.Itoa(j)] = "router_" + strconv.Itoa(i)
								w8.WriteString(str)
							}

/*
							if existing["site_" + strconv.Itoa(i*2)]["router_" + strconv.Itoa(j)][outgoingNode] == "" {
								str = "$site(" + strconv.Itoa(i*2) + ") add-route [$ns link $router(" + strconv.Itoa(j) + ") $router(" + idxOutgoingNode + ")] $router(" + strconv.Itoa(i) + ")\n"
								existing["site_" + strconv.Itoa(i*2)]["router_" + strconv.Itoa(j)][outgoingNode] = "router_" + strconv.Itoa(i)
								existing["site_" + strconv.Itoa(i*2)][outgoingNode]["router_" + strconv.Itoa(j)] = "router_" + strconv.Itoa(i)
								w8.WriteString(str)
							}


							if existing["router_" + strconv.Itoa(i)]["site_" + strconv.Itoa(j*2)]["router_" + strconv.Itoa(j)] == "" {
								str = "$router(" + strconv.Itoa(i) + ") add-route [$ns link $site(" + strconv.Itoa(j*2) + ") $router(" + strconv.Itoa(j) + ")] $router(" + idxNextHop + ")\n"
								existing["router_" + strconv.Itoa(i)]["site_" + strconv.Itoa(j*2)]["router_" + strconv.Itoa(j)] = nameNextHop
								existing["router_" + strconv.Itoa(i)]["router_" + strconv.Itoa(j)]["site_" + strconv.Itoa(j*2)] = nameNextHop
								w8.WriteString(str)
							}

							if existing["site_" + strconv.Itoa(i*2)]["site_" + strconv.Itoa(j*2)]["router_" + strconv.Itoa(j)] == "" {
								str = "$site(" + strconv.Itoa(i*2) + ") add-route [$ns link $site(" + strconv.Itoa(j*2) + ") $router(" + strconv.Itoa(j) + ")] $router(" + strconv.Itoa(i) + ")\n"
								existing["site_" + strconv.Itoa(i*2)]["site_" + strconv.Itoa(j*2)]["router_" + strconv.Itoa(j)] = "router_" + strconv.Itoa(i)
								existing["site_" + strconv.Itoa(i*2)]["router_" + strconv.Itoa(j)]["site_" + strconv.Itoa(j*2)] = "router_" + strconv.Itoa(i)
								w8.WriteString(str)
							}

							if existing["site_" + strconv.Itoa(i*2)]["site_" + strconv.Itoa(j*2+1)]["router_" + strconv.Itoa(j)] == "" {
								str = "$site(" + strconv.Itoa(i*2) + ") add-route [$ns link $site(" + strconv.Itoa(j*2+1) + ") $router(" + strconv.Itoa(j) + ")] $router(" + strconv.Itoa(i) + ")\n"
								existing["site_" + strconv.Itoa(i*2)]["site_" + strconv.Itoa(j*2+1)]["router_" + strconv.Itoa(j)] = "router_" + strconv.Itoa(i)
								existing["site_" + strconv.Itoa(i*2)]["router_" + strconv.Itoa(j)]["site_" + strconv.Itoa(j*2+1)] = "router_" + strconv.Itoa(i)
								w8.WriteString(str)
							}

							// +1
							if existing["site_" + strconv.Itoa(i*2+1)]["router_" + strconv.Itoa(j)][outgoingNode] == "" {
								str = "$site(" + strconv.Itoa(i*2+1) + ") add-route [$ns link $router(" + strconv.Itoa(j) + ") $router(" + idxOutgoingNode + ")] $router(" + strconv.Itoa(i) + ")\n"
								existing["site_" + strconv.Itoa(i*2+1)]["router_" + strconv.Itoa(j)][outgoingNode] = "router_" + strconv.Itoa(i)
								existing["site_" + strconv.Itoa(i*2+1)][outgoingNode]["router_" + strconv.Itoa(j)] = "router_" + strconv.Itoa(i)
								w8.WriteString(str)
							}


							if existing["router_" + strconv.Itoa(i)]["site_" + strconv.Itoa(j*2+1)]["router_" + strconv.Itoa(j)] == "" {
								str = "$router(" + strconv.Itoa(i) + ") add-route [$ns link $site(" + strconv.Itoa(j*2+1) + ") $router(" + strconv.Itoa(j) + ")] $router(" + idxNextHop + ")\n"
								existing["router_" + strconv.Itoa(i)]["site_" + strconv.Itoa(j*2+1)]["router_" + strconv.Itoa(j)] = nameNextHop
								existing["router_" + strconv.Itoa(i)]["router_" + strconv.Itoa(j)]["site_" + strconv.Itoa(j*2+1)] = nameNextHop
								w8.WriteString(str)
							}

							if existing["site_" + strconv.Itoa(i*2+1)]["site_" + strconv.Itoa(j*2+1)]["router_" + strconv.Itoa(j)] == "" {
								str = "$site(" + strconv.Itoa(i*2+1) + ") add-route [$ns link $site(" + strconv.Itoa(j*2+1) + ") $router(" + strconv.Itoa(j) + ")] $router(" + strconv.Itoa(i) + ")\n"
								existing["site_" + strconv.Itoa(i*2+1)]["site_" + strconv.Itoa(j*2+1)]["router_" + strconv.Itoa(j)] = "router_" + strconv.Itoa(i)
								existing["site_" + strconv.Itoa(i*2+1)]["router_" + strconv.Itoa(j)]["site_" + strconv.Itoa(j*2+1)] = "router_" + strconv.Itoa(i)
								w8.WriteString(str)
							}

							if existing["site_" + strconv.Itoa(i*2+1)]["site_" + strconv.Itoa(j*2)]["router_" + strconv.Itoa(j)] == "" {
								str = "$site(" + strconv.Itoa(i*2+1) + ") add-route [$ns link $site(" + strconv.Itoa(j*2) + ") $router(" + strconv.Itoa(j) + ")] $router(" + strconv.Itoa(i) + ")\n"
								existing["site_" + strconv.Itoa(i*2+1)]["site_" + strconv.Itoa(j*2)]["router_" + strconv.Itoa(j)] = "router_" + strconv.Itoa(i)
								existing["site_" + strconv.Itoa(i*2+1)]["router_" + strconv.Itoa(j)]["site_" + strconv.Itoa(j*2)] = "router_" + strconv.Itoa(i)
								w8.WriteString(str)
							}
*/

						}
					} else {
						// compare which distance is better: though here or through the primary
						if strconv.Itoa(j) != idxOutgoingNode {
							if shortest[namei][namej] < shortest[namei][firstLink[namej]]+shortest[firstLink[namej]][namej] && existing[namei][namej][outgoingNode] == "" {
								str = "$router(" + strconv.Itoa(i) + ") add-route [$ns link $router(" + strconv.Itoa(j) + ") $router(" + idxOutgoingNode + ")] $router(" + idxNextHop + ")\n"
								existing[namei][namej][outgoingNode] = nameNextHop
								existing[namei][outgoingNode][namej] = nameNextHop
								w8.WriteString(str)

								if existing["site_" + strconv.Itoa(i)]["router_" + strconv.Itoa(j)][outgoingNode] == "" {
									str = "$site(" + strconv.Itoa(i) + ") add-route [$ns link $router(" + strconv.Itoa(j) + ") $router(" + idxOutgoingNode + ")] $router(" + strconv.Itoa(i) + ")\n"
									existing["site_" + strconv.Itoa(i)]["router_" + strconv.Itoa(j)][outgoingNode] = "router_" + strconv.Itoa(i)
									existing["site_" + strconv.Itoa(i)][outgoingNode]["router_" + strconv.Itoa(j)] = "router_" + strconv.Itoa(i)
									w8.WriteString(str)
								}


								if existing["router_" + strconv.Itoa(i)]["site_" + strconv.Itoa(j)]["router_" + strconv.Itoa(j)] == "" {
									str = "$router(" + strconv.Itoa(i) + ") add-route [$ns link $site(" + strconv.Itoa(j) + ") $router(" + strconv.Itoa(j) + ")] $router(" + idxNextHop + ")\n"
									existing["router_" + strconv.Itoa(i)]["site_" + strconv.Itoa(j)]["router_" + strconv.Itoa(j)] = nameNextHop
									existing["router_" + strconv.Itoa(i)]["router_" + strconv.Itoa(j)]["site_" + strconv.Itoa(j)] = nameNextHop
									w8.WriteString(str)
								}

								if existing["site_" + strconv.Itoa(i)]["site_" + strconv.Itoa(j)]["router_" + strconv.Itoa(j)] == "" {
									str = "$site(" + strconv.Itoa(i) + ") add-route [$ns link $site(" + strconv.Itoa(j) + ") $router(" + strconv.Itoa(j) + ")] $router(" + strconv.Itoa(i) + ")\n"
									existing["site_" + strconv.Itoa(i)]["site_" + strconv.Itoa(j)]["router_" + strconv.Itoa(j)] = "router_" + strconv.Itoa(i)
									existing["site_" + strconv.Itoa(i)]["router_" + strconv.Itoa(j)]["site_" + strconv.Itoa(j)] = "router_" + strconv.Itoa(i)
									w8.WriteString(str)
								}

								/*
								if existing["site_" + strconv.Itoa(i*2)]["router_" + strconv.Itoa(j)][outgoingNode] == "" {
									str = "$site(" + strconv.Itoa(i*2) + ") add-route [$ns link $router(" + strconv.Itoa(j) + ") $router(" + idxOutgoingNode + ")] $router(" + strconv.Itoa(i) + ")\n"
									existing["site_" + strconv.Itoa(i*2)]["router_" + strconv.Itoa(j)][outgoingNode] = "router_" + strconv.Itoa(i)
									existing["site_" + strconv.Itoa(i*2)][outgoingNode]["router_" + strconv.Itoa(j)] = "router_" + strconv.Itoa(i)
									w8.WriteString(str)
								}


								if existing["router_" + strconv.Itoa(i)]["site_" + strconv.Itoa(j*2)]["router_" + strconv.Itoa(j)] == "" {
									str = "$router(" + strconv.Itoa(i) + ") add-route [$ns link $site(" + strconv.Itoa(j*2) + ") $router(" + strconv.Itoa(j) + ")] $router(" + idxNextHop + ")\n"
									existing["router_" + strconv.Itoa(i)]["site_" + strconv.Itoa(j*2)]["router_" + strconv.Itoa(j)] = nameNextHop
									existing["router_" + strconv.Itoa(i)]["router_" + strconv.Itoa(j)]["site_" + strconv.Itoa(j*2)] = nameNextHop
									w8.WriteString(str)
								}

								if existing["site_" + strconv.Itoa(i*2)]["site_" + strconv.Itoa(j*2)]["router_" + strconv.Itoa(j)] == "" {
									str = "$site(" + strconv.Itoa(i*2) + ") add-route [$ns link $site(" + strconv.Itoa(j*2) + ") $router(" + strconv.Itoa(j) + ")] $router(" + strconv.Itoa(i) + ")\n"
									existing["site_" + strconv.Itoa(i*2)]["site_" + strconv.Itoa(j*2)]["router_" + strconv.Itoa(j)] = "router_" + strconv.Itoa(i)
									existing["site_" + strconv.Itoa(i*2)]["router_" + strconv.Itoa(j)]["site_" + strconv.Itoa(j*2)] = "router_" + strconv.Itoa(i)
									w8.WriteString(str)
								}

								if existing["site_" + strconv.Itoa(i*2)]["site_" + strconv.Itoa(j*2+1)]["router_" + strconv.Itoa(j)] == "" {
									str = "$site(" + strconv.Itoa(i*2) + ") add-route [$ns link $site(" + strconv.Itoa(j*2+1) + ") $router(" + strconv.Itoa(j) + ")] $router(" + strconv.Itoa(i) + ")\n"
									existing["site_" + strconv.Itoa(i*2)]["site_" + strconv.Itoa(j*2+1)]["router_" + strconv.Itoa(j)] = "router_" + strconv.Itoa(i)
									existing["site_" + strconv.Itoa(i*2)]["router_" + strconv.Itoa(j)]["site_" + strconv.Itoa(j*2+1)] = "router_" + strconv.Itoa(i)
									w8.WriteString(str)
								}

								// +1
								if existing["site_" + strconv.Itoa(i*2+1)]["router_" + strconv.Itoa(j)][outgoingNode] == "" {
									str = "$site(" + strconv.Itoa(i*2+1) + ") add-route [$ns link $router(" + strconv.Itoa(j) + ") $router(" + idxOutgoingNode + ")] $router(" + strconv.Itoa(i) + ")\n"
									existing["site_" + strconv.Itoa(i*2+1)]["router_" + strconv.Itoa(j)][outgoingNode] = "router_" + strconv.Itoa(i)
									existing["site_" + strconv.Itoa(i*2+1)][outgoingNode]["router_" + strconv.Itoa(j)] = "router_" + strconv.Itoa(i)
									w8.WriteString(str)
								}


								if existing["router_" + strconv.Itoa(i)]["site_" + strconv.Itoa(j*2+1)]["router_" + strconv.Itoa(j)] == "" {
									str = "$router(" + strconv.Itoa(i) + ") add-route [$ns link $site(" + strconv.Itoa(j*2+1) + ") $router(" + strconv.Itoa(j) + ")] $router(" + idxNextHop + ")\n"
									existing["router_" + strconv.Itoa(i)]["site_" + strconv.Itoa(j*2+1)]["router_" + strconv.Itoa(j)] = nameNextHop
									existing["router_" + strconv.Itoa(i)]["router_" + strconv.Itoa(j)]["site_" + strconv.Itoa(j*2+1)] = nameNextHop
									w8.WriteString(str)
								}

								if existing["site_" + strconv.Itoa(i*2+1)]["site_" + strconv.Itoa(j*2+1)]["router_" + strconv.Itoa(j)] == "" {
									str = "$site(" + strconv.Itoa(i*2+1) + ") add-route [$ns link $site(" + strconv.Itoa(j*2+1) + ") $router(" + strconv.Itoa(j) + ")] $router(" + strconv.Itoa(i) + ")\n"
									existing["site_" + strconv.Itoa(i*2+1)]["site_" + strconv.Itoa(j*2+1)]["router_" + strconv.Itoa(j)] = "router_" + strconv.Itoa(i)
									existing["site_" + strconv.Itoa(i*2+1)]["router_" + strconv.Itoa(j)]["site_" + strconv.Itoa(j*2+1)] = "router_" + strconv.Itoa(i)
									w8.WriteString(str)
								}

								if existing["site_" + strconv.Itoa(i*2+1)]["site_" + strconv.Itoa(j*2)]["router_" + strconv.Itoa(j)] == "" {
									str = "$site(" + strconv.Itoa(i*2+1) + ") add-route [$ns link $site(" + strconv.Itoa(j*2) + ") $router(" + strconv.Itoa(j) + ")] $router(" + strconv.Itoa(i) + ")\n"
									existing["site_" + strconv.Itoa(i*2+1)]["site_" + strconv.Itoa(j*2)]["router_" + strconv.Itoa(j)] = "router_" + strconv.Itoa(i)
									existing["site_" + strconv.Itoa(i*2+1)]["router_" + strconv.Itoa(j)]["site_" + strconv.Itoa(j*2)] = "router_" + strconv.Itoa(i)
									w8.WriteString(str)
								}
								*/
							} else {
								// it'll be filled in on the other side
							}
						}

					}
				}
			}
		}


	}


	/*
	for i := 0; i < R; i++ {
		str := "$site(" + strconv.Itoa(i*2) + ") add-route [$ns link $site(" + strconv.Itoa(i*2+1) + ") $router(" + strconv.Itoa(i) + ")] $router(" + strconv.Itoa(i) + ")\n"
		w8.WriteString(str)
		str = "$site(" + strconv.Itoa(i*2+1) + ") add-route [$ns link $site(" + strconv.Itoa(i*2) + ") $router(" + strconv.Itoa(i) + ")] $router(" + strconv.Itoa(i) + ")\n"
		w8.WriteString(str)
	}
	*/

	/*
	for i := 0; i < N; i++ {

		namei := "site_" + strconv.Itoa(i)
		for j := 0; j < N; j++ {
			if i == j {
				continue
			}

			namej := "router_" + strconv.Itoa(j)

			if shortest[namei][namej] > 10000 {
				panic("no route!")
			}

			path := findPath(namei, namej, next)

			fmt.Println(namei, namej, path, shortest[namei][namej])

			nameNextHop := path[1]
			nameLastHop := path[len(path)-2]
			idxNextHop := strings.Split(nameNextHop, "_")[1]
			idxLastHop := strings.Split(nameLastHop, "_")[1]

			str := ""

			// this is the shortest path to that point

			// if the path there is not defined yet, define it now
			if len(path) > 2 {

				if existing[namei][namej][nameLastHop] == "" {
					str = "$site(" + strconv.Itoa(i) + ") add-route [$ns link $router(" + strconv.Itoa(j) + ") $router(" + idxLastHop + ")] $router(" + idxNextHop + ")\n"
					existing[namei][namej][nameLastHop] = nameNextHop
					existing[namei][nameLastHop][namej] = nameNextHop
					w8.WriteString(str)
				}
			}

			// but which of the outgoing links og namej, aka dest, should i also add?

			// path to main interface has to be shortest

			// is the main interface defined?
			if firstLink[namej] == "" {
				firstLink[namej] = nameLastHop
			}

			// look at all outgoing links of namej
			for outgoingNode, exists := range directions[namej] {
				if exists && outgoingNode != nameLastHop {
					idxOutgoingNode := strings.Split(outgoingNode, "_")[1]
					if firstLink[namej] != outgoingNode && existing[namei][namej][outgoingNode] == "" {

						str = "$site(" + strconv.Itoa(i) + ") add-route [$ns link $router(" + strconv.Itoa(j) + ") $router(" + idxOutgoingNode + ")] $router(" + idxNextHop + ")\n"
						existing[namei][namej][outgoingNode] = nameNextHop
						existing[namei][outgoingNode][namej] = nameNextHop
						w8.WriteString(str)
					} else {
						// compare which distance is better: though here or through the primary
						if shortest[namei][namej] < shortest[namei][firstLink[namej]]+shortest[firstLink[namej]][namej] && existing[namei][namej][outgoingNode] == "" {
							str = "$site(" + strconv.Itoa(i) + ") add-route [$ns link $router(" + strconv.Itoa(j) + ") $router(" + idxOutgoingNode + ")] $router(" + idxNextHop + ")\n"
							existing[namei][namej][outgoingNode] = nameNextHop
							existing[namei][outgoingNode][namej] = nameNextHop
							w8.WriteString(str)
						} else {
							// it'll be filled in on the other side
						}

					}
				}
			}
		}
	}

	for i := 0; i < N; i++ {

		namei := "site_" + strconv.Itoa(i)
		for j := 0; j < N; j++ {
			if i == j {
				continue
			}


			namej := "site_" + strconv.Itoa(j)



			if shortest[namei][namej] > 10000 {
				panic("no route!")
			}

			path := findPath(namei, namej, next)

			fmt.Println(namei, namej, path, shortest[namei][namej])

			nameNextHop := path[1]
			nameLastHop := path[len(path)-2]
			idxNextHop := strings.Split(nameNextHop, "_")[1]
			idxLastHop := strings.Split(nameLastHop, "_")[1]

			str := ""

			// this is the shortest path to that point

			// if the path there is not defined yet, define it now
			if len(path) > 2 {

				if existing[namei][namej][nameLastHop] == "" {
					str = "$site(" + strconv.Itoa(i) + ") add-route [$ns link $site(" + strconv.Itoa(j) + ") $router(" + idxLastHop + ")] $router(" + idxNextHop + ")\n"
					existing[namei][namej][nameLastHop] = nameNextHop
					existing[namei][nameLastHop][namej] = nameNextHop
					w8.WriteString(str)
				}
			}

			// but which of the outgoing links og namej, aka dest, should i also add?

			// path to main interface has to be shortest

			// is the main interface defined?
			if firstLink[namej] == "" {
				firstLink[namej] = nameLastHop
			}

			// look at all outgoing links of namej
			for outgoingNode, exists := range directions[namej] {
				if exists && outgoingNode != nameLastHop {
					idxOutgoingNode := strings.Split(outgoingNode, "_")[1]
					if firstLink[namej] != outgoingNode && existing[namei][namej][outgoingNode] == ""{

						str = "$site(" + strconv.Itoa(i) + ") add-route [$ns link $site(" + strconv.Itoa(j) + ") $router(" + idxOutgoingNode + ")] $router(" + idxNextHop + ")\n"
						existing[namei][namej][outgoingNode] = nameNextHop
						existing[namei][outgoingNode][namej] = nameNextHop
						w8.WriteString(str)
					} else {
						// compare which distance is better: though here or through the primary
						if shortest[namei][namej] < shortest[namei][firstLink[namej]]+shortest[firstLink[namej]][namej] && existing[namei][namej][outgoingNode] == "" {
							str = "$site(" + strconv.Itoa(i) + ") add-route [$ns link $site(" + strconv.Itoa(j) + ") $router(" + idxOutgoingNode + ")] $router(" + idxNextHop + ")\n"
							existing[namei][namej][outgoingNode] = nameNextHop
							existing[namei][outgoingNode][namej] = nameNextHop
							w8.WriteString(str)
						} else {
							// it'll be filled in on the other side
						}

					}
				}
			}

		}
	}

	*/
	fmt.Println("1 2", shortest["site_0"]["site_1"], shortest["site_0"]["site_2"])

	w8.Flush()
	file8.Close()

	return firstLink
}


func ReadFileLineByLine(configFilePath string) (func() string, error) {
	f, err := os.Open(configFilePath)
	//defer close(f)

	if err != nil {
		return func() string {return ""}, err
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
		fmt.Print(e)
		panic(e)
	}
}

func main() {
	dist, N, R, directions, ipLInks := readNodes()
	//dist, N, _,_ := readNodes()
	//pings := readPings()

	file7, _ := os.Create("shortest.txt")
	w7 := bufio.NewWriter(file7)
	shortest, next := floydWarshall(N, R, dist)
	//shortest, _ := floydWarshall(N, dist)

	for n1, m := range shortest {
		for n2, d := range m {
			if !strings.Contains(n1, "router") && !strings.Contains(n2, "router") {
				w7.WriteString("ping " + n1 + " " + n2 + " = " + fmt.Sprintf("%.2f", d) + "\n")
			}
		}
	}

	w7.Flush()
	file7.Close()

	//ping node_34 node_15 = 16.64
	//ping node_34 node_24 = 20.64
	//ping node_34 node_41 = 12.32




	fmt.Println(shortest)

	/*
	for i := 0 ; i < len(dist); i++ {
		namei := "node_" + strconv.Itoa(i)
		for j := 0; j < len(dist); j++ {
			namej := "node_" + strconv.Itoa(j)

			if pings[namei][namej] > shortest[namei][namej] {
				//fmt.Println("TIV ", namei, namej, pings[namei][namej], shortest[namei][namej])
			}
		}
	}*/

	firstLinks := printShortestNsPaths(N, R, shortest, next, directions)


	fmt.Println(ipLInks)

	for i := 0 ; i < N ; i++ {
		namei := "node_" + strconv.Itoa(i)
		//fmt.Println("here on", namei, "with main", firstLinks[namei])
		fmt.Println(namei, ipLInks[namei][firstLinks[namei]])
	}


}
