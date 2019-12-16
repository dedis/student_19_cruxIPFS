package operations

import (
	"context"
	"time"

	"github.com/ipfs/ipfs-cluster/api/rest/client"
)

func readFile(c client.Client, filename string) time.Duration {
	ctx := context.Background()

	sh := c.IPFS(ctx)
	start := time.Now()
	_, err := sh.Cat(filename)
	t := time.Now()
	checkErr(err)
	return t.Sub(start)

	/*
		buffer := make([]byte, 1024)
		n, err := rc.Read(buffer)
		checkErr(err)
		str := string(buffer[:n])
		fmt.Printf("\ncat "+filename+":\n%s", str)
	*/
}
