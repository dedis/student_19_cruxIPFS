package operations

//"github.com/ipfs/ipfs-cluster/api/rest/client"
//ma "github.com/multiformats/go-multiaddr"

import (
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	cruxIPFS "github.com/dedis/student_19_cruxIPFS"
	"github.com/ipfs/ipfs-cluster/api/rest/client"
	"go.dedis.ch/onet/v3/log"
)

var nodes = make(map[string]*Node)

// Read the given filename from the given node
func Read(node, filename string) time.Duration {
	if len(nodes) == 0 {
		nodes = LoadClusterInstances(cruxIPFS.SaveFile)
	}
	if n, ok := nodes[node]; ok {
		var min time.Duration = math.MaxInt64
		//var max time.Duration
		mutex := &sync.Mutex{}
		wg := sync.WaitGroup{}
		wg.Add(len(n.Clients))

		for _, c := range n.Clients {
			go func(c0 client.Client, m *sync.Mutex) {
				t := readFile(c0, filename)
				m.Lock()
				if t < min {
					min = t
				}
				m.Unlock()
				log.Lvl1("time", t)
				wg.Done()
			}(c, mutex)
		}
		wg.Wait()
		return min
	}
	panic(node + "do not exist")
}

// Write the given filename from the given node
func Write(node, filename string) (string, time.Duration) {
	if len(nodes) == 0 {
		nodes = LoadClusterInstances(cruxIPFS.SaveFile)
	}
	if n, ok := nodes[node]; ok {
		var sum time.Duration
		mutex := &sync.Mutex{}
		wg := sync.WaitGroup{}
		wg.Add(len(n.Clients))
		name := ""

		for _, c := range n.Clients {
			go func(c0 client.Client, m *sync.Mutex) {
				n, t := writeFile(c0, filepath.Join(fileFolder, filename))
				m.Lock()
				sum += t
				name = n
				m.Unlock()
				log.Lvl1("time", t)

				wg.Done()
			}(c, mutex)
		}
		wg.Wait()
		return name, sum
	}
	panic(node + "do not exist")
}

// NewFile write new file to disk
func NewFile(filename string) {
	os.Mkdir(fileFolder, defaultFileMode)
	str := strings.Repeat("abcd", 256)
	ioutil.WriteFile(filepath.Join(fileFolder, filename), []byte(str),
		defaultFileMode)
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func printNodes() {
	str := ""
	for _, n := range nodes {
		str += n.Name + " : "
		for _, c := range n.Clients {
			str += fmt.Sprint(c, ", ")
		}
	}
	log.Lvl1(str)
}
