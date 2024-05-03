package block

import (
	"errors"
)

// Blockchain represents the blockchain
type Blockchain struct {
	Blocks []*Block
}

// AddBlock adds a new block to the chain
func (bc *Blockchain) AddBlock(newBlock *Block) error {
	// Verify the integrity of the new block
	if !bc.verifyBlock(newBlock) {
		return errors.New("block verification failed")
	}

	// Append the new block to the chain
	bc.Blocks = append(bc.Blocks, newBlock)
	return nil
}

// Verify the integrity of the block
func (bc *Blockchain) verifyBlock(newBlock *Block) bool {
	// Check if the blockchain is empty
	if len(bc.Blocks) == 0 {
		return true
	}
	// Get the hash of the last block in the chain
	lastBlock := bc.Blocks[len(bc.Blocks)-1]
	lastBlockHash := lastBlock.CalculateHash()

	// Verify that the parent hash of the new block matches the hash of the last block
	return newBlock.Header.ParentHash == lastBlockHash 
}
