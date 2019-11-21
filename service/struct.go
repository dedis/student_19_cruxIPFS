package service

import (
	"sync"

	"go.dedis.ch/onet/v3"
)

// storageID reflects the data we're storing - we could store more
// than one structure.
var storageID = []byte("main")

// Service is our template-service
type Service struct {
	// We need to embed the ServiceProcessor, so that incoming messages
	// are correctly handled.
	*onet.ServiceProcessor

	storage *storage
}

// storage is used to save our data.
type storage struct {
	Count int
	sync.Mutex
}
