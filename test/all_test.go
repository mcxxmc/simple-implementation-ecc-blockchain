package test

import "testing"

func TestAll(t *testing.T) {
	TestBlockchain(t)
	TestClientBasic(t)
	TestDatacenter(t)
}
