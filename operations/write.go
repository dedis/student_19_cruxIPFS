package operations

import (
	"context"
	"os"
	"time"

	"github.com/ipfs/ipfs-cluster/api"
	"github.com/ipfs/ipfs-cluster/api/rest/client"
	"go.dedis.ch/onet/v3/log"
)

// WriteFile to the cluster
func writeFile(c client.Client, path string) (string, time.Duration) {
	ctx := context.Background()

	_, err := os.Stat(path)
	if err != nil {
		log.Lvl1(err)
	}
	out := make(chan *api.AddedOutput)
	paths := []string{path}
	start := time.Now()
	go func() {
		err := c.Add(ctx, paths, api.DefaultAddParams(), out)
		if err != nil {
			log.Lvl1(err)
		}
	}()
	ao := <-out
	t := time.Now()
	//fmt.Printf("\nAdded %s: %s\n", filepath.Base(path), ao.Name)
	//fmt.Println(ao.Cid, ao.Size, ao.Bytes)
	if ao == nil {
		log.Lvl1("ao==nil")
		return "", t.Sub(start)
	}
	return ao.Cid.String(), t.Sub(start)
}
