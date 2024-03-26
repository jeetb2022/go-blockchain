package block

import "Blockchain_Project/transaction"

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


