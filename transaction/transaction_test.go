package transaction

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"fmt"
	"testing"
)

func TestSignTransaction(t *testing.T) {
	// Generate sender's private key
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Errorf("Error generating private key: %v", err)
		return
	}
	fmt.Println("Generated sender's private key")

	// Convert the sender's address to [20]byte
	var from [20]byte
	copy(from[:], privateKey.PublicKey.X.Bytes())
	fmt.Println("Converted sender's address to [20]byte")

	// Recipient's address
	var to [20]byte
	fmt.Println("Set recipient's address")

	// Create and sign the transaction
	signedTransaction, err := SignTransaction(privateKey, to, 100, 0)
	if err != nil {
		t.Errorf("Error signing transaction: %v", err)
		return
	}
	fmt.Println("Created and signed the transaction")

	// Verify that the signature components are not nil
	if signedTransaction.V == nil || signedTransaction.R == nil || signedTransaction.S == nil {
		t.Error("Signature components are nil")
		return
	}
	fmt.Println("Signature components are not nil")

	// Optionally, add more assertions to validate the signed transaction
}
