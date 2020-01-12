package main

import (
	"fmt"
	"io/ioutil"
	"math"
	"strconv"
	"strings"
)

const lowBound float64 = 150
const highBound float64 = 300

type pingRes struct {
	Node1 string
	Node2 string
	Ping  float64
}

func main() {
	//folder := "K3N20D121remoteO2000raftnew"
	/*
		folder := "K3N20D121remoteO2000crdt"
		read("c", folder)
		read("v", folder)

		folder = "K3N20D121remoteO2000raft"
		read("c", folder)
		read("v", folder)
	*/

	folder := "K3N20D100remoteO101raft"
	read("c", folder)
	read("v", folder)

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

func read(sim, folder string) {
	filename := folder + "/output_" + sim + ".txt"

	fmt.Println("\n", folder)

	b, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	pings := parsePings(folder)

	lowW := make([]int, 0)
	highW := make([]int, 0)
	lowR := make([]int, 0)
	highR := make([]int, 0)

	write := make([]string, 0)
	read := make([]string, 0)
	writeMax := make([]string, 0)
	readMax := make([]string, 0)

	writeList := make([]int, 0)
	writeMaxList := make([]int, 0)
	readList := make([]int, 0)
	readMaxList := make([]int, 0)

	lines := strings.Split(string(b), "\n")
	for i, l := range lines {
		if strings.Contains(l, "maxoptime") {
			j := 1
			if strings.Contains(lines[i+j], "optime") {
				j++
			}
			nodes := strings.Split(strings.Split(l, "maxoptime")[1], " ")[0]
			writeTimesStr := strings.Split(lines[i+j], " ")
			min := math.MaxInt64
			max := 0
			for _, s := range writeTimesStr[:len(writeTimesStr)-1] {
				tmp, err := strconv.Atoi(s)
				if err != nil {
					fmt.Println(lines[i+1])
					panic(err)
				}
				if tmp < min {
					min = tmp
				}
				if tmp > max {
					max = tmp
				}
			}
			writeList = append(writeList, min)
			writeMaxList = append(writeMaxList, max)
			write = append(write, "writeoptime"+nodes+" "+strconv.Itoa(min))
			writeMax = append(writeMax, "maxwriteoptime"+nodes+" "+strconv.Itoa(max))

			j++
			readTimesStr := strings.Split(lines[i+j], " ")
			min2 := math.MaxInt64
			max2 := 0
			for _, s := range readTimesStr[:len(writeTimesStr)-1] {
				tmp, err := strconv.Atoi(s)
				if err != nil {
					panic(err)
				}
				if tmp < min2 {
					min2 = tmp
				}
				if tmp > max2 {
					max2 = tmp
				}
			}
			readList = append(readList, min2)
			readMaxList = append(readMaxList, max2)
			read = append(read, "readoptime"+nodes+" "+strconv.Itoa(min2))
			readMax = append(readMax, "maxreadoptime"+nodes+" "+strconv.Itoa(max2))

			for _, p := range pings {
				if p.Node1 != p.Node2 && strings.Contains(l, p.Node1) && strings.Contains(l, p.Node2) {
					if p.Ping < lowBound {
						lowR = append(lowR, min2)
						lowW = append(lowW, min)
					} else if p.Ping > highBound {
						highR = append(highR, min2)
						highW = append(highW, min)
					}
					break
				}
			}

		}
	}
	strWrite := ""
	for _, l := range write {
		strWrite += fmt.Sprintln(l)
	}
	strRead := ""
	for _, l := range read {
		strRead += fmt.Sprintln(l)
	}
	strMaxWrite := ""
	for _, l := range writeMax {
		strMaxWrite += fmt.Sprintln(l)
	}
	strMaxRead := ""
	for _, l := range readMax {
		strMaxRead += fmt.Sprintln(l)
	}

	ioutil.WriteFile(folder+"/data/write_"+sim+".txt", []byte(strWrite), 0777)
	ioutil.WriteFile(folder+"/data/read_"+sim+".txt", []byte(strRead), 0777)
	ioutil.WriteFile(folder+"/data/maxwrite_"+sim+".txt", []byte(strMaxWrite), 0777)
	ioutil.WriteFile(folder+"/data/maxread_"+sim+".txt", []byte(strMaxRead), 0777)

	fmt.Println()

	avgR := 0.0
	sd := 0.0
	for _, m := range readList {
		avgR += float64(m)
	}
	avgR /= float64(len(readList))
	for _, m := range readList {
		sd += math.Pow(float64(m)-avgR, 2)
	}
	sd = math.Sqrt(sd / float64(len(readList)))
	fmt.Println("Read average "+sim+":", int64(avgR))
	fmt.Println("Read standard deviation "+sim+":", int64(sd))

	avgW := 0.0
	sd = 0.0
	for _, m := range writeList {
		avgW += float64(m)
	}
	avgW /= float64(len(writeList))
	for _, m := range writeList {
		sd += math.Pow(float64(m)-avgW, 2)
	}
	sd = math.Sqrt(sd / float64(len(writeList)))
	fmt.Println("Write average "+sim+":", int64(avgW))
	fmt.Println("Write standard deviation "+sim+":", int64(sd))

	avg := 0.0
	sd = 0.0
	for _, m := range writeMaxList {
		avg += float64(m)
	}
	avg /= float64(len(writeMaxList))
	for _, m := range writeMaxList {
		sd += math.Pow(float64(m)-avg, 2)
	}
	sd = math.Sqrt(sd / float64(len(writeMaxList)))
	fmt.Println("Max write average "+sim+":", int64(avg))
	fmt.Println("Max write standard deviation "+sim+":", int64(sd))

	avg = 0.0
	sd = 0.0
	for _, m := range readMaxList {
		avg += float64(m)
	}
	avg /= float64(len(readMaxList))
	for _, m := range readMaxList {
		sd += math.Pow(float64(m)-avg, 2)
	}
	sd = math.Sqrt(sd / float64(len(readMaxList)))
	fmt.Println("Max read average "+sim+":", int64(avg))
	fmt.Println("Max read standard deviation "+sim+":", int64(sd))

	if len(lowR) > 0 {
		avgLowR := 0.0
		sd := 0.0
		for _, m := range lowR {
			avgLowR += float64(m)
		}
		avgLowR /= float64(len(lowR))
		for _, m := range lowR {
			sd += math.Pow(float64(m)-avgLowR, 2)
		}
		sd = math.Sqrt(sd / float64(len(lowR)))
		fmt.Println("Read low average "+sim+":", int64(avgLowR))
		fmt.Println("Read low standard deviation "+sim+":", int64(sd))
	}

	if len(highR) > 0 {
		avgHighR := 0.0
		sd := 0.0
		for _, m := range highR {
			avgHighR += float64(m)
		}
		avgHighR /= float64(len(highR))
		for _, m := range highR {
			sd += math.Pow(float64(m)-avgHighR, 2)
		}
		sd = math.Sqrt(sd / float64(len(highR)))
		fmt.Println("Read high average "+sim+":", int64(avgHighR))
		fmt.Println("Read high standard deviation "+sim+":", int64(sd))
	}

	if len(lowW) > 0 {
		avgLowW := 0.0
		sd := 0.0
		for _, m := range lowW {
			avgLowW += float64(m)
		}
		avgLowW /= float64(len(lowW))
		for _, m := range lowW {
			sd += math.Pow(float64(m)-avgLowW, 2)
		}
		sd = math.Sqrt(sd / float64(len(lowW)))
		fmt.Println("Write low average "+sim+":", int64(avgLowW))
		fmt.Println("Write low standard deviation "+sim+":", int64(sd))
	}

	if len(highW) > 0 {
		avgHighW := 0.0
		sd := 0.0
		for _, m := range highW {
			avgHighW += float64(m)
		}
		avgHighW /= float64(len(highW))
		for _, m := range highW {
			sd += math.Pow(float64(m)-avgHighW, 2)
		}
		sd = math.Sqrt(sd / float64(len(highW)))
		fmt.Println("Write high average "+sim+":", int64(avgHighW))
		fmt.Println("Write high standard deviation "+sim+":", int64(sd))
	}

}
