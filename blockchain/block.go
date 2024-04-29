package block

import (
	"Blockchain_Project/transaction"
	"Blockchain_Project/utils"
	"crypto/sha256"
)

// CalculateHash calculates the hash of the block header
func (b *Block) CalculateHash() [32]byte {
	headerBytes := utils.Serialize(b.Header)
	hash := sha256.Sum256(headerBytes)
	return hash
}

type Block struct {
	Header       *BlockHeader
	Transactions []*transaction.Transaction
}

type BlockHeader struct {
	ParentHash       [32]byte
	Miner            [20]byte
	StateRoot        [32]byte
	TransactionsRoot [32]byte
	Difficulty       uint64
	TotalDifficulty  uint64
	Number           uint64
	Timestamp        uint64
	ExtraData        []byte
	Nonce            uint64
}

// NewBlock creates a new block
func NewBlock(parentHash [32]byte, miner [20]byte, stateRoot [32]byte, transactionsRoot [32]byte, difficulty uint64, totalDifficulty uint64, number uint64, timestamp uint64, extraData []byte, nonce uint64, transactions []*transaction.Transaction) *Block {
	header := &BlockHeader{
		ParentHash:       parentHash,
		Miner:            miner,
		StateRoot:        stateRoot,
		TransactionsRoot: transactionsRoot,
		Difficulty:       difficulty,
		TotalDifficulty:  totalDifficulty,
		Number:           number,
		Timestamp:        timestamp,
		ExtraData:        extraData,
		Nonce:            nonce,
	}
	return &Block{
		Header:       header,
		Transactions: transactions,
	}
}