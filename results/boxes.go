package main

import (
	"fmt"
	"io/ioutil"
	"math"
	"strconv"
	"strings"
)

var box = []float64{0, 50, 100, 150, 200, 250, 300, 1000}

type pingRes struct {
	Node1 string
	Node2 string
	Ping  float64
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	//fmt.Println(" Cruxified IPFS")
	//printBoxes("K3N20D150remoteO2000crdt", "c")
	fmt.Println(" Vanilla IPFS")
	printBoxes("K3N20D150remoteO2000c", "v")

}

func parsePings(folder string) []pingRes {
	filename := folder + "/data/pings.txt"
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	lines := strings.Split(string(b), "\n")
	toRet := make([]pingRes, 0)
	for _, l := range lines {
		tokens := strings.Split(l, " ")
		if len(tokens) > 4 {
			p, err := strconv.ParseFloat(tokens[4], 64)
			if err != nil {
				panic(err)
			}
			toRet = append(toRet, pingRes{
				Node1: tokens[1],
				Node2: tokens[2],
				Ping:  p,
			})
		}
	}
	return toRet
}

func getNumbers(folder, filename string) ([]int, []int, []int, []int) {

	b, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	reads := make([][]int, len(box))
	writes := make([][]int, len(box))

	for i := 0; i < len(box); i++ {
		reads[i] = make([]int, 0)
		writes[i] = make([]int, 0)
	}

	read := make([]string, 0)
	write := make([]string, 0)

	pings := parsePings(folder)

	lines := strings.Split(string(b), "\n")
	for i, l := range lines {
		if strings.Contains(l, "minoptime") {
			if strings.Contains(lines[i+2], "readoptime") {
				nodes := strings.Split(strings.Split(l, "minoptime")[1], " ")[0]

				wr, err := strconv.Atoi(strings.Split(l, " ")[1])
				checkErr(err)
				r, err := strconv.Atoi(strings.Split(lines[i+2], " ")[1])
				checkErr(err)
				w := wr - r
				writes[len(box)-1] = append(writes[len(box)-1], w)
				reads[len(box)-1] = append(reads[len(box)-1], r)

				write = append(write, "writeoptime"+nodes+" "+strconv.Itoa(w))
				read = append(read, "readoptime"+nodes+" "+strconv.Itoa(r))

				for _, p := range pings {
					if p.Node1 != p.Node2 && strings.Contains(l, p.Node1) &&
						strings.Contains(l, p.Node2) {

						for j := 0; j < len(box)-1; j++ {
							if p.Ping > box[j] && p.Ping < box[j+1] {
								writes[j] = append(writes[j], w)
								reads[j] = append(reads[j], r)

								if j == len(box)-2 {
									fmt.Println(l)
									fmt.Println(lines[i+3])
									fmt.Println(lines[i+4])
								}

								break
							}
						}
					}
				}
			}
		}
	}
	wRet := make([]int, len(box))
	rRet := make([]int, len(box))
	wSd := make([]int, len(box))
	rSd := make([]int, len(box))

	for i := 0; i < len(box); i++ {
		sumW := 0
		sumR := 0
		stdevW := 0.0
		stdevR := 0.0
		for j := 0; j < len(writes[i]); j++ {
			sumW += writes[i][j]
			sumR += reads[i][j]
		}
		if len(writes[i]) == 0 {
			wRet[i] = 0
			rRet[i] = 0
		} else {
			wRet[i] = sumW / len(writes[i])
			rRet[i] = sumR / len(reads[i])

			for j := 0; j < len(writes[i]); j++ {
				stdevW += math.Pow(float64(writes[i][j]-wRet[i]), 2)
				stdevR += math.Pow(float64(reads[i][j]-rRet[i]), 2)
			}
			wSd[i] = int(math.Sqrt(stdevW / float64(len(writes[i]))))
			rSd[i] = int(math.Sqrt(stdevR / float64(len(reads[i]))))
		}
	}

	return wRet, wSd, rRet, rSd
}

func printBoxes(folder, sim string) {
	filename := folder + "/output_" + sim + ".txt"

	wAvg, wStdev, rAvg, rStdev := getNumbers(folder, filename)

	sum := make([]int, len(box))
	for i := 0; i < len(wAvg); i++ {
		sum[i] = wAvg[i] + rAvg[i]
	}

	fmt.Println("Write+Read average:", sum)
	fmt.Println("Write average:", wAvg)
	fmt.Println("Write stdev:", wStdev)
	fmt.Println("Read average:", rAvg)
	fmt.Println("Read stdev:", rStdev)
	fmt.Println()
}
