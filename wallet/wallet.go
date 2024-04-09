package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/hex"
	"fmt"
	"os"
)

func createWallet() {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		panic(err)
	}

	// Serialize the private key to a byte slice
	privateKeyBytes, err := x509.MarshalECPrivateKey(privateKey)
	if err != nil {
		panic(err)
	}

	// Encode the byte slice to a hexadecimal string
	privateKeyStr := hex.EncodeToString(privateKeyBytes)

	// Set the private key as an environment variable
	os.Setenv("PRIVATE_KEY", privateKeyStr)

	fmt.Println("Generated sender's private key")
}

func loadPrivateKey() (*ecdsa.PrivateKey, error) {
	// Retrieve the private key string from the environment variable
	privateKeyStr := os.Getenv("PRIVATE_KEY")
	if privateKeyStr == "" {
		return nil, fmt.Errorf("private key not found in environment variables")
	}

	// Decode the hexadecimal string to a byte slice
	privateKeyBytes, err := hex.DecodeString(privateKeyStr)
	if err != nil {
		return nil, err
	}

	// Deserialize the byte slice to an *ecdsa.PrivateKey object
	privateKey, err := x509.ParseECPrivateKey(privateKeyBytes)
	if err != nil {
		return nil, err
	}

	return privateKey, nil
}
