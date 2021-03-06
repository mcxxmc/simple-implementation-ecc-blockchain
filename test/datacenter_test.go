package test

import (
	"fmt"
	"github.com/mcxxmc/simple-implementation-ecc-blockchain/client"
	"github.com/mcxxmc/simple-implementation-ecc-blockchain/datacenter"
	"github.com/mcxxmc/simple-implementation-ecc/ecc"
	"github.com/mcxxmc/simple-implementation-ecc/galois"
	"testing"
)

func TestDatacenter(t *testing.T) {
	ep := ecc.NewElliptic(1, 1, 23)
	ep.SetGeneratorPoint(3, 10)
	dc, err := datacenter.NewDataCenter(ep, "datacenter", 0)	// assumes id 0 - 100 is reserved for datacenters
	if err != nil {
		t.Fatal(err)
	}
	err = dc.RandomInitialization()
	if err != nil {
		t.Fatal(err)
	}
	dcPubKey, err := dc.RequestPublicKey(0)
	if err != nil {
		t.Fatal(err)
	}
	if !galois.PointEqual(dcPubKey, dc.Self.GetPublicKey()) {
		t.Fatal("datacenter pubKey does not match")
	}

	alice := client.NewClient(ep, "alice", 101)
	bob := client.NewClient(ep, "bob", 102)
	alice.RandomInitialization()
	bob.RandomInitialization()
	for dc.PublicKeyExists(alice.GetPublicKey()) {
		alice.RandomInitialization()
	}
	for dc.PublicKeyExists(bob.GetPublicKey()) {
		bob.RandomInitialization()
	}
	fmt.Println("datacenter private key: ", dc.Self.Ecdh.PrivateKey, ", alice private key: ", alice.Ecdh.PrivateKey,
		", bob private key: ", bob.Ecdh.PrivateKey)

	err = dc.RequestRegisterNewUser(101, alice.GetPublicKey())
	if err != nil {
		t.Fatal(err)
	}
	err = dc.RequestRegisterNewUser(102, bob.GetPublicKey())
	if err != nil {
		t.Fatal(err)
	}

	alicePubKey, err := dc.RequestPublicKey(101)
	if !galois.PointEqual(alicePubKey, *alice.PubKey) {
		t.Fatal("alice key does not match")
	}
	bobPubKey, err := dc.RequestPublicKey(102)
	if !galois.PointEqual(bobPubKey, *bob.PubKey) {
		t.Fatal("bob key does not match")
	}
	voidPubKey, err := dc.RequestPublicKey(103)
	if !voidPubKey.IsNone || err == nil {
		t.Fatal("this key should not exist")
	}

	aliceOldPrivKey := alice.Ecdh.PrivateKey
	for alice.Ecdh.PrivateKey == aliceOldPrivKey || dc.PublicKeyExists(alice.GetPublicKey()) {
		alice.RandomInitialization()	// set a new private key
	}
	oldSharedKey := ecc.Calculate(dcPubKey, aliceOldPrivKey, ep)
	err = dc.RequestUpdateUserPubKey(101, oldSharedKey, alice.GetPublicKey())
	if err != nil {
		t.Fatal(err)
	}
	// try the same again
	err = dc.RequestUpdateUserPubKey(101, oldSharedKey, alice.GetPublicKey())
	if err == nil {
		t.Fatal("1. should not be able to update")
	}
	// again and again
	newSharedKey := ecc.Calculate(dcPubKey, alice.Ecdh.PrivateKey, ep)
	err = dc.RequestUpdateUserPubKey(101, newSharedKey, alice.GetPublicKey())
	if err == nil {
		t.Fatal("2. should not be able to update")
	}
	err = dc.RequestUpdateUserPubKey(101, newSharedKey, bob.GetPublicKey())
	if err == nil {
		t.Fatal("3. should not be able to update")
	}

	err = dc.ForceWrite()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("blockchain size: ", dc.Bc.Size)
	for i, block := range dc.Bc.Blocks {
		fmt.Println("version number: ", i)
		fmt.Println(string(block.Content))
	}
}
