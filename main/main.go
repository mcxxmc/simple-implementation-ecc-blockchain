package main

import (
	"fmt"
	ecc "github.com/mcxxmc/simple-implementation-ecc/ecc"
	"math/rand"
	"time"
)

func main() {
	ep := ecc.SampleElliptic()
	ep.SetGeneratorPoint(15, 13)
	ecdh := ecc.NewInstanceECDH(ep)
	rand.Seed(time.Now().Unix())
	ecdh.RandomlyPicksPrivateKey()
	fmt.Println(ecdh.PublicKey())
}
