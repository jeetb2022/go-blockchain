package utils

import (
	"Blockchain_Project/database"
	"Blockchain_Project/transaction"
	"fmt"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"golang.org/x/crypto/sha3"
)

type Block struct {
	Header       *Header
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

func ValidateAddress(address [20]byte) (bool, error) {
	_, err := database.GetAccountFromDB(address)

	if err != nil {

		fmt.Println("Invalid account address")
		return false, err
	} else {
		fmt.Println("Valid account address")
		return true, nil
	}
}

func StateRoot() common.Hash {
	return HashList(database.GetStateRoot())
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

func CalculateTransactionsRoot(transactions []transaction.SignedTransaction) common.Hash {
	var txHashes []common.Hash
	for _, tx := range transactions {
		txHash := rlpHash(tx)
		txHashes = append(txHashes, txHash)
	}

	transactionsRoot := HashList(txHashes)

	return transactionsRoot
}
