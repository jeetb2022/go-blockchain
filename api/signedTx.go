package api

import (
	"Blockchain_Project/transaction"
	"fmt"
	
	"math/big"
	"os"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"golang.org/x/crypto/sha3"
)

// type Tx struct {
// 	To    common.Address
// 	Value uint64
// 	Nonce uint64
// }

// type SignedTx struct {
// 	To      common.Address
// 	Value   uint64
// 	Nonce   uint64
// 	V, R, S *big.Int // signature values
// }

func SignTxn(tx transaction.Transaction) transaction.SignedTransaction {
	// Sign transaction

	// hash of txn
	h := Hash(&tx)
	privatekey, err := crypto.HexToECDSA(os.Getenv("PRIVATE_KEY"))
	if err != nil {
		panic(err)
	}

	sig, err := crypto.Sign(h[:], privatekey)
	if err != nil {
		panic(err)
	}

	R, S, V := decodeSignature(sig)
	signedTx := transaction.SignedTransaction{
		To:    tx.To,
		Value: tx.Value,
		Nonce: tx.Nonce,
		V:     V,
		R:     R,
		S:     S,
	}
	// encodedSignedTxn, err := rlp.EncodeToBytes(signedTx)
	// if err != nil {
	// 	panic(err)
	// }

	// fmt.Println(hex.EncodeToString(encodedSignedTxn))

	return signedTx
}

// HashSigned returns the tx hash
func HashSigned(tx *transaction.SignedTransaction) common.Hash {
	return rlpHash(tx)
}

func Hash(tx *transaction.Transaction) common.Hash {
	return rlpHash([]interface{}{
		tx.To,
		tx.Value,
		tx.Nonce,
	})
}

func rlpHash(x interface{}) (h common.Hash) {
	sha := hasherPool.Get().(crypto.KeccakState)
	defer hasherPool.Put(sha)
	sha.Reset()
	rlp.Encode(sha, x)
	sha.Read(h[:])

	return h
}

// hasherPool holds LegacyKeccak256 hashers for rlpHash.
var hasherPool = sync.Pool{
	New: func() interface{} { return sha3.NewLegacyKeccak256() },
}

// decodeSignature decodes the signature into v, r, and s values
func decodeSignature(sig []byte) (r, s, v *big.Int) {
	if len(sig) != crypto.SignatureLength {
		panic(fmt.Sprintf("wrong size for signature: got %d, want %d", len(sig), crypto.SignatureLength))
	}
	r = new(big.Int).SetBytes(sig[:32])
	s = new(big.Int).SetBytes(sig[32:64])
	v = new(big.Int).SetBytes([]byte{sig[64] + 27})
	return r, s, v
}
