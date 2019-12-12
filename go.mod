module github.com/dedis/student_19_cruxIPFS

go 1.13

require (
	github.com/BurntSushi/toml v0.3.1
	github.com/dedis/cothority_template v0.0.0-20191121084815-f73b5bf67b5d
	github.com/satori/go.uuid v1.2.0
	github.com/stretchr/testify v1.4.0
	go.dedis.ch/cothority/v3 v3.3.2
	go.dedis.ch/kyber/v3 v3.0.9
	go.dedis.ch/onet/v3 v3.0.31
	gopkg.in/urfave/cli.v1 v1.20.0
)

replace go.dedis.ch/onet/v3 => /mnt/guillaume/Documents/workspace/go/src/github.com/dedis/onet
