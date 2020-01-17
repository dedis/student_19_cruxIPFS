# operations

This folder contains the Crux client interacting with IPFS Cluster. 

## [const.go](const.go)

Contains the constants used by the Crux client.

## [operations.go](operations.go)

Contains the basic commands implementation for the Crux client, such as `Write(f)` or `Read(f)`.

## [performance.go](performance.go)

Contains a few tests, testing interaction latency between pairs of hosts. The first node of the pair using the `Write(f)` function and the second one using the `Read(f)` function. The interaction time of those two calls is measured and outputed.

## [save.go](save.go)

Contains a `Save(state)` and `Load(file)` methods. The `Save(state)` method saves the state of running IPFS and IPFS Cluster daemons, with IP and ports. `Load(file)` reads the saved file and load the running IPFS and IPFS daemon details, in order to interact with them.

## [struct.go](struct.go)

Contains the structures used by the Crux client.
