package old

import (
	"crypto/sha256"
	"fmt"
)

// Block a block of the blockchain.
type Block struct {
	Prev      *Block
	Next      *Block
	Hash      [32]byte			// 256-bit hash based on the current block and the previous block
	Data      []byte				// can be raw or encrypted
	Timestamp string
	Index     string // assumes to be int in hex form
	OwnerId   string // id of the owner of this block
}

// Copy returns a deep copy of the current block, with Prev and Next set to nil.
func (b *Block) Copy() *Block {
	hash := [32]byte{}
	data := make([]byte, len(b.Data))
	for i, byt := range b.Hash {
		hash[i] = byt
	}
	for i, byt := range b.Data {
		data[i] = byt
	}
	return &Block{
		Hash:      hash,
		Data:      data,
		Timestamp: b.Timestamp,
		Index:     b.Index,
		OwnerId:   b.OwnerId,
	}
}

// Hash calculates and returns a 256-bit hash based on the hash of the previous block and the data, timestamp, Index and
// owner ID of the current block.
func Hash(block *Block) [32]byte {
	tmp := ""
	if block.Prev != nil {
		tmp += string(block.Prev.Hash[:])
	}
	tmp += fmt.Sprintf("%v", block.Data) + block.Timestamp + block.Index + block.OwnerId
	return sha256.Sum256([]byte(tmp))
}
