package blockchain

import (
	"Blockchain_Project/transaction"
	"io"
	"log"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/joho/godotenv"
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

// func SignHeader(header *Header)  {
// 	PrivateKey, err := crypto.HexToECDSA(os.Getenv("PRIVATE_KEY"))
// 	if err != nil {
// 		panic(err)
// 	}
// 	SgnHdr, err := crypto.Sign(SealHash(header).Bytes(), PrivateKey)
// 	if err != nil {
// 		panic(err)
// 	}

// 	header.ExtraData = SgnHdr

//		return common.Address(header.ExtraData)
//	}
func SignHeader(h Header) []byte {
	// Sign the header
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	pvtKey := os.Getenv("PRIVATE_KEY")
	privateKey, err := crypto.HexToECDSA(pvtKey)
	// fmt.Println("Private Key: ", privateKey)
	if err != nil {
		panic(err)
	}
	sig, err := crypto.Sign(SealHash(&h).Bytes(), privateKey)
	if err != nil {
		panic(err)
	}
	// fmt.Println("Signature: ", sig, "Length: ", len(sig))

	signedBytes, err := rlp.EncodeToBytes(sig)
	if err != nil {
		panic(err)
	}
	// fmt.Println("Signed Bytes: ", signedBytes, "Length: ", len(signedBytes))
	return signedBytes
}

func SealHash(header *Header) (hash common.Hash) {
	hasher := sha3.NewLegacyKeccak256()
	encodeHeader(hasher, header)
	hasher.Sum(hash[:0])

	return hash
}

func encodeHeader(w io.Writer, header *Header) {
	enc := []interface{}{
		header.ParentHash,
		header.Miner,
		header.StateRoot,
		header.TransactionsRoot,
		header.Number,
		header.Timestamp,
	}

	if err := rlp.Encode(w, enc); err != nil {
		panic("can't encode: " + err.Error())
	}
}
