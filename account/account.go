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

func ValidateAddress(mineraddr [20]byte) bool {
	myString := hex.EncodeToString(mineraddr[:])

	// Get the private key from environment variable or generate it if not set
	privHex := os.Getenv("PRIVATE_KEY")
	if privHex == "" {
		privHex = GeneratePrivAndPubKey()
		// Set the PRIVATE_KEY environment variable
		os.Setenv("PRIVATE_KEY", privHex)
		err := WritePrivateKeyToEnvFile(privHex)
		// Write the private key to the .env file
		if err != nil {
			fmt.Println("Error writing PRIVATE_KEY to .env file:", err)
			return false
		}
	}

	fmt.Println("Address:", myString)

	return true
}

func GenerateRandomHex(length int) (string, error) {
	// Determine the number of bytes needed for the specified length of hex string
	numBytes := length / 2
	if length%2 != 0 {
		numBytes++
	}

	// Generate random bytes
	randomBytes := make([]byte, numBytes)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", err
	}

	// Convert random bytes to hexadecimal string
	hexString := hex.EncodeToString(randomBytes)

	// Truncate the string to the desired length
	if len(hexString) > length {
		hexString = hexString[:length]
	}

	return hexString, nil
}

func GeneratePrivAndPubKey() string {
	privateHex := "ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"
	randomHex, err := GenerateRandomHex(64) // Generate a random hex string of length 32
	if err != nil {
		panic(err)
	}
	fmt.Println("Random Hex:", randomHex)

	privateKey, err := crypto.HexToECDSA(randomHex)
	if err != nil {
		panic(err)
	}

	publicKey := crypto.PubkeyToAddress(privateKey.PublicKey)
	address := publicKey.Hex()
	fmt.Println("Address:", address) // this is the user public address (which can be shared to anyone)

	return privateHex
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
