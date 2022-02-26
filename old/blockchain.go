package old

import (
	"bytes"
	"math/big"
	"time"
)

// Blockchain a simple blockchain structure, with .Top pointing to the latest block.
type Blockchain struct {
	Top 		*Block
}

// generates the id for the new block, which is assumed to be an int in hex form.
//
// The id of the first block is always assumed to be 0.
func (bc *Blockchain) nextBlockId() string {
	if bc.Top == nil {
		return big.NewInt(0).Text(16)
	}
	i := new(big.Int)
	i.SetString(bc.Top.Index, 16)
	return i.Add(i, big.NewInt(1)).Text(16)
}

// AddNewRecord adds a new record to the top of the blockchain. The data can be either raw or encrypted.
func (bc *Blockchain) AddNewRecord(data []byte, ownerId string) {
	block := &Block{
		Prev:      bc.Top,
		Data:      data,
		Timestamp: time.Now().String(),
		Index:     bc.nextBlockId(),
		OwnerId:   ownerId,
	}
	block.Hash = Hash(block)
	bc.Top.Next = block
	bc.Top = block
}

// Copy returns a deep copy of the current blockchain.
func (bc *Blockchain) Copy() *Blockchain {
	bc2 := &Blockchain{}
	pointer := bc.Top
	for pointer.Prev != nil {	// go to the bottom
		pointer = pointer.Prev
	}
	var prev, cur *Block
	for pointer != nil {
		prev = cur
		cur = pointer.Copy()
		cur.Prev = prev
		if prev != nil {
			prev.Next = cur
		}
		pointer = pointer.Next
	}
	bc2.Top = cur
	return bc2
}

// ValidateBlockchain validates a blockchain by checking the hash of each block;
//
// If a block is found to be "suspicious", this function will return false and returns a pointer to that block.
func ValidateBlockchain(bc *Blockchain) (bool, *Block) {
	pointer := bc.Top
	for pointer != nil {
		hash := Hash(pointer)
		if bytes.Equal(pointer.Hash[:], hash[:]) {
			return false, pointer
		}
	}
	return true, nil
}
