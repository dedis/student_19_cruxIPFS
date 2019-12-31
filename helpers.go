package template

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

// CheckErr checks for an error and prints it
func CheckErr(e error) {
	if e != nil && e != io.EOF {
		fmt.Println(e)
		panic(e)
	}
}

// ReadFileLineByLine reads a file line by line
func ReadFileLineByLine(configFilePath string) func() string {
	f, err := os.Open(configFilePath)
	if err != nil {
		fmt.Println(err, configFilePath)
	}

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
