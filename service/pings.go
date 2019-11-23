package service

import (
	"errors"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"go.dedis.ch/onet/v3/log"
	"go.dedis.ch/onet/v3/network"

	cruxIPFS "github.com/dedis/student_19_cruxIPFS"
)

func (s *Service) ExecReqPings(env *network.Envelope) error {
	// Parse message
	req, ok := env.Msg.(*cruxIPFS.ReqPings)
	if !ok {
		log.Error(s.ServerIdentity(), "failed to cast to ReqPings")
		return errors.New(s.ServerIdentity().String() + " failed to cast to ReqPings")
	}

	// wait for pings to be finished
	for !s.DonePing {
		time.Sleep(5 * time.Second)
	}

	reply := ""
	myName := s.Nodes.GetServerIdentityToName(s.ServerIdentity())
	// build reply
	for peerName, pingTime := range s.OwnPings {
		//if peerName == myName {
		//	reply += myName + " " + peerName + " " + "0.0"
		//} else {
		reply += myName + " " + peerName + " " + fmt.Sprintf("%f", pingTime) + "\n"
		//}
	}

	log.Lvl3("sending", reply)
	requesterIdentity := s.Nodes.GetByName(req.SenderName).ServerIdentity

	e := s.SendRaw(requesterIdentity, &cruxIPFS.ReplyPings{Pings: reply, SenderName: myName})
	if e != nil {
		panic(e)
	}

	return e
}

func (s *Service) ExecReplyPings(env *network.Envelope) error {
	fmt.Println("RepPing")

	// Parse message
	req, ok := env.Msg.(*cruxIPFS.ReplyPings)
	if !ok {
		log.Error(s.ServerIdentity(), "failed to cast to ReplyPings")
		return errors.New(s.ServerIdentity().String() + " failed to cast to ReplyPings")
	}

	// process ping output
	//log.LLvl1("resp=", req.Pings)
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

			//if _, ok := s.PingDistances[dst]; !ok {
			//	s.PingDistances[dst] = make(map[string]float64)
			//}

			s.PingDistances[src][dst] += pingRes
			//s.PingDistances[dst][src] += pingRes
			s.PingDistances[src][src] = 0.0

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
		// measure pings to other nodes
		s.measureOwnPings()
		s.DonePing = true

		s.PingMapMtx.Lock()
		// fill ownping in pingdistances
		for name, dist := range s.OwnPings {
			src := s.Nodes.GetServerIdentityToName(s.ServerIdentity())
			dst := name

			if _, ok := s.PingDistances[src]; !ok {
				s.PingDistances[src] = make(map[string]float64)
			}

			s.PingDistances[src][dst] = dist
			s.PingDistances[src][src] = 0.0
		}
		s.PingMapMtx.Unlock()

		log.LLvl1(s.Nodes.GetServerIdentityToName(s.ServerIdentity()), "finished ping own meas with len", len(s.OwnPings))

		// ask for pings from others
		for _, node := range s.Nodes.All {
			if node.Name != s.Nodes.GetServerIdentityToName(s.ServerIdentity()) {
				fmt.Println("Request ping to", node.ServerIdentity)
				//s.InitRequest(serviceReq)
				e := s.SendRaw(node.ServerIdentity, &cruxIPFS.ReqPings{SenderName: s.Nodes.GetServerIdentityToName(s.ServerIdentity())})
				if e != nil {
					panic(e)
				}
			}
		}

		// wait for ping replies from everyone but myself
		fmt.Println("len(s.Nodes.All)", len(s.Nodes.All))
		for s.NrPingAnswers != len(s.Nodes.All)-1 {
			fmt.Println(s.NrPingAnswers)
			time.Sleep(5 * time.Second)
		}

		// prints
		observerNode := s.Nodes.GetServerIdentityToName(s.ServerIdentity())
		pingDistStr := observerNode + "pingDistStr--------> "

		// divide all pings by 2
		/*
			for i := 0 ; i < 45 ; i++ {
				name1 := "node_" + strconv.Itoa(i)
				for j := 0; j < 45; j++ {
					name2 := "node_" + strconv.Itoa(j)
					s.PingDistances[name1][name2] = s.PingDistances[name1][name2] / 2.0
				}
			}
		*/

		for i := 0; i < 45; i++ {
			name1 := "node_" + strconv.Itoa(i)
			for j := 0; j < 45; j++ {
				name2 := "node_" + strconv.Itoa(j)
				pingDistStr += name1 + "-" + name2 + "=" + fmt.Sprintf("%f", s.PingDistances[name1][name2])
			}
			pingDistStr += "\n"
		}

		log.LLvl1(pingDistStr)

		// check that there are enough pings
		if len(s.PingDistances) < len(s.Nodes.All) {
			log.Lvl1(s.Nodes.GetServerIdentityToName(s.ServerIdentity()), " too few pings 1")
		}
		for _, m := range s.PingDistances {
			if len(m) < len(s.Nodes.All) {
				log.Lvl1(s.Nodes.GetServerIdentityToName(s.ServerIdentity()), " too few pings 2")
				log.LLvl1(m)
			}
		}

		log.LLvl1(s.Nodes.GetServerIdentityToName(s.ServerIdentity()), "has all pings, starting tree gen")

		// check TIV just once
		// painful, but run floyd warshall and build the static routes
		// ifrst, let's solve the k means mistery

		if s.Nodes.GetServerIdentityToName(s.ServerIdentity()) == "node_0" {
			for n1, m := range s.PingDistances {
				for n2, d := range m {
					//bestDist
					for k := 0; k < 20; k++ {
						namek := "node_" + strconv.Itoa(k)
						if d > s.PingDistances[n1][namek]+s.PingDistances[namek][n2] {
							log.LLvl1("TIV!", n1, n2, "through", namek, "original:", s.PingDistances[n1][n2], ">", s.PingDistances[n1][namek], "+", s.PingDistances[namek][n2])
							break
						}
					}
				}
			}
		}

		// ping node_0 node_1 = 19.314
		if s.Nodes.GetServerIdentityToName(s.ServerIdentity()) == "node_0" {
			for n1, m := range s.PingDistances {
				for n2, d := range m {
					log.LLvl1("ping ", n1, n2, "=", d)
				}
			}
		}
	} else {
		// read from file lines of fomrm "ping node_19 node_7 = 32.317"
		//readLine, _ := ReadFileLineByLine("pings10_2.txt")
		readLine := cruxIPFS.ReadFileLineByLine("pings.txt")

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
			//log.LLvl1(tokens)
			pingTime, err := strconv.ParseFloat(tokens[4], 64)
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

func (s *Service) measureOwnPings() {
	myName := s.Nodes.GetServerIdentityToName(s.ServerIdentity())
	for _, node := range s.Nodes.All {

		if node.ServerIdentity.String() != s.ServerIdentity().String() {

			log.Lvl2(myName, "meas ping to ", s.Nodes.GetServerIdentityToName(node.ServerIdentity))

			for {

				peerName := node.Name
				pingCmdStr := "ping -W 150 -q -c 3 -i 1 " + node.ServerIdentity.Address.Host() + " | tail -n 1"
				pingCmd := exec.Command("sh", "-c", pingCmdStr)
				pingOutput, err := pingCmd.Output()
				if err != nil {
					log.Fatal("couldn't measure ping")
				}

				if strings.Contains(string(pingOutput), "pipe") || len(strings.TrimSpace(string(pingOutput))) == 0 {
					log.Lvl1(s.Nodes.GetServerIdentityToName(s.ServerIdentity()), "retrying for", peerName, node.ServerIdentity.Address.Host(), node.ServerIdentity.String())
					log.LLvl1("retry")
					continue
				}

				processPingCmdStr := "echo " + string(pingOutput) + " | cut -d ' ' -f 4 | cut -d '/' -f 1-2 | tr '/' ' '"
				processPingCmd := exec.Command("sh", "-c", processPingCmdStr)
				processPingOutput, _ := processPingCmd.Output()
				if string(pingOutput) == "" {
					//log.Lvl1("empty ping ", myName + " " + peerName)
				} else {
					//log.Lvl1("%%%%%%%%%%%%% ping ", myName + " " + peerName, "output ", string(pingOutput), "processed output ", string(processPingOutput))
				}

				//log.Lvl1("%%%%%%%%%%%%% ping ", s.Nodes.GetServerIdentityToName(s.ServerIdentity()) + " " + peerName, "output ", string(pingOutput), "processed output ", string(processPingOutput))

				strPingOut := string(processPingOutput)

				pingRes := strings.Split(strPingOut, "/")
				//log.LLvl1("pingRes", pingRes)

				avgPing, err := strconv.ParseFloat(pingRes[5], 64)
				if err != nil {
					log.Fatal("Problem when parsing pings")
				}

				s.OwnPings[peerName] = float64(avgPing / 2.0)

				break
			}

		}
	}
}
