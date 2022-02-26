package main

import (
	"fmt"
	ecc "github.com/mcxxmc/simple-implementation-ecc/ecc"
)

func main() {
	ep := ecc.SampleElliptic()
	ep.SetGeneratorPoint(15, 13)
	ecdh := ecc.NewInstanceECDH(ep)
	ecdh.RandomlyPicksPrivateKey()
	fmt.Println(ecdh.PublicKey())
}
