module github.com/dedis/student_19_cruxIPFS

go 1.13

require (
	github.com/BurntSushi/toml v0.3.1
	github.com/coreos/etcd v3.3.18+incompatible // indirect
	github.com/dedis/cothority_template v0.0.0-20191121084815-f73b5bf67b5d
	github.com/ipfs/ipfs-cluster v0.11.0
	github.com/multiformats/go-multiaddr v0.0.4
	github.com/satori/go.uuid v1.2.0
	go.dedis.ch/cothority/v3 v3.3.2
	go.dedis.ch/kyber/v3 v3.0.11
	go.dedis.ch/onet/v3 v3.0.31
	go.etcd.io/etcd v3.3.18+incompatible
)

//replace go.dedis.ch/onet/v3 => /mnt/guillaume/Documents/workspace/go/src/github.com/dedis/onet
