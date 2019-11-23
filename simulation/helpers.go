package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/dedis/student_19_cruxIPFS/gentree"
	"go.dedis.ch/onet/v3/log"
)

// CheckErr checks for an error and prints it
func CheckErr(e error) {
	if e != nil && e != io.EOF {
		fmt.Print(e)
		panic(e)
	}
}

// ReadFileLineByLine reads a file line by line
func ReadFileLineByLine(configFilePath string) func() string {
	wd, err := os.Getwd()
	log.Lvl1(wd)
	f, err := os.Open(configFilePath)
	//defer close(f)
	CheckErr(err)
	reader := bufio.NewReader(f)
	//defer close(reader)
	var line string
	return func() string {
		if err == io.EOF {
			return ""
		}
		line, err = reader.ReadString('\n')
		CheckErr(err)
		line = strings.Split(line, "\n")[0]
		return line
	}
}

// CreateNode with the given parameters
func CreateNode(Name string, x float64, y float64, IP string, level int) *gentree.LocalityNode {
	var myNode gentree.LocalityNode

	myNode.X = x
	myNode.Y = y
	myNode.Name = Name
	myNode.IP[IP] = true
	myNode.Level = level
	myNode.ADist = make([]float64, 0)
	myNode.PDist = make([]string, 0)
	myNode.Cluster = make(map[string]bool)
	myNode.Bunch = make(map[string]bool)
	myNode.Rings = make([]string, 0)
	return &myNode
}
