package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"sync"
)

func main() {

	os.RemoveAll("/home/guillaume/cluster0")
	os.RemoveAll("/home/guillaume/cluster1")

	os.Mkdir("/home/guillaume/cluster0", 0777)
	os.Mkdir("/home/guillaume/cluster1", 0777)

	var wg sync.WaitGroup
	for i := 0; i < 2; i++ {
		wg.Add(1)
		go func(i string) {
			cmd := "CLUSTER_LISTENMULTIADDRESS=\"/ip4/0.0.0.0/tcp/999" + i + "\" ipfs-cluster-service -c /home/guillaume/cluster" + i + " init"
			//os.Setenv("CLUSTER_LISTENMULTIADDRESS", "/ip4/0.0.0.0/tcp/999"+i)
			//o, err := exec.Command("ipfs-cluster-service", "-c", "/home/guillaume/cluster"+i, "init").Output()
			o, err := exec.Command("bash", "-c", cmd).Output()
			fmt.Println(string(o))
			if err != nil {
				fmt.Println(err)
			}
			wg.Done()

		}(strconv.Itoa(i))
	}
	wg.Wait()

}
