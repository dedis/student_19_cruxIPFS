module github.com/dedis/student_19_cruxIPFS

require (
	github.com/BurntSushi/toml v0.3.1
	//github.com/dedis/student_19_cruxIPFS v0.0.0-20191031080145-689e36070d7b // indirect
	github.com/stretchr/testify v1.4.0
	github.com/urfave/cli v1.22.0
	go.dedis.ch/cothority/v3 v3.3.1
	go.dedis.ch/kyber/v3 v3.0.7
	go.dedis.ch/onet/v3 v3.0.26
	go.dedis.ch/protobuf v1.0.9
	golang.org/x/sys v0.0.0-20190912141932-bc967efca4b8
	gopkg.in/dedis/onet.v2 v2.0.0-20181115163211-c8f3724038a7 // indirect
	gopkg.in/urfave/cli.v1 v1.20.0
)

replace go.dedis.ch/onet/v3 => /mnt/guillaume/Documents/workspace/go/src/github.com/dedis/onet

go 1.13
