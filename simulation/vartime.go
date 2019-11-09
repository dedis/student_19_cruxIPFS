// +build vartime

package main

import (
	"go.dedis.ch/cothority"
	"go.dedis.ch/kyber"
	"go.dedis.ch/kyber/pairing"
	"go.dedis.ch/kyber/pairing/bn256"
)

func init() {
	cothority.Suite = struct {
		pairing.Suite
		kyber.Group
	}{
		Suite: bn256.NewSuite(),
		Group: bn256.NewSuiteG2(),
	}
}
