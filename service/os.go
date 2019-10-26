package service

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// default file mode for files that the program writes to the system
const defaultFileMode os.FileMode = 0777

// CreateEmptyDir create an empty directory at the given path
func CreateEmptyDir(path string) error {
	// remove existing dir if any
	os.RemoveAll(path)
	// create the empty directorys
	err := os.Mkdir(path, defaultFileMode)
	return err
}

// GetNextAvailablePort return n available ports between pmin and pmax
// crash if error
func GetNextAvailablePort(pmin, pmax, n int) []int {
	cmd := "comm -23 <(seq \"" + strconv.Itoa(pmin) + "\" \"" +
		strconv.Itoa(pmax) + "\" | sort) <(ss -tan | awk '{print $4}' | cut" +
		"-d':' -f2 | grep '[0-9]\\{1,5\\}' | sort -u) | head -n " +
		strconv.Itoa(n)
	out, err := exec.Command("bash", "-c", cmd).Output()
	// if command fails, crash the program
	if err != nil {
		fmt.Println(err)
		fmt.Println(out)
		os.Exit(1)
	}
	ret := make([]int, 0)
	for i, s := range strings.Split(string(out), "\n") {
		if i >= n {
			break
		}
		p, err := strconv.Atoi(s)
		// if error while reading the ports, crash the program
		if err != nil {
			fmt.Println(out)
			fmt.Println(err)
			fmt.Println("Fatal: cannot parse available ports")
			os.Exit(1)
		}
		ret = append(ret, p)
	}
	return ret
}
