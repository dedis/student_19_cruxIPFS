package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func main() {
	pmin := 14000
	pmax := 15000
	n := 3
	cmd := "comm -23 <(seq \"" + strconv.Itoa(pmin) + "\" \"" + strconv.Itoa(pmax) +
		"\" | sort) <(ss -tan | awk '{print $4}' | cut -d':' -f2 | " +
		"grep '[0-9]\\{1,5\\}' | sort -u) | head -n " + strconv.Itoa(n)
	out, err := exec.Command("bash", "-c", cmd).Output()
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
		if err != nil {
			fmt.Println(err)
			return
		}
		ret = append(ret, p)
	}
	fmt.Println(ret)
}
