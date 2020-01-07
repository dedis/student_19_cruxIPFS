package service

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"go.dedis.ch/onet/v3/log"
	"go.dedis.ch/onet/v3/network"

	cruxIPFS "github.com/dedis/student_19_cruxIPFS"
)

// ExecReqPings sends all its own pings distances to the node that requested it
func (s *Service) ExecReqPings(env *network.Envelope) error {

	// Parse message
	req, ok := env.Msg.(*ReqPings)
	if !ok {
		log.Error(s.ServerIdentity(), "failed to cast to ReqPings")
		return errors.New(s.ServerIdentity().String() +
			" failed to cast to ReqPings")
	}

	// wait for pings to be finished
	for !s.DonePing {
		time.Sleep(5 * time.Second)
	}

	reply := ""
	myName := s.Nodes.GetServerIdentityToName(s.ServerIdentity())
	// build reply
	for peerName, pingTime := range s.OwnPings {
		reply += myName + " " + peerName + " " +
			fmt.Sprintf("%f", pingTime) + "\n"
	}

	log.Lvl3("sending", reply)
	requesterIdentity := s.Nodes.GetByName(req.SenderName).ServerIdentity

	e := s.SendRaw(requesterIdentity, &ReplyPings{
		Pings:      reply,
		SenderName: myName,
	})

	if e != nil {
		panic(e)
	}

	return e
}

// ExecReplyPings handle replies of other nodes ping distances to get the full
// array of distances between each nodes
func (s *Service) ExecReplyPings(env *network.Envelope) error {

	// Parse message
	req, ok := env.Msg.(*ReplyPings)
	if !ok {
		log.Error(s.ServerIdentity(), "failed to cast to ReplyPings")
		return errors.New(s.ServerIdentity().String() +
			" failed to cast to ReplyPings")
	}

	// process ping output
	s.PingMapMtx.Lock()
	lines := strings.Split(req.Pings, "\n")

	for _, line := range lines {
		if line != "" {
			//log.LLvl1("line=", line)
			words := strings.Split(line, " ")
			src := words[0]
			dst := words[1]
			pingRes, err := strconv.ParseFloat(words[2], 64)
			if err != nil {
				log.Error("Problem when parsing pings")
			}

			if _, ok := s.PingDistances[src]; !ok {
				s.PingDistances[src] = make(map[string]float64)
			}
			if _, ok := s.PingDistances[dst]; !ok {
				s.PingDistances[dst] = make(map[string]float64)
			}

			s.PingDistances[src][dst] += pingRes
			s.PingDistances[dst][src] += pingRes
		}
	}

	s.PingMapMtx.Unlock()

	s.PingAnswerMtx.Lock()
	s.NrPingAnswers++
	s.PingAnswerMtx.Unlock()

	return nil

}

func (s *Service) getPings(readFromFile bool) {
	if !readFromFile {
		if s.Name == Node0 {
			log.Lvl1("Computing new ping distances")
		}

		// measure pings to other nodes
		s.measureOwnPings()
		s.DonePing = true

		s.PingMapMtx.Lock()
		src := s.Nodes.GetServerIdentityToName(s.ServerIdentity())
		for name, dist := range s.OwnPings {
			dst := name

			if _, ok := s.PingDistances[src]; !ok {
				s.PingDistances[src] = make(map[string]float64)
			}
			if _, ok := s.PingDistances[dst]; !ok {
				s.PingDistances[dst] = make(map[string]float64)
			}
			s.PingDistances[src][dst] += dist
			s.PingDistances[dst][src] += dist
			s.PingDistances[src][src] = 0.0
			s.PingDistances[dst][dst] = 0.0
		}
		s.PingMapMtx.Unlock()

		log.LLvl3(s.Nodes.GetServerIdentityToName(s.ServerIdentity()),
			"finished ping own meas with len", len(s.OwnPings))

		// ask for pings from others
		for _, node := range s.Nodes.All {
			if node.Name !=
				s.Nodes.GetServerIdentityToName(s.ServerIdentity()) {

				e := s.SendRaw(node.ServerIdentity, &ReqPings{
					SenderName: s.Nodes.GetServerIdentityToName(
						s.ServerIdentity())})

				if e != nil {
					panic(e)
				}

			}
		}

		// wait for ping replies from everyone but myself
		for s.NrPingAnswers != len(s.Nodes.All)-1 {
			time.Sleep(5 * time.Second)
		}

		// divide all pings by 2
		for i := 0; i < len(s.Nodes.All); i++ {
			name1 := NodeName + strconv.Itoa(i)
			for j := 0; j < len(s.Nodes.All); j++ {
				name2 := NodeName + strconv.Itoa(j)
				s.PingDistances[name1][name2] =
					s.PingDistances[name1][name2] / 2.0
			}
		}

		for i := 0; i < len(s.Nodes.All); i++ {
			name1 := NodeName + strconv.Itoa(i)
			for j := 0; j < len(s.Nodes.All); j++ {
				name2 := NodeName + strconv.Itoa(j)
				if s.PingDistances[name1][name2] !=
					s.PingDistances[name2][name1] {

					log.Lvl1("Error: ping not symmetric")
				}
			}
		}

		// check that there are enough pings
		if len(s.PingDistances) < len(s.Nodes.All) {
			log.Lvl1(s.Nodes.GetServerIdentityToName(s.ServerIdentity()),
				" too few pings 1")
		}
		for _, m := range s.PingDistances {
			if len(m) < len(s.Nodes.All) {
				log.Lvl1(s.Nodes.GetServerIdentityToName(s.ServerIdentity()),
					" too few pings 2")
				log.LLvl1(m)
			}
		}

		log.Lvl3(s.Nodes.GetServerIdentityToName(s.ServerIdentity()),
			"has all pings, starting tree gen")

		/*
			// check TIV just once
			// painful, but run floyd warshall and build the static routes
			// ifrst, let's solve the k means mistery

			if s.Nodes.GetServerIdentityToName(s.ServerIdentity()) == Node0 {
				for n1, m := range s.PingDistances {
					for n2, d := range m {
						//bestDist
						for k := 0; k < len(s.Nodes.All); k++ {
							namek := NodeName + strconv.Itoa(k)
							if d > s.PingDistances[n1][namek]+
								s.PingDistances[namek][n2] {
								log.LLvl1("TIV!", n1, n2, "through", namek,
									"original:", s.PingDistances[n1][n2], ">",
									s.PingDistances[n1][namek], "+",
									s.PingDistances[namek][n2])
								break
							}
						}
					}
				}
			}
		*/

		s.writePingsToFile()

	} else {
		// read from file lines of form "node_19 node_7 : 321"
		readLine := cruxIPFS.ReadFileLineByLine(PingsFile)

		for true {
			line := readLine()
			if line == "" {
				break
			}

			if strings.HasPrefix(line, "#") {
				continue
			}

			tokens := strings.Split(line, " ")
			src := tokens[0]
			dst := tokens[1]
			pingTime, err := strconv.ParseFloat(tokens[3], 64)
			if err != nil {
				log.Error("Problem when parsing pings")
			}

			if _, ok := s.PingDistances[src]; !ok {
				s.PingDistances[src] = make(map[string]float64)
			}

			s.PingDistances[src][dst] = pingTime
		}
	}
}

// measureOwnPings compute ping distances to all nodes from self
func (s *Service) measureOwnPings() {
	myName := s.Nodes.GetServerIdentityToName(s.ServerIdentity())
	for _, node := range s.Nodes.All {

		if node.ServerIdentity.String() != s.ServerIdentity().String() {

			log.Lvl2(myName, "meas ping to ",
				s.Nodes.GetServerIdentityToName(node.ServerIdentity))

			for {

				peerName := node.Name
				pingCmdStr := "ping -W 150 -q -c 3 -i 5 " +
					node.ServerIdentity.Address.Host() + " | tail -n 1"
				pingCmd := exec.Command("sh", "-c", pingCmdStr)
				pingOutput, err := pingCmd.Output()
				if err != nil {
					log.Fatal("couldn't measure ping")
				}

				if strings.Contains(string(pingOutput), "pipe") ||
					len(strings.TrimSpace(string(pingOutput))) == 0 {

					log.Lvl1(s.Nodes.GetServerIdentityToName(
						s.ServerIdentity()), "retrying for", peerName,
						node.ServerIdentity.Address.Host(),
						node.ServerIdentity.String())
					log.LLvl1("retry")
					continue
				}

				processPingCmdStr := "echo " + string(pingOutput) +
					" | cut -d ' ' -f 4 | cut -d '/' -f 1-2 | tr '/' ' '"
				processPingCmd := exec.Command("sh", "-c", processPingCmdStr)
				processPingOutput, _ := processPingCmd.Output()

				strPingOut := string(processPingOutput)

				pingRes := strings.Split(strPingOut, "/")

				avgPing, err := strconv.ParseFloat(pingRes[5], 64)
				if err != nil {
					log.Fatal("Problem when parsing pings")
				}

				s.OwnPings[peerName] = float64(avgPing)
				//s.OwnPings[peerName] = float64(avgPing / 2.0)

				break
			}

		}
	}
}

// printDistances print ping distances in an array
func (s *Service) printDistances(str string) {
	str += "\n       | "
	for i := 0; i < len(s.Nodes.All); i++ {
		str += "  " + NodeName + strconv.Itoa(i) + "  |"
	}
	str += "\n"
	for i := 0; i < len(s.Nodes.All); i++ {
		name1 := NodeName + strconv.Itoa(i)
		str += name1 + " | "
		for j := 0; j < len(s.Nodes.All); j++ {
			name2 := NodeName + strconv.Itoa(j)
			str += fmt.Sprintf(" %f |", s.PingDistances[name1][name2])
		}
		str += "\n"
	}
	log.Lvl1(str)
}

// printPings print ping distances in a list
func (s *Service) printPings() {
	str := "\n"
	for i := 0; i < len(s.Nodes.All); i++ {
		name1 := NodeName + strconv.Itoa(i)
		for j := 0; j < len(s.Nodes.All); j++ {
			name2 := NodeName + strconv.Itoa(j)
			// ping node_1 node_4 = 65.0345
			str += fmt.Sprintf("ping %s %s = %f\n", name1, name2,
				s.PingDistances[name1][name2])
		}
	}
	log.Lvl1(str)
}

// writePingsToFile write ping distances to a file
func (s *Service) writePingsToFile() {
	str := ""
	for i := 0; i < len(s.Nodes.All); i++ {
		name1 := NodeName + strconv.Itoa(i)
		for j := 0; j < len(s.Nodes.All); j++ {
			name2 := NodeName + strconv.Itoa(j)
			// node_5 node_7 : 108.2233
			str += fmt.Sprintf("%s %s : %f\n", name1, name2,
				s.PingDistances[name1][name2])
		}
	}
	ioutil.WriteFile(PingsFile, []byte(str), defaultFileMode)
}
