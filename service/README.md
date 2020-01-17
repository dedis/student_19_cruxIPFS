# service

This folder contains all the services used to start IPFS, configure it and interact with Onet.

## [clusterbootstrap.go](clusterbootstrap.go)

This file contains the `ClusterBootstrapProtocol`. This protocol is called is called by the protocol `StartIPFSProtocol`, and is run for each ARA on the leader, this one will start an IPFS Cluster instance, and communicate its bootstrap address to all the other members of the ARA. The other members will start an IPFS Cluster instance with the bootstrap address of the leader to join the ARA.

## [const.go](const.go)

This file contains the constants used in the [service](.) folder. 

## [helpers.go](helpers.go)

This file contains helpers methods used in the [service](.) folder. Most of the methods directly interact with the os, for instance creating directory or getting an unused port.

## [ipfs.go](ipfs.go)

This file contains all the code setting up IPFS and IPFS Cluster daemons, and starting them.

## [ping.go](ping.go)

This file contains the code used to compute the ping distances between the hosts. The ping distances can be computed: each host ping every other host in the system and then, they share their distances to all other hosts with all peers in the system. The ping distance between each pair of hosts can also be loaded from a text file.

## [service.go](service.go)

This file contains the Onet service. The `Setup` service initializes all the data structures, computes ping distances and builds the ARAs.

## [startara.go](startara.go)

This file contains the protocol `StartARAProtocol`. This protocol is called by `StartInstancesProtocol` on the root of an ARA. It starts an IPFS daemon, then an IPFS Cluster daemon on the ARA leader, and then broadcast its bootstrap address to all other cluster members that will join the cluster using this bootstrap address. The other cluster members will also start one new IPFS and one new IPFS Cluster instance to join the ARA.

## [startinstances.go](startinstances.go)

This file contains the protocol `StartInstancesProtocol`. This protocol is called on the root of the Onet tree, and for each ARA calls the protocol `StartARAProtocol`, in order to start one IPFS and IPFS Cluster daemon for each ARA membership on each node. 

## [startipfs.go](startipfs.go)

This file contains the protocol `StartIPFSProtocol`. This protocol is run on the root of the Onet tree, and starts a single IPFS daemon on each host. Each host return its IPFS bootstrap address to the leader. Then this protocol start an instance of `ClusterBootstrapProtocol` for each ARA on the ARA leader.

## [struct.go](struct.go)

This file contains the structures used in the [service](.) folder. 
