package main

import (
	"Blockchain_Project/api"
	"Blockchain_Project/cli"
	"Blockchain_Project/network"
	"Blockchain_Project/transaction"
	"Blockchain_Project/txpool"
	"Blockchain_Project/validation"
	"crypto/rand"
	"fmt"
	"math/big"
	"net/http"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/joho/godotenv"
)

var tp = txpool.NewTransactionPool()

func randomBigInt() *big.Int {
	var numBytes [32]byte
	rand.Read(numBytes[:])
	return new(big.Int).SetBytes(numBytes[:])
}
func main() {
	// Load environment variables
	godotenv.Load()
	defer os.Exit(0)
	validation.GetTxPool(tp)
	// Create a new instance of the CLI client

	cmd := cli.Client{}
	for i := 0; i < 100; i++ {
		transaction := transaction.Transaction{
			To:    common.HexToAddress("0x1234567890123456789012345678901234567890"),
			Value: 1000,
			Nonce: 1,
		}
		signedTx := api.SignTxn(transaction)

		tp.AddTransactionToTxPool(&signedTx)
		// transactions[i] = transaction
	}
	api.GetTxPool(tp)
	network.GetTxPool(tp)

	TimerWithCallback := func() {
		// tp.GetAllTransactions()
	}

	// Start the ticker to execute the callback function every 2 seconds
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	go func() {
		for {
			select {
			case <-ticker.C:
				TimerWithCallback()
			}
		}
	}()

	// Start the HTTP server in a goroutine
	go func() {

		// Register API handlers
		http.HandleFunc("/sendTx", api.SendTxHandler)
		http.HandleFunc("/sendUnsignedTx", api.SendUnsignedTxHandler)
		http.HandleFunc("/blockNumber", api.BlockNumberHandler)
		http.HandleFunc("/getNonce", api.GetNonceHandler)
		http.HandleFunc("/getBalance", api.GetBalanceHandler)
		http.HandleFunc("/getKnownHosts", api.GetKnownHostHandler)
		http.HandleFunc("/addAddress", api.CreateAccountHandler)

		// Start the HTTP server
		fmt.Println("Server is running on port 8005")
		if err := http.ListenAndServe(":8005", nil); err != nil {
			fmt.Printf("Failed to start HTTP server: %v\n", err)
		}
	}()
	cmd.Run()
	// // Generate 100 demo transactions
	// for i := 0; i < 100; i++ {
	// 	sgnTx := &transaction.SignedTransaction{
	// 		To:    common.Address{byte(i)},
	// 		Value: uint64(rand.Intn(1000)),
	// 		Nonce: uint64(i),
	// V:     big.NewInt(int64(rand.Intn(1000))),
	// R:     big.NewInt(int64(rand.Intn(1000))),
	// S:     big.NewInt(int64(rand.Intn(1000))),
	// 	}

	// 	err := tp.AddTransactionToTxPool(sgnTx)
	// 	if err != nil {
	// 		fmt.Println("Error adding transaction to pool:", err)
	// 		continue
	// 	}

	// 	encodedTx, err := rlp.EncodeToBytes(tp.Transactions[i])
	// 	if err != nil {
	// 		fmt.Println("Error encoding transaction:", err)
	// 		continue
	// 	}

	// 	// Add the serialized transaction to LevelDB
	// 	err = database.AddSignedTransactionToLevelDB(encodedTx)
	// 	if err != nil {
	// 		fmt.Println("Error adding transaction to LevelDB:", err)
	// 	}
	// }
	// database.PrintAllData()
}
