package main

import (
	"strings"

	"go.dedis.ch/onet/log"
)

// ParseSimpleNodesFile parse a simple nodes description file to a list of
// string IP addresses
func (s *IPFSSimulation) ParseSimpleNodesFile(filename string,
	nrMachines, nrHosts int) []string {
	list := make([]string, nrHosts)

	readLine, _ := ReadFileLineByLine(filename)

	i := 0

	for {
		line := readLine()
		//fmt.Println(line)
		if line == "" {
			//fmt.Println("end")
			break
		}

		if strings.HasPrefix(line, "#") {
			continue
		}

		// take a set of nodes with complete latencies

		startFullAddr := 0
		if s.Rounds > 1 {
			startFullAddr = nrMachines
		}

		if i >= startFullAddr && i < startFullAddr+nrMachines {
			tokens := strings.Split(line, " ")
			//IP := tokens[2]
			IP := tokens[1]

			tokens2 := strings.Split(IP, ",")

			j := 0
			for k, t := range tokens2 {
				if t != "" && k < s.Rounds {

					log.Lvl2("i=", i, "j=", j, "sum is", i+j)

					list[i-startFullAddr+j] = t
					j += startFullAddr
				}
			}
		}
		if i == startFullAddr+nrMachines {
			break
		}
		i++
	}
	return list
}
