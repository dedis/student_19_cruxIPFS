package operations

import (
	"context"
	"os"
	"sync"
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
	cids := make(chan string, 10)
	out := make(chan *api.AddedOutput, 1)
	var wg sync.WaitGroup
	wg.Add(1)
	go func(ch chan string) {
		defer wg.Done()
		for v := range out {
			if v == nil {
				ch <- ""
				return
			}
			ch <- v.Cid.String()
		}
	}(cids)

	paths := []string{path}
	start := time.Now()
	c.Add(ctx, paths, api.DefaultAddParams(), out)
	wg.Wait()
	name := <-cids
	t := time.Now()
	if name == "" {
		log.Lvl1("nil return after write")
	}
	return name, t.Sub(start)
}
