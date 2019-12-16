package main

import (
	"fmt"

	"github.com/dedis/student_19_cruxIPFS/operations"
)

func main() {
	clusters := operations.LoadClusterInstances("../simulation/save.txt")
	fmt.Println(clusters)
}
