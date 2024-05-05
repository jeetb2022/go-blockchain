package api

import (
	"Blockchain_Project/database"
	"Blockchain_Project/network"
	"Blockchain_Project/transaction"
	"Blockchain_Project/txpool"
	"bytes"
	"crypto/rand"
	"fmt"
	"math/big"
	math_rand "math/rand"
	"net/http"

	"github.com/ethereum/go-ethereum/common"
)

var tp *txpool.TransactionPool

func GenerateRandomAddress() common.Address {
	// Generate a random byte slice of length 20 (address size)
	randomBytes := make([]byte, 20)
	_, err := rand.Read(randomBytes)
	if err != nil {
		// Handle error
		panic(err)
	}

	// Convert the byte slice to a common.Address
	randomAddress := common.BytesToAddress(randomBytes)
	return randomAddress
}

// Handler function to handle /sendTx endpoint
func SendTxHandler(w http.ResponseWriter, r *http.Request) {
	tx := transaction.SignedTransaction{
		Nonce: 43,
		To:    GenerateRandomAddress(),
		Value: 1000,
		V:     big.NewInt(int64(math_rand.Intn(1000))),
		R:     big.NewInt(int64(math_rand.Intn(1000))),
		S:     big.NewInt(int64(math_rand.Intn(1000))),
	}
	if !transaction.VerifyTx(tx) {
		http.Error(w, "Transaction verification failed", http.StatusBadRequest)
	}
	tp.AddTransactionToTxPool(&tx)
	// tp.GetAllTransactions()
	w.WriteHeader(http.StatusOK)
	message := "Transaction added scessfully"
	_, err := w.Write([]byte(message))
	if err != nil {
		// Handle error if unable to write to response
		panic(err)
	}

}

// Handler function to handle /blockNumber endpoint
func BlockNumberHandler(w http.ResponseWriter, r *http.Request) {
	// Get the recent most block number from the blockchain
	recentBlockNumber, err := database.GetCurrentHeight()
	if err != nil {
		http.Error(w, "Failed to retrieve recent block number", http.StatusInternalServerError)
		return
	}
	// Send response with the block number
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Recent most block number: %d", recentBlockNumber)
}

// // Handler function to handle /block endpoint with query parameter "hash"
// func blockByHashHandler(w http.ResponseWriter, r *http.Request) {
// 	// Parse query parameter "hash" to get the block hash
// 	blockHash := r.URL.Query().Get("hash")

// 	// Get the block with the specified block hash from the blockchain
// 	block := getBlockByHash(blockHash) // You need to implement this function

// 	// Serialize the block into JSON
// 	blockJSON, err := json.Marshal(block)
// 	if err != nil {
// 		// Handle error
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	// Send response with the block JSON
// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusOK)
// 	w.Write(blockJSON)
// }

// // Handler function to handle /tx endpoint with query parameter "hash"
// func txHandler(w http.ResponseWriter, r *http.Request) {
// 	// Parse query parameter "hash" to get the transaction hash
// 	tx := transaction.SignedTransaction{}
// 	if !transaction.VerifyTx(tx) {
// 		http.Error(w, "Transaction verification failed", http.StatusBadRequest)
// 	}
// 	tp.AddTransactionToTxPool(&tx)
// 	w.WriteHeader(http.StatusOK)
// 	message := "Transaction added sucessfully"
// 	_, err := w.Write([]byte(message))
// 	if err != nil {
// 		// Handle error if unable to write to response
// 		panic(err)
// 	}

// }
func GetTxPool(p *txpool.TransactionPool) {
	tp = p
}

// Handler function to handle /getNonce endpoint with query parameter "address"
func GetNonceHandler(w http.ResponseWriter, r *http.Request) {
	// Parse query parameter "address" to get the account address
	addressString := r.URL.Query().Get("address")
	var address [20]byte
	copy(address[:], addressString)

	// Get the account from the database
	account, err := database.GetAccountFromDB(address)
	if err != nil {
		// Handle error
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Get the nonce from the retrieved account
	nonce := account.Nonce

	// Send response with the nonce
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Nonce of account %s: %d", addressString, nonce)
}

// Handler function to handle /getBalance endpoint with query parameter "address"
func GetBalanceHandler(w http.ResponseWriter, r *http.Request) {
	// Parse query parameter "address" to get the account address
	addressString := r.URL.Query().Get("address")
	var address [20]byte
	copy(address[:], addressString)

	// Get the account from the database
	account, err := database.GetAccountFromDB(address)
	if err != nil {
		// Handle error
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Get the balance from the retrieved account
	balance := account.Balance

	// Send response with the balance
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Balance of account %s: %d", addressString, balance)
}

// Handler function to handle /getknownhost endpoint
func GetKnownHostHandler(w http.ResponseWriter, r *http.Request) {
	
	// Parse query parameter "address" to get the account address
	addrList := network.GetPeerAddrs()
	var buf bytes.Buffer
	for _, addr := range addrList {
		buf.WriteString(addr)
		buf.WriteByte('\n') // Add a newline separator
	}

	// Convert the buffer to a byte slice
	addrBytes := buf.Bytes()

	// message := "Transaction added scessfully"
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte(addrBytes))
	if err != nil {
		// Handle error if unable to write to response
		panic(err)
	}

	// address := r.URL.Query().Get("address")

	// // Your implementation to check if the address is known goes here
	// if address == "0x1234567890" {
	// 	// Send response
	// 	w.WriteHeader(http.StatusOK)
	// 	fmt.Fprintf(w, "Address %s is known", address)
	// 	return
	// }
	// Send response
	// fmt.Fprintf(w, "Address %s is known or unknown", address)
}
