module github.com/dedis/student_19_cruxIPFS

go 1.13

require (
	github.com/BurntSushi/toml v0.3.1
	github.com/dedis/cothority_template v0.0.0-20191121084815-f73b5bf67b5d
	github.com/muesli/clusters v0.0.0-20190807044042-ba9c57dd9228
	github.com/muesli/kmeans v0.0.0-20190917235210-80dfc71e6c5a
	github.com/satori/go.uuid v1.2.0
	github.com/stretchr/testify v1.4.0
	go.dedis.ch/cothority/v3 v3.3.2
	go.dedis.ch/kyber/v3 v3.0.11
	go.dedis.ch/onet/v3 v3.0.31
	gopkg.in/dedis/onet.v1 v1.0.0-20180206090940-2ca76e69d0fc
	gopkg.in/dedis/onet.v2 v2.0.0-20181115163211-c8f3724038a7
	gopkg.in/urfave/cli.v1 v1.20.0
)

//replace go.dedis.ch/onet/v3 => /mnt/guillaume/Documents/workspace/go/src/github.com/dedis/onet
