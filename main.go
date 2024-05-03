package main

import (
	"Blockchain_Project/cli"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	defer os.Exit(0)
	cmd := cli.Client{}
	cmd.Run()
}
