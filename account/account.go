package account

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"

	"github.com/ethereum/go-ethereum/crypto"
)

type Account struct {
	Address [20]byte
	Nonce   uint64
	Balance uint64
}


func GeneratePrivAndPubKey() string {
	randomHex, err := GenerateRandomHex(64) // Generate a random hex string of length 32
	if err != nil {
		panic(err)
	}
	// fmt.Println("Random Hex:", randomHex)
	privateKey, err := crypto.HexToECDSA(randomHex)
	if err != nil {
		panic(err)
	}

	publicKey := crypto.PubkeyToAddress(privateKey.PublicKey)
	address := publicKey.Hex()
	os.Setenv("PRIVATE_KEY",randomHex)
	
	if err:=WritePrivateKeyToEnvFile(randomHex); err!=nil{
		panic(err)
	}
	fmt.Println("Address:", address) // this is the user public address (which can be shared to anyone)

	return address
}


func GenerateRandomHex(length int) (string, error) {
	// Determine the number of bytes needed
	byteLength := length / 2
	if length%2 != 0 {
		byteLength++
	}

	// Generate random bytes
	bytes := make([]byte, byteLength)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}

	// Convert bytes to hexadecimal string
	hexString := hex.EncodeToString(bytes)

	// Truncate to desired length
	hexString = hexString[:length]

	return hexString, nil
}


func WritePrivateKeyToEnvFile(privateKey string) error {
	envFile := ".env"
	file, err := os.OpenFile(envFile, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = fmt.Fprintf(file, "PRIVATE_KEY=%s\n", privateKey)
	if err != nil {
		return err
	}

	return nil
}
