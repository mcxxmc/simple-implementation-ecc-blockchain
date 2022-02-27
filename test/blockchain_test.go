package test

import (
	"crypto/rand"
	"github.com/mcxxmc/simple-implementation-ecc-blockchain/blockchain"
	"io"
	"testing"
)

const additionalBlocks = 10
const testDataSize = 128

func TestBlockchain(t *testing.T) {
	bc, err := blockchain.NewBlockchain()
	if err != nil {
		t.Error(err)
		return
	}
	for i := 0; i < additionalBlocks; i ++ {
		data := make([]byte, testDataSize)
		_, err = io.ReadFull(rand.Reader, data)
		if err != nil {
			t.Error("test blockchain #1")
			t.Error(err)
			continue
		}
		nonce, err := blockchain.NewNonce()
		if err != nil {
			t.Error("test blockchain #2")
			t.Error(err)
			continue
		}
		bc.AddNewBlock(data, nonce)
	}
	b, _ := blockchain.Verify(bc)
	if !b {
		t.Error("fail to verify blockchain")
	}
	if bc.Size != 1 + additionalBlocks {
		t.Error("wrong blockchain size; expecting ", 1 + additionalBlocks, ", got ", bc.Size)
	}
	bc.Blocks[2].Content = []byte("")
	b, modified := blockchain.Verify(bc)
	if b || modified != bc.Blocks[2] {
		t.Error("fail to detect change(s) in blockchain")
	}
}
