package transaction

import (
	"math/big"
)

type Transaction struct {
	From    [20]byte
	To      [20]byte
	Value   uint64
	Nonce   uint64
	V, R, S *big.Int
}
