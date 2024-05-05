package utils

import (
	"Blockchain_Project/database"
	"fmt"
	"sync"
	"Blockchain_Project/blockchain"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"golang.org/x/crypto/sha3"
)

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
	return hashList(database.GetStateRoot())
}

func hashList(hashes []common.Hash) common.Hash {
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

func CalculateTransactionsRoot(block *blockchain.Block) common.Hash {
	var txHashes []common.Hash
	for _, tx := range block.Transactions {
		txHash := rlpHash(tx)
		txHashes = append(txHashes, txHash)
	}

	transactionsRoot := hashList(txHashes)

	return transactionsRoot
}
