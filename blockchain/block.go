package blockchain

import (
	"Blockchain_Project/transaction"
	"github.com/ethereum/go-ethereum/common"
)

type Block struct {
	Header       *Header
	Transactions []*transaction.Transaction
}

type Header struct {
	ParentHash       common.Hash
	Miner            common.Address
	StateRoot        common.Hash
	TransactionsRoot common.Hash
	Number           uint64
	Timestamp        uint64
	ExtraData        []byte
}