package main

import (
    "Blockchain_Project/cli"
    "os"
)

func main() {
	defer os.Exit(0)
	cmd := cli.Client{}
	cmd.Run()
}

