package service

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"go.dedis.ch/onet/v3/log"
)

// StartIPFS starts an IPFS instance for the given service
func (s *Service) StartIPFS() {

	// get config path
	// e.g $GOPATH/src/github.com/dedis/student_19_cruxIPFS/simulation/build
	pwd, err := os.Getwd()
	checkErr(err)
	//pwd := "/users/guissou"
	configPath := filepath.Join(pwd, ConfigsFolder)

	// set the port range allocated to s
	s.MinPort = BaseHostPort + s.getNodeID()*MaxPortNumberPerHost
	s.MaxPort = s.MinPort + MaxPortNumberPerHost

	s.ConfigPath = filepath.Join(configPath, s.Name)

	// create own config home folder and ipfs config folder
	s.MyIPFSPath = filepath.Join(s.ConfigPath, IPFSFolder)
	checkErr(CreateEmptyDir(s.MyIPFSPath))

	/*
		ou, _ := exec.Command("whoami").Output()
		log.Lvl1(string(ou))
		ou, _ = exec.Command("uname", "-a").Output()
		log.Lvl1(string(ou))
		ou, _ = exec.Command("ls", "/bin").Output()
		log.Lvl1(string(ou))
	*/
	ou, err := exec.Command("ipfs", "-c"+s.MyIPFSPath, "init").Output()
	if err != nil {
		log.Lvl1(string(ou))
		log.Error(err)
	}

	// init ipfs in the desired folder
	//	exec.Command("ipfs", "-c"+s.MyIPFSPath, "init").Run()
	/*
		o, err := exec.Command("ipfs", "-c"+s.MyIPFSPath, "init", ">log.out").Output()
		log.Lvl1("ipfs", "-c "+s.MyIPFSPath, "init")
		log.Lvl1(string(o))
		if err != nil {
			log.Error(err)
		}
	*/

	// edit the ip in the config file
	s.EditIPFSConfig()

	// start ipfs daemon
	go func() {
		exec.Command("ipfs", "-c"+s.MyIPFSPath, "daemon").Run()
		fmt.Println("ipfs at ip", s.Name, "crashed")
	}()
	// wait until it has started
	time.Sleep(IPFSStartupTime)
}
