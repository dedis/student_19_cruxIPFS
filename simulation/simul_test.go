package main_test

import (
	"testing"

	"go.dedis.ch/onet/log"
	"go.dedis.ch/onet/simul"
)

func TestMain(m *testing.M) {
	log.MainTest(m)
}

func TestSimulation(t *testing.T) {
	simul.Start("protocol.toml", "service.toml")
}
