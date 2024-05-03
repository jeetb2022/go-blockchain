package txpool

import (
	"Blockchain_Project/transaction"
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
)

type TransactionPool struct {
	Transactions []*transaction.SignedTransaction
}

// NewTransactionPool creates a new transaction pool
func NewTransactionPool() *TransactionPool {
	return &TransactionPool{
		Transactions: make([]*transaction.SignedTransaction, 0),
	}
}

// AddTransaction adds a new transaction to the pool after validation
func (tp *TransactionPool) AddTransactionToTxPool(tx *transaction.SignedTransaction) error {
	// Basic validation
	if tx == nil {
		return errors.New("transaction is nil")
	}
	if tx.To == (common.Address{}) {
		return errors.New("invalid recipient address")
	}
	if tx.Value == 0 {
		return errors.New("invalid transaction value")
	}
	if tx.Nonce == 0 {
		return errors.New("invalid transaction nonce")
	}
	if tx.V == nil || tx.R == nil || tx.S == nil {
		return errors.New("signature components missing")
	}
	// Add to pool
	tp.Transactions = append(tp.Transactions, tx)
	return nil
}

// PrintAllTransactions prints all transactions in the pool
func (tp *TransactionPool) GetAllTransactions() {
	fmt.Println("Transactions in the pool:")
	for _, tx := range tp.Transactions {
		fmt.Printf("To: %s\n", tx.To.Hex())
		fmt.Printf("Value: %d\n", tx.Value)
		fmt.Printf("Nonce: %d\n", tx.Nonce)
		fmt.Printf("V: %s\n", tx.V.String())
		fmt.Printf("R: %s\n", tx.R.String())
		fmt.Printf("S: %s\n", tx.S.String())
		fmt.Println("---------------------------------")
	}
}
