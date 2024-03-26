package cli

import (
	"Blockchain_Project/account"
	"flag"
	"fmt"
	"os"
)

type Client struct{}

func (cli *Client) validateArgs() {
	fmt.Println(os.Args)
	// if len(os.Args) < 2 {
	// 	cli.printUsage()
	// 	runtime.Goexit()
	// }
}

func (cli *Client) printUsage() {

}

func (cli *Client) Run() {
	fmt.Println("Running the CLI command")

	testFlag := flag.NewFlagSet("test2", flag.ExitOnError)
	testFlagVal := testFlag.String(
		"test",
		"",
		"Enable mining mode and send reward to ADDRESS",
	)
	testFlagVal2 := testFlag.String(
		"test2",
		"",
		"Enable mining mode and send reward to ADDRESS",
	)
	testFlag.Parse(os.Args[1:])
	fmt.Println("test flag value :", *testFlagVal2)
	fmt.Println("test flag value :", *testFlagVal)
	type Address [20]byte
	address := Address{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20}

	newAccount := &account.Account{
		Address: address,
		Nonce:   13,
		Balance: 10000,
	}
	fmt.Println((os.Args))
	fmt.Println(newAccount.Address)
	fmt.Println(newAccount.Nonce)
	fmt.Println(newAccount.Balance)
}
