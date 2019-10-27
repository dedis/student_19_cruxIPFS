package main

import (
	"fmt"
	"os/exec"
	"time"
)

func main() {

	/*
		//array := []string{"-c", "ipfs -c /home/guillaume/.ipfs_test/myfolder/Node1 daemon"}
		array := []string{"-c", "echo coucou"}

		var procAttr os.ProcAttr
		procAttr.Files = []*os.File{nil, nil, nil}
		p, err := os.StartProcess("/bin/bash", array, &procAttr)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(p.Pid)*/

	go func() {
		o, _ := exec.Command("bash", "-c", "ipfs -c/home/guillaume/.ipfs_test/myfolder/Node1 daemon &").Output()
		fmt.Println(string(o))
	}()
	fmt.Println("Continuing")
	time.Sleep(30 * time.Second)
}
