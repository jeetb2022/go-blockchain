package api

import (
	"Blockchain_Project/account"
	"Blockchain_Project/database"
	"Blockchain_Project/network"
	"Blockchain_Project/transaction"
	"Blockchain_Project/txpool"
	"Blockchain_Project/validation"
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
)

var tp *txpool.TransactionPool

func GetTxPool(p *txpool.TransactionPool) {
	tp = p
}

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
	var tx transaction.SignedTransaction
	var req struct {
		SignedTxn string `json:"signed"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Println(req.SignedTxn)
	encodedSignedTxn := req.SignedTxn
	decodedString, err := hex.DecodeString(encodedSignedTxn)

	if err := rlp.DecodeBytes(decodedString, &tx); err != nil {
		panic(err)
	}
	// Decode the request body into the SignedTransaction struct
	if !validation.ValidateTransaction(&tx) {
		http.Error(w, "Transaction verification failed", http.StatusBadRequest)
		return
	}
	// Add the transaction to the transaction pool
	tp.AddTransactionToTxPool(&tx)

	// Send the transaction over the network
	network.SendTransaction(tx)

	w.WriteHeader(http.StatusOK)
	message := "Transaction added successfully"
	_, err = w.Write([]byte(message))
	if err != nil {
		// Handle error if unable to write to response
		panic(err)
	}
}

func SendUnsignedTxHandler(w http.ResponseWriter, r *http.Request) {
	var tx transaction.SignedTransaction

	err := json.NewDecoder(r.Body).Decode(&tx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	if !validation.ValidateTransaction(&tx) {
		http.Error(w, "Transaction verification failed", http.StatusBadRequest)
		return
	}
	tp.AddTransactionToTxPool(&tx)
	network.SendTransaction(tx)
	// tp.GetAllTransactions()
	w.WriteHeader(http.StatusOK)
	message := "Transaction added scessfully"
	_, err = w.Write([]byte(message))
	if err != nil {
		panic(err)
	}
}

// func SendTxHandler(w http.ResponseWriter, r *http.Request) {
// 	// Create a new signed transaction
// 	tx := transaction.SignedTransaction{
// 		Nonce: 43,
// 		To:    GenerateRandomAddress(),
// 		Value: 1000,
// 		V:     big.NewInt(int64(math_rand.Intn(1000))),
// 		R:     big.NewInt(int64(math_rand.Intn(1000))),
// 		S:     big.NewInt(int64(math_rand.Intn(1000))),
// 	}

// 	// Validate the transaction
// 	// if !validation.ValidateTransaction(&tx) {
// 	// 	http.Error(w, "Transaction verification failed", http.StatusBadRequest)
// 	// 	return
// 	// }

// 	// Add the transaction to the transaction pool
// 	tp.AddTransactionToTxPool(&tx)
// 	type TransactionResponse struct {
// 		Nonce uint64 `json:"nonce"`
// 		To    string `json:"to"`
// 		Value uint64 `json:"value"`
// 		V     int64  `json:"v"`
// 		R     int64  `json:"r"`
// 		S     int64  `json:"s"`
// 	}
// 	// Serialize the transaction data into a response struct
// 	resp := TransactionResponse{
// 		Nonce: tx.Nonce,
// 		To:    tx.To.String(),
// 		Value: tx.Value,
// 		V:     tx.V.Int64(),
// 		R:     tx.R.Int64(),
// 		S:     tx.S.Int64(),
// 	}

// 	// Encode the response object as JSON
// 	w.Header().Set("Content-Type", "application/json") // Set content type to JSON
// 	if err := json.NewEncoder(w).Encode(resp); err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// }

// // Handler function to handle /block endpoint with query parameter "hash"
// func blockByHashHandler(w http.ResponseWriter, r *http.Request) {
// 	// Parse query parameter "hash" to get the block hash
// blockHash := r.URL.Query().Get("hash")

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

func CreateAccountHandler(w http.ResponseWriter, r *http.Request) {
	// Parse request parameters to get Nonce, Balance, and Address values
	nonceStr := r.URL.Query().Get("nonce")
	balanceStr := r.URL.Query().Get("balance")
	address := r.URL.Query().Get("address")

	// Convert Nonce and Balance values to uint64
	nonce, err := strconv.ParseUint(nonceStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid nonce", http.StatusBadRequest)
		return
	}
	balance, err := strconv.ParseUint(balanceStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid balance", http.StatusBadRequest)
		return
	}

	// Generate account struct with provided Nonce, Balance, and Address
	var newAccount *account.Account
	newAccount, err = GenerateAccountWithAddress(address, nonce, balance)

	if err != nil {
		http.Error(w, "Failed to generate account", http.StatusInternalServerError)
		return
	}

	// Add the newly created account to the database
	err = database.AddAccountToDB(newAccount.Address, newAccount)
	if err != nil {
		http.Error(w, "Failed to add account to database", http.StatusInternalServerError)
		return
	}

	// Serialize the account data into JSON
	accountJSON, err := json.Marshal(newAccount)
	if err != nil {
		http.Error(w, "Failed to serialize account", http.StatusInternalServerError)
		return
	}

	// Send response with the account JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(accountJSON)
}

func GenerateAccountWithAddress(address string, nonce, balance uint64) (*account.Account, error) {
	// Convert public address to [20]byte format
	addressBytes := common.HexToAddress(address).Bytes()

	// Create and return the Account struct with provided nonce, balance, and address
	return &account.Account{
		Address: common.BytesToAddress(addressBytes),
		Nonce:   nonce,
		Balance: balance,
	}, nil
}

func GenerateAccount(nonce, balance uint64) (*account.Account, error) {
	// Generate private key and public address with provided Nonce and Balance
	privateKey := account.GeneratePrivAndPubKey()

	// Convert public address to [20]byte format
	addressBytes := common.HexToAddress(privateKey).Bytes()

	// Create and return the Account struct with provided nonce, balance, and address
	return &account.Account{
		Address: common.BytesToAddress(addressBytes),
		Nonce:   nonce,
		Balance: balance,
	}, nil
}

// GenesisBlockHandler is a handler function to add the genesis block into the database
func GenesisBlockHandler(w http.ResponseWriter, r *http.Request) {

	// Create a new genesis block
	database.GenesisBlock()

	// Respond with success message
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Genesis block added successfully")
}

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
// blockHash := r.URL.Query().Get("hash")

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

// GetBalanceHandler handles the /getBalance endpoint with query parameter "address"
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

	// Define the response struct
	type BalanceResponse struct {
		Address string `json:"address"`
		Balance uint64 `json:"balance"`
	}

	// Create the response object
	resp := BalanceResponse{
		Address: addressString,
		Balance: account.Balance,
	}

	// Encode the response object as JSON
	w.Header().Set("Content-Type", "application/json") // Set content type to JSON
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
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

// Handler function to handle /block endpoint with query parameter "number"
func BlockHandler(w http.ResponseWriter, r *http.Request) {
	// Parse query parameter "number" to get the block number
	blockNumberString := r.URL.Query().Get("number")
	blockNumber, err := strconv.ParseUint(blockNumberString, 10, 64)
	if err != nil {
		// Handle invalid block number format
		http.Error(w, "Invalid block number format", http.StatusBadRequest)
		return
	}

	// Get the block from the database
	block, err := database.RetrieveBlockHash(blockNumber)
	if err != nil {
		// Handle error
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Encode the block object as JSON
	w.Header().Set("Content-Type", "application/json") // Set content type to JSON
	if err := json.NewEncoder(w).Encode(block); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// Handler function to handle /block?hash={hash} endpoint
func GetBlockByHashHandler(w http.ResponseWriter, r *http.Request) {
	// Parse query parameter "hash" to get the block hash
	hash := r.URL.Query().Get("hash")
	fmt.Println("Here", hash)
	// Get the block from the database based on the hash
	block, err := database.GetBlockByHash([]byte(hash))
	if err != nil {
		// Handle error (block not found)
		http.Error(w, "Block not found", http.StatusNotFound)
		return
	}

	// Return the block data as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(block)
}

// Handler function to handle /tx?hash={hash} endpoint
func GetTransactionByHashHandler(w http.ResponseWriter, r *http.Request) {
	// Parse query parameter "hash" to get the transaction hash
	hash := r.URL.Query().Get("hash")

	// Get the transaction from the database based on the hash

	tx, err := database.GetSignedTransactionFromLevelDB([]byte(hash))
	if err != nil {
		// Handle error (transaction not found)
		http.Error(w, "Transaction not found", http.StatusNotFound)
		return
	}

	// Return the transaction data as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tx)
}
