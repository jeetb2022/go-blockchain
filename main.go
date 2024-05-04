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

    // Run the CLI client
    cmd.Run()
}
