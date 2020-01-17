package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"go.dedis.ch/onet/v3"
	"go.dedis.ch/onet/v3/log"
	"go.dedis.ch/onet/v3/network"

	"github.com/dedis/student_19_cruxIPFS/gentree"
	"github.com/dedis/student_19_cruxIPFS/service"
)

// ReadNodeInfo read node information
func (s *IPFSSimulation) ReadNodeInfo(isLocalTest bool,
	config onet.SimulationConfig) {
	_, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	s.ReadNodesFromFile(nodesFile, config)
}

// ReadNodesFromFile read nodes information from a text file
func (s *IPFSSimulation) ReadNodesFromFile(filename string,
	config onet.SimulationConfig) {
	s.Nodes.All = make([]*gentree.LocalityNode, 0)

	readLine := service.ReadFileLineByLine(filename)

	for i := 0; i < len(config.Roster.List); i++ {
		line := readLine()
		if line == "" {
			break
		}

		if strings.HasPrefix(line, "#") {
			continue
		}

		tokens := strings.Split(line, " ")
		name, levelstr := tokens[0], tokens[4]

		level, err := strconv.Atoi(levelstr)

		if err != nil {
			log.Error(err)
		}

		myNode := gentree.CreateNode(name, level)
		s.Nodes.All = append(s.Nodes.All, myNode)
	}
}

func (s *IPFSSimulation) initializeMaps(config *onet.SimulationConfig,
	isLocalTest bool) map[*network.ServerIdentity]string {

	s.Nodes.ServerIdentityToName = make(map[network.ServerIdentityID]string)
	ServerIdentityToName := make(map[*network.ServerIdentity]string)

	if isLocalTest {
		for i := range s.Nodes.All {
			treeNode := config.Tree.List()[i]
			s.Nodes.All[i].ServerIdentity = treeNode.ServerIdentity
			s.Nodes.ServerIdentityToName[treeNode.ServerIdentity.ID] =
				s.Nodes.All[i].Name
			ServerIdentityToName[treeNode.ServerIdentity] = s.Nodes.All[i].Name
		}
	} else {
		for _, treeNode := range config.Tree.List() {
			serverIP := treeNode.ServerIdentity.Address.Host()
			node := s.Nodes.GetByIP(serverIP)
			node.ServerIdentity = treeNode.ServerIdentity
			s.Nodes.ServerIdentityToName[treeNode.ServerIdentity.ID] = node.Name
			ServerIdentityToName[treeNode.ServerIdentity] = node.Name
		}
	}

	return ServerIdentityToName
}

func parseParams() {
	b, err := ioutil.ReadFile(detailsFile)
	if err != nil {
		panic(err)
	}
	truestr := "true"
	for _, l := range strings.Split(string(b), "\n") {
		if strings.Contains(l, "pings") {
			if strings.Split(l, "=")[1] == truestr {
				computePings = true
			} else {
				computePings = false
			}
		}
		if strings.Contains(l, "mode") {
			mode = strings.Split(l, "=")[1]
		}
		if strings.Contains(l, "cruxified") {
			if strings.Split(l, "=")[1] == truestr {
				cruxified = true
			} else {
				cruxified = false
			}
		}
		if strings.Contains(l, "remote") {
			if strings.Split(l, "=")[1] == truestr {
				remote = true
			} else {
				remote = false
			}
		}
		if strings.Contains(l, "ops") {
			nOps, err = strconv.Atoi(strings.Split(l, "=")[1])
			if err != nil {
				panic(err)
			}
		}
	}
}
