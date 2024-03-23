package cli

import (
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
	fmt.Println((os.Args))
}
