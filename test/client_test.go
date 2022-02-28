package test

import (
	"crypto/rand"
	"github.com/mcxxmc/simple-implementation-ecc-blockchain/client"
	"github.com/mcxxmc/simple-implementation-ecc/ecc"
	"io"
	"testing"
)

const testClientLoop = 100

func TestClientBasic(t *testing.T) {
	// test for communications between 2 clients
	ep := ecc.NewElliptic(1, 1, 23)
	ep.SetGeneratorPoint(3, 10)
	alice := client.NewClient(ep, "alice", 1)
	bob := client.NewClient(ep, "bob", 2)
	alice.RandomInitialization()
	bob.RandomInitialization()

	for i := 0; i < testClientLoop; i ++ {
		bytes := make([]byte, 512)
		_, err := io.ReadFull(rand.Reader, bytes)
		if err != nil {
			t.Error(err)
			continue
		}
		msg := string(bytes)

		// from alice to bob
		encrypted, err := alice.EncryptMsg(msg)
		if err != nil {
			t.Error(err)
			continue
		}
		decrypted, err := bob.DecryptMsg(encrypted)
		if err != nil {
			t.Error(err)
			continue
		}
		if decrypted != msg {
			t.Error("wrong msg!")
		}

		// from bob to alice
		encrypted, err = bob.EncryptMsg(msg)
		if err != nil {
			t.Error(err)
			continue
		}
		decrypted, err = alice.DecryptMsg(encrypted)
		if err != nil {
			t.Error(err)
			continue
		}
		if decrypted != msg {
			t.Error("wrong msg!")
		}
	}
}
