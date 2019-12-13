package template_test

import (
	"testing"

	// We need to include the service so it is started.

	_ "github.com/dedis/student_19_cruxIPFS/service"
	"go.dedis.ch/kyber/v3/suites"
	"go.dedis.ch/onet/v3/log"
)

var tSuite = suites.MustFind("Ed25519")

func TestMain(m *testing.M) {
	log.MainTest(m)
}
