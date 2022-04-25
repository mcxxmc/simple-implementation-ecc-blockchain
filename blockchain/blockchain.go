package blockchain

import (
	"crypto/rand"
	"io"
	"time"
)

// Blockchain a simple blockchain structure.
//
// Please use NewBlockchain() as the constructor (which gives a genesis block).
type Blockchain struct {
	Blocks 		[]*Block	// where all the blocks are held, sequentially.
	Size		int			// number of blocks in the chain
}

// NextBlockId returns the id of the next block; currently, it assumes that the id of a new block increases by 1 each
// time, so the id of a block is in fact its index in the chain.
func (bc *Blockchain) NextBlockId() int {
	return bc.Size
}

// AddNewBlock adds a new block to the top of the blockchain, given content and nonce.
func (bc *Blockchain) AddNewBlock(content []byte, nonce [32]byte) {
	top := bc.Blocks[bc.Size - 1]
	header := &Header{
		Version:        bc.NextBlockId(),
		Timestamp:      time.Now().String(),
		PreviousHash:   Hash(top),
		MerkleRootHash: "",			// not implemented	// todo
		Nonce:          nonce,
		TargetHash:     "",			// not implemented
	}
	block := &Block{
		PrevBlock: top,
		Header:    header,
		Content:   content,
	}
	bc.Blocks = append(bc.Blocks, block)
	bc.Size ++
}

// Vote votes if the another blockchain is valid and equal to this one.
//
// This process may be simplified but it should work. If any of the previous block is tampered, the final hash of the
// block on top should always be different, even though mathematically, there is a super tiny chance for a collision.
func (bc *Blockchain) Vote(another *Blockchain) bool {
	if bc.Size != another.Size || Hash(bc.Blocks[bc.Size - 1]) != Hash(another.Blocks[another.Size - 1]) {
		return false
	}
	b, _ :=  Verify(another)
	return b
}

// Verify verifies a blockchain. If found, it will return the first invalid block.
//
// Note that it cannot 100% verify the top block; that should be achieved through "voting" across the networks.
func Verify(bc *Blockchain) (bool, *Block) {
	for i := 1; i < bc.Size; i ++ {
		block := bc.Blocks[i]
		if block.PrevBlock != bc.Blocks[i - 1] || block.Header.PreviousHash != Hash(block.PrevBlock) {
			return false, block.PrevBlock
		}
	}
	return true, nil
}

// NewNonce generates a random new nonce.
func NewNonce() ([32]byte, error) {
	tmp := make([]byte, 32)
	var err error
	_, err = io.ReadFull(rand.Reader, tmp)
	var nonce [32]byte
	copy(nonce[:], tmp)
	return nonce, err
}

// NewBlockchain returns the pointer to a new blockchain with a genesis block.
func NewBlockchain() (*Blockchain, error) {
	nonce, err := NewNonce()
	if err != nil {
		return nil, err
	}
	header := &Header{
		Version:        0,
		Timestamp:      time.Now().String(),
		PreviousHash:   "",
		MerkleRootHash: "",
		Nonce:          nonce,
		TargetHash:     "",
	}
	block := &Block{
		PrevBlock: nil,
		Header:    header,
		Content:   []byte("genesis block"),
	}
	return &Blockchain{
		Blocks: []*Block{block},
		Size:   1,
	}, err
}
