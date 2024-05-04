package main

import (
	"Blockchain_Project/api"
	"Blockchain_Project/cli"
	"fmt"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	godotenv.Load()
	defer os.Exit(0)

	// Create a new instance of the CLI client
	cmd := cli.Client{}

	// Start the HTTP server in a goroutine
	go func() {
		// Register API handlers
		http.HandleFunc("/sendTx", api.SendTxHandler)
		http.HandleFunc("/blockNumber", api.BlockNumberHandler)
		http.HandleFunc("/getNonce", api.GetNonceHandler)
		http.HandleFunc("/getBalance", api.GetBalanceHandler)

		// Start the HTTP server
		fmt.Println("Server is running on port 8000")
		if err := http.ListenAndServe(":8000", nil); err != nil {
			fmt.Printf("Failed to start HTTP server: %v\n", err)
		}
	}()
	cmd.Run()

	// tp := txpool.NewTransactionPool()

	// // Generate 100 demo transactions
	// for i := 0; i < 100; i++ {
	// 	sgnTx := &transaction.SignedTransaction{
	// 		To:    common.Address{byte(i)},
	// 		Value: uint64(rand.Intn(1000)),
	// 		Nonce: uint64(i),
	// 		V:     big.NewInt(int64(rand.Intn(1000))),
	// 		R:     big.NewInt(int64(rand.Intn(1000))),
	// 		S:     big.NewInt(int64(rand.Intn(1000))),
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
