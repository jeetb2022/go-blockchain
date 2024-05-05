package validation

import (
	"Blockchain_Project/blockchain"
	"Blockchain_Project/database"
	"Blockchain_Project/utils"
	"errors"
	"fmt"
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

	if block.Header.Miner != blockchain.SignHeader(block.Header) {
		return false, errors.New("block's miner address is not the address that signed the block")
	}

	stateRoot := utils.StateRoot()
	if block.Header.StateRoot != stateRoot {
		return false, errors.New("block's state root does not match the root of the state trie")
	}

	transactionsRoot := utils.CalculateTransactionsRoot(block)
	if block.Header.TransactionsRoot != transactionsRoot {
		return false, errors.New("block's transactions root does not match the root of the transactions trie")
	}

	blockNumber, err := database.GetCurrentHeight()
	if block.Header.Number != blockNumber+1 {
		return false, errors.New("block's number is not the next number in the chain")
	}

	return true, nil
}
