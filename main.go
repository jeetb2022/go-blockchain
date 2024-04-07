package main

import (
	"Blockchain_Project/cli"
	"Blockchain_Project/network"
	"context"
	"os"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	ctx := context.Background()
	defer os.Exit(0)
	cmd := cli.Client{}
	cmd.Run()
	network.Run(ctx)
}
