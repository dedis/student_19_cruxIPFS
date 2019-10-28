package main

import (
	"fmt"
	"strings"

	"github.com/dedis/student_19_cruxIPFS/service"
)

func main() {

	filepath := "/home/guillaume/ipfs_test/myfolder/Node0/cluster0/service.json"
	conf, err := service.ReadConfig(filepath)
	if err != nil {
		fmt.Println("here")
		fmt.Println(err)
	}

	iSecret := strings.Index(conf, "secret")
	nLine := strings.Index(conf[iSecret:], "\n")

	conf = strings.ReplaceAll(conf, conf[iSecret:iSecret+nLine], "secret\": \""+"4E1100xDEADBEEF"+"\",")

	iPeername := strings.Index(conf, "peername")
	nLine = strings.Index(conf[iPeername:], "\n")
	conf = strings.ReplaceAll(conf, conf[iPeername:iPeername+nLine], "peername\": \""+"guissou"+"\",")

	err = service.WriteConfig(filepath, conf)
	if err != nil {
		fmt.Println("there")
		fmt.Println(err)
	}
}
