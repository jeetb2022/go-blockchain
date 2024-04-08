package transaction

import (
	"Blockchain_Project/utils"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"math/big"
)

// SignTransaction signs a transaction using the private key of the sender
func SignTransaction(privateKey *ecdsa.PrivateKey, from [20]byte, to [20]byte, value, nonce uint64) (*Transaction, error) {
	transaction := Transaction{
		From:  from,
		To:    to,
		Value: value,
		Nonce: nonce,
	}

	serialized := utils.Serialize(transaction)
	hash := sha256.Sum256(serialized)

	r, s, err := ecdsa.Sign(rand.Reader, privateKey, hash[:])
	if err != nil {
		return nil, err
	}

	return &Transaction{
		From:  from,
		To:    to,
		Value: value,
		Nonce: nonce,
		V:     new(big.Int).SetInt64(int64(27 + privateKey.PublicKey.Curve.Params().N.BitLen() - 8)),
		R:     r,
		S:     s,
	}, nil
}
