package transaction

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type Transaction struct {
	To    common.Address
	Value uint64
	Nonce uint64
}
type SignedTransaction struct {
	To      common.Address
	Value   uint64
	Nonce   uint64
	V, R, S *big.Int
}

func VerifyTx(tx SignedTransaction) bool {
	status := true
	// Apply Verify Logic
	
	return status
}
