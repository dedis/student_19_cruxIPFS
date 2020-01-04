package service

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// default file mode for files that the program writes to the system
const defaultFileMode os.FileMode = 0777

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

// CreateEmptyDir create an empty directory at the given path
func CreateEmptyDir(path string) error {
	// remove existing dir if any
	err := os.RemoveAll(path)
	if err != nil {
		return err
	}

	// create the empty directorys
	err = os.MkdirAll(path, defaultFileMode)
	return err
}

// GetNextAvailablePorts return n available ports between pmin and pmax
// crash if error
func GetNextAvailablePorts(pmin, pmax, n int) (*[]int, error) {
	cmd := "comm -23 <(seq \"" + strconv.Itoa(pmin) + "\" \"" +
		strconv.Itoa(pmax) + "\" | sort) <(ss -tan | awk '{print $4}' | cut " +
		"-d':' -f2 | grep '[0-9]\\{1,5\\}' | sort -u) | head -n " +
		strconv.Itoa(n)
	out, err := exec.Command("bash", "-c", cmd).Output()
	// if command fails, crash the program
	if err != nil {
		return nil, err
	}
	ret := make([]int, 0)
	for i, s := range strings.Split(string(out), "\n") {
		if i >= n {
			break
		}
		p, err := strconv.Atoi(s)
		// if error while reading the ports, crash the program
		if err != nil {
			return nil, errors.New("cannot parse available ports")
		}
		ret = append(ret, p)
	}
	if len(ret) != n {
		return nil, errors.New("Couldn't get" + strconv.Itoa(n) +
			"available ports")
	}
	return &ret, nil
}

// ReadConfig read a config file given as parameter and returns a string
func ReadConfig(file string) (string, error) {
	dat, err := ioutil.ReadFile(file)
	if err != nil {
		return "", err
	}
	return string(dat), nil
}

// WriteConfig write string to a file
func WriteConfig(path string, config string) error {
	// Open file using WRITE only permission.
	file, err := os.OpenFile(path, os.O_WRONLY, defaultFileMode)
	if err != nil {
		return err
	}
	defer file.Close()

	// Write some text line-by-line to file.
	_, err = file.WriteString(config)
	if err != nil {
		return err
	}

	// Save file changes.
	err = file.Sync()
	if err != nil {
		return err
	}
	return nil
}

func genSecret() string {
	// generate random secret
	key := make([]byte, 32)
	_, err := rand.Read(key)
	checkErr(err)
	return hex.EncodeToString(key)
}
