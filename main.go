package main

import (
	"Blockchain_Project/cli"
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	// os.Setenv("SomeVar", "SomeValue")

	someVar := os.Getenv("SomeVar")
	fmt.Println(someVar)
	defer os.Exit(0)
	cmd := cli.Client{}
	cmd.Run()
}
