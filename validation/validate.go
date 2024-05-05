package validation

import (
	"Blockchain_Project/blockchain"
	"Blockchain_Project/database"
	"Blockchain_Project/transaction"
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
