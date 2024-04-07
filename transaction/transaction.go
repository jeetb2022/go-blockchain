package transaction

import (
	"math/big"
)

type Transaction struct {
	To      [20]byte
	Value   uint64
	Nonce   uint64
	V, R, S *big.Int
}

var tx_pool []Transaction

func AddToPool(tx Transaction) {
	tx_pool = append(tx_pool, tx)
}

func GetPool() []Transaction {
	return tx_pool
}