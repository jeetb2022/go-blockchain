package validation

import (
	"Blockchain_Project/blockchain"
	"Blockchain_Project/database"
	"Blockchain_Project/transaction"
	"Blockchain_Project/txpool"
	"Blockchain_Project/utils"
	"errors"
	"fmt"
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"golang.org/x/crypto/sha3"
)

var tp *txpool.TransactionPool

func GetTxPool(p *txpool.TransactionPool) {
	tp = p
}
func ValidateAddress(address [20]byte) (bool) {
	_, err := database.GetAccountFromDB(address)

	if err != nil {

		fmt.Println("Invalid account address")
		return false
	} else {
		fmt.Println("Valid account address")
		return true
	}
}
func CalculateTransactionsRoot(block *blockchain.Block) common.Hash {
	var txHashes []common.Hash
	for _, tx := range block.Transactions {
		txHash := rlpHash(tx)
		txHashes = append(txHashes, txHash)
	}

	transactionsRoot := HashList(txHashes)
	return transactionsRoot
}

func HashList(hashes []common.Hash) common.Hash {
	for len(hashes) > 1 {
		var newHashes []common.Hash
		for i := 0; i < len(hashes); i += 2 {
			if i+1 < len(hashes) {
				combinedHash := append(hashes[i][:], hashes[i+1][:]...)
				newHash := common.BytesToHash(combinedHash)
				newHashes = append(newHashes, newHash)
			} else {
				newHashes = append(newHashes, hashes[i])
			}
		}
		hashes = newHashes
	}
	return hashes[0]
}

func ValidateBlock(block *blockchain.Block) (bool, error) {

	if block.Header == nil {
		return false, errors.New("block's header is nil")
	}

	parentHash, err := database.GetLastBlockHash()
	if err != nil {
		return false, errors.New("failed to get the last block's hash")
	}
	if block.Header.ParentHash != database.RlpHash(parentHash) {
		return false, errors.New("block's parent hash does not match the last block's hash")
	}

	// if block.Header.Miner != bloc {
	// 	return false, errors.New("block's miner address is not the address that signed the block")
	// }

	stateRoot := utils.StateRoot()
	if block.Header.StateRoot != stateRoot {
		return false, errors.New("block's state root does not match the root of the state trie")
	}

	transactionsRoot := CalculateTransactionsRoot(block)
	if block.Header.TransactionsRoot != transactionsRoot {
		return false, errors.New("block's transactions root does not match the root of the transactions trie")
	}

	blockNumber, err := database.GetCurrentHeight()
	if block.Header.Number != blockNumber+1 {
		return false, errors.New("block's number is not the next number in the chain")
	}

	return true, nil
}
func rlpHash(x interface{}) (h common.Hash) {
	sha := hasherPool.Get().(crypto.KeccakState)
	defer hasherPool.Put(sha)
	sha.Reset()
	rlp.Encode(sha, x)
	sha.Read(h[:])

	return h
}

var hasherPool = sync.Pool{
	New: func() interface{} { return sha3.NewLegacyKeccak256() },
}

func Hash(tx *transaction.Transaction) common.Hash {
	return rlpHash([]interface{}{
		tx.To,
		tx.Value,
		tx.Nonce,
	})
}

func recoverPlain(sighash common.Hash, R, S, Vb *big.Int, homestead bool) (common.Address, error) {
	if Vb.BitLen() > 8 {
		// return common.Address{}, ErrInvalidSig
		panic("invalid signature")
	}
	V := byte(Vb.Uint64() - 27)
	if !crypto.ValidateSignatureValues(V, R, S, homestead) {
		// return common.Address{}, ErrInvalidSig
		panic("invalid signature")
	}
	// encode the signature in uncompressed format
	r, s := R.Bytes(), S.Bytes()
	sig := make([]byte, crypto.SignatureLength)
	copy(sig[32-len(r):32], r)
	copy(sig[64-len(s):64], s)
	sig[64] = V
	// recover the public key from the signature
	pub, err := crypto.Ecrecover(sighash[:], sig)
	if err != nil {
		return common.Address{}, err
	}
	if len(pub) == 0 || pub[0] != 4 {
		return common.Address{}, errors.New("invalid public key")
	}
	var addr common.Address
	copy(addr[:], crypto.Keccak256(pub[1:])[12:])
	return addr, nil
}

func ValidateTransaction(trans *transaction.SignedTransaction) bool {
	tx := &transaction.Transaction{
		To:    trans.To,
		Value: trans.Value,
		Nonce: trans.Nonce,
	}
	sender, err := recoverPlain(Hash(tx), trans.R, trans.S, trans.V, true)
	if err != nil {
		return false
	}
	// Check if the sender's account exists in the database
	fmt.Println(sender)
	senderAccount, err := database.GetAccountFromDB(sender)
	if err != nil {
		fmt.Println("Error retrieving sender's account from database:", err)
		return false
	}

	// Check if the sender's account has sufficient balance
	if senderAccount.Balance < tx.Value {
		fmt.Println("Insufficient balance in sender's account")
		return false
	}

	return true
}

func decodeSignature(sig []byte) (r, s, v *big.Int) {
	if len(sig) != crypto.SignatureLength {
		panic(fmt.Sprintf("wrong size for signature: got %d, want %d", len(sig), crypto.SignatureLength))
	}
	r = new(big.Int).SetBytes(sig[:32])
	s = new(big.Int).SetBytes(sig[32:64])
	v = new(big.Int).SetBytes([]byte{sig[64] + 27})
	return r, s, v
}

// func GetMinedBlock() *blockchain.Block {
// 	// Create a new block

// 	block := &blockchain.Block{
// 		Transactions: make([]*transaction.SignedTransaction, 0),
// 	}

// 	// Maximum number of transactions to add to the block
// 	maxTransactions := 10
// 	var pickedTransactions []transaction.SignedTransaction
// 	// Iterate through transactions in the transaction pool
// 	for i := 0; i < len(tp.Transactions) && i < maxTransactions; i++ {
// 		if len(tp.Transactions) <= 0 {
// 			break
// 		}
// 		tx := tp.Transactions[i]

// 		// Validate the transaction
// 		// if ValidateTransaction(tx) {
// 		// 	txToAppend := &transaction.SignedTransaction{
// 		// 		To:    tx.To,
// 		// 		Value: tx.Value,
// 		// 		Nonce: tx.Nonce,
// 		// 		V:     tx.V,
// 		// 		R:     tx.R,
// 		// 		S:     tx.S,
// 		// 	}

// 		// Add the transaction to the block
// 		block.Transactions = append(block.Transactions, tx)
// 		pickedTransactions := append(pickedTransactions, *tx)
// 		pickedTransactions = pickedTransactions

// 		// Remove the transaction from the pool
// 		tp.Transactions = append(tp.Transactions[:i], tp.Transactions[i+1:]...)

// 		// Update the index to handle the removed transaction
// 		i--

// 	}

// 	stateRoot := utils.StateRoot()
// 	transactionRoot := utils.CalculateTransactionsRoot(pickedTransactions)
// 	parentHash, err := database.GetLastBlockHash()
// 	parentHashBytes := database.RlpHash(parentHash)
// 	currentHeight, err := database.GetCurrentHeight()
// 	if err != nil {
// 		panic(err)
// 	}
// 	if err != nil {
// 		panic(err)
// 	}
// 	timestamp := time.Now().Unix()
// 	minerAddr := network.GetMinerAddr()
// 	blockHeader := &blockchain.Header{
// 		ParentHash:       parentHashBytes,
// 		Miner:            minerAddr,
// 		StateRoot:        stateRoot,
// 		TransactionsRoot: transactionRoot,
// 		Number:           currentHeight,
// 		Timestamp:        uint64(timestamp),
// 	}
// 	extradata := blockchain.SignHeader(blockHeader)
// 	blockHeader.ExtraData = extradata[:]
// 	block.Header = blockHeader
// 	return block
// }
