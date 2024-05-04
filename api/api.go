package api

import (
	"Blockchain_Project/database"
	"Blockchain_Project/transaction"
	"fmt"
	"net/http"
)

// Handler function to handle /sendTx endpoint
func SendTxHandler(w http.ResponseWriter, r *http.Request) {
	// Parse query parameter "signed" to get the signed transaction
	signedTx := r.URL.Query().Get("signed")

	// Decode the signed transaction from base64 or hex or any other encoding if needed

	// Deserialize the signed transaction (RLP decoding)
	var tx transaction.Transaction
	fmt.Println(tx)
	fmt.Println(signedTx)
	// Decode the signed transaction into the tx struct

	// Process the transaction (e.g., add it to the transaction pool, validate, etc.)

	// Send response
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Transaction received and processed successfully")
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
// 	txHash := r.URL.Query().Get("hash")

// 	// Get the transaction with the specified transaction hash from the blockchain

// 	// Serialize the transaction into JSON
// 	txJSON, _ := json.Marshal(transaction)

// 	// Send response with the transaction JSON
// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusOK)
// 	w.Write(txJSON)
// }

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
	address := r.URL.Query().Get("address")

	// Your implementation to check if the address is known goes here
	if address == "0x1234567890" {
		// Send response
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Address %s is known", address)
		return
	}
	// Send response
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Address %s is known or unknown", address)
}