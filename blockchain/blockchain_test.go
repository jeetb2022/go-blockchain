package block

import (
	"fmt"
	"testing"
)

func TestBlockchain(t *testing.T) {
	// Create a new blockchain
	bc := &Blockchain{}

	// Create some dummy blocks
	block1 := NewBlock([32]byte{}, [20]byte{}, [32]byte{}, [32]byte{}, 0, 0, 1, 0, []byte{}, 0, nil)
	block2 := NewBlock(block1.CalculateHash(), [20]byte{}, [32]byte{}, [32]byte{}, 0, 0, 2, 0, []byte{}, 0, nil)
	block3 := NewBlock(block2.CalculateHash(), [20]byte{}, [32]byte{}, [32]byte{}, 0, 0, 2, 0, []byte{}, 0, nil)

	// Add blocks to the blockchain
	err := bc.AddBlock(block1)
	if err != nil {
		t.Errorf("Failed to add block 1: %v", err)
	} else {
		fmt.Println("Succesfully added block 1")
	}

	err = bc.AddBlock(block2)
	if err != nil {
		t.Errorf("Failed to add block 2: %v", err)
	} else {
		fmt.Println("Succesfully added block 2")
	}

	err = bc.AddBlock(block3)
	if err != nil {
		t.Errorf("Failed to add block 3: %v", err)
	} else {
		fmt.Println("Succesfully added block 3")
	}

	// Verify the integrity of the blockchain
	if len(bc.Blocks) != 3 {
		t.Errorf("Expected blockchain length to be 3, got %d", len(bc.Blocks))
	}
}
