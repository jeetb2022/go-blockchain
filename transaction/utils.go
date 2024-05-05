package transaction

import (
	"crypto/ecdsa"
	"errors"
	"fmt"
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"golang.org/x/crypto/sha3"
)

// SignTransaction signs a transaction using the private key of the sender
func SignTransaction(privateKey *ecdsa.PrivateKey, to [20]byte, value, nonce uint64) (*SignedTransaction, error) {
	transaction := Transaction{
		To:    to,
		Value: value,
		Nonce: nonce,
	}

	h := Hash(&transaction)
	signature, err := crypto.Sign(h[:], privateKey)

	if err != nil {
		panic(err)
	}

	r, s, v := decodeSignature(signature)

	signedTransaction := SignedTransaction{
		To:    transaction.To,
		Value: transaction.Value,
		Nonce: transaction.Nonce,
		V:     v,
		R:     r,
		S:     s,
	}
	fmt.Println("Tx hash", HashSigned(&signedTransaction).Hex())

	return &signedTransaction, nil
}

func Hash(tx *Transaction) common.Hash {
	return rlpHash([]interface{}{
		tx.To,
		tx.Value,
		tx.Nonce,
	})
}

func HashSigned(tx *SignedTransaction) common.Hash {
	return rlpHash(tx)
}

func rlpHash(x interface{}) (h common.Hash) {
	sha := hasherPool.Get().(crypto.KeccakState)
	defer hasherPool.Put(sha)
	sha.Reset()
	rlp.Encode(sha, x)
	sha.Read(h[:])

	return h
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

var hasherPool = sync.Pool{
	New: func() interface{} { return sha3.NewLegacyKeccak256() },
}

func recoverPlain(sighash common.Hash, R, S, Vb *big.Int, homestead bool) (common.Address, error) {
	if Vb.BitLen() > 8 {
		// return common.Address{}, ErrInvalidSig
		panic("invalid signature @utils")
	}
	V := byte(Vb.Uint64() - 27)
	if !crypto.ValidateSignatureValues(V, R, S, homestead) {
		// return common.Address{}, ErrInvalidSig
		panic("invalid signature " )
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

func SerializeTransaction(msg SignedTransaction) ([]byte, error) {
	// Encode the message struct
	encodedMsg, err := rlp.EncodeToBytes(msg)
	if err != nil {
		return nil, err
	}
	return encodedMsg, nil
}
func DeserializeTransaction(msg []byte) (SignedTransaction, error) {
	// Decode the message struct
	var transaction SignedTransaction
	err := rlp.DecodeBytes(msg, &transaction)
	if err != nil {
		panic(err)
	}
	return transaction, nil
}
