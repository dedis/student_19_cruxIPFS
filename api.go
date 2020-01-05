package cruxIPFS

/*
The api.go defines the methods that can be called from the outside. Most
of the methods will take a roster so that the service knows which nodes
it should work with.

This part of the service runs on the client or the app.
*/

import (
	"go.dedis.ch/onet/v3"
)

// Client is a structure to communicate with the template
// service
type Client struct {
	*onet.Client
}

/*
// NewClient instantiates a new template.Client
func NewClient() *Client {
	return &Client{Client: onet.NewClient(cothority.Suite, service.ServiceName)}
}
*/
