package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/dedis/paper_crux/dsn_exp/gentree"
	"go.dedis.ch/onet/log"
)

func checkErr(e error) {
	if e != nil && e != io.EOF {
		fmt.Print(e)
		panic(e)
	}
}

// ReadFileLineByLine reads a file line by line
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

// ReadNodes from the config
func (s *IPFSSimulation) ReadNodes(filename string) {

	s.Nodes.All = make([]*gentree.LocalityNode, 0)

	readLine, _ := ReadFileLineByLine(filename)

	lineNr := 0

	for true {
		line := readLine()
		//fmt.Println(line)
		if line == "" {
			//fmt.Println("end")
			break
		}

		if strings.HasPrefix(line, "#") {
			continue
		}

		tokens := strings.Split(line, " ")
		coords := strings.Split(tokens[1], ",")

		name, xStr, yStr, IP, levelStr := tokens[0], coords[0], coords[1], tokens[2], tokens[3]
		x, _ := strconv.ParseFloat(xStr, 64)
		y, _ := strconv.ParseFloat(yStr, 64)
		//name, IP, levelStr := tokens[0], tokens[1], tokens[2]

		//x := 0.0
		//y := 0.0
		level, err := strconv.Atoi(levelStr)

		if err != nil {
			log.Lvl1("Error", err)

		}

		//	log.Lvl1("reqd node level", name, level_str, "lvl", level)

		myNode := CreateLocalityNode(name, x, y, IP, level)
		fmt.Println(myNode.ServerIdentity)
		s.Nodes.All = append(s.Nodes.All, myNode)

		// TODO hack!!!
		if lineNr > 45 {
			s.Nodes.All[lineNr].Level = s.Nodes.All[lineNr%45].Level
		}
		lineNr++

	}

	//log.Lvlf1("our nodes are %v", s.Nodes.All)
}

// CreateLocalityNode creates a locality node
func CreateLocalityNode(Name string, x float64, y float64, IP string, level int) *gentree.LocalityNode {
	var myNode gentree.LocalityNode

	myNode.X = x
	myNode.Y = y
	myNode.Name = Name
	myNode.IP = make(map[string]bool)

	tokens := strings.Split(IP, ",")
	for _, t := range tokens {
		myNode.IP[t] = true
	}

	myNode.Level = level
	myNode.ADist = make([]float64, 0)
	myNode.PDist = make([]string, 0)
	myNode.Cluster = make(map[string]bool)
	myNode.Bunch = make(map[string]bool)
	myNode.Rings = make([]string, 0)
	return &myNode
}
