package txpool

import (
	"Blockchain_Project/transaction"
	"fmt"

	// "sync"

	// "time"

	"github.com/ethereum/go-ethereum/common"
	// "github.com/ethereum/go-ethereum/crypto"
	// "github.com/ethereum/go-ethereum/rlp"
	// "golang.org/x/crypto/sha3"
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
	// if tx == nil {
	// 	return errors.New("transaction is nil")
	// }
	// if tx.To == (common.Address{}) {
	// 	return errors.New("invalid recipient address")
	// }
	// if tx.Value == 0 {
	// 	return errors.New("invalid transaction value")
	// }
	// if tx.Nonce == 0 {
	// 	return errors.New("invalid transaction nonce")
	// }
	// if tx.V == nil || tx.R == nil || tx.S == nil {
	// 	return errors.New("signature components missing")
	// }
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

type Block struct {
	Header       Header
	Transactions []*transaction.SignedTransaction
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

// func rlpHash(x interface{}) (h common.Hash) {
// 	sha := hasherPool.Get().(crypto.KeccakState)
// 	defer hasherPool.Put(sha)
// 	sha.Reset()
// 	rlp.Encode(sha, x)
// 	sha.Read(h[:])

// 	return h
// }

// var hasherPool = sync.Pool{
// 	New: func() interface{} { return sha3.NewLegacyKeccak256() },
// }

// func calculateStateRoot() common.Hash {

// }

// func calculateTransactionsRoot(transactions []*transaction.SignedTransaction) common.Hash {

// }

// func (tp *TransactionPool) CreateBlocks(miner common.Address, extraData []byte) {
// 	blockNumber := uint64(0)
// 	for {
// 		time.Sleep(2 * time.Second)

// 		var transactions []*transaction.SignedTransaction
// 		if len(tp.Transactions) > 10 {
// 			transactions = tp.Transactions[:10]    // taking first 10 transactions
// 			tp.Transactions = tp.Transactions[10:] // removing first 10 transactions from pool
// 		} else {
// 			transactions = tp.Transactions
// 			tp.Transactions = nil
// 		}

// 		// Create the block header
// 		header := Header{
// 			ParentHash:       calculateParentHash(), // Replace with your function
// 			Miner:            miner,
// 			StateRoot:        calculateStateRoot(),                    // Replace with your function
// 			TransactionsRoot: calculateTransactionsRoot(transactions), // Replace with your function
// 			Number:           blockNumber,
// 			Timestamp:        uint64(time.Now().Unix()),
// 			ExtraData:        extraData,
// 		}

// 		block := &Block{
// 			Header:       header,
// 			Transactions: transactions,
// 		}

// 		fmt.Println("Created block with transactions:", block.Transactions)

// 		blockNumber++

// 		if len(tp.Transactions) == 0 {
// 			break
// 		}
// 	}
// }
