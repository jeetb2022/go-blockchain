package cli

import (
	// "Blockchain_Project/account"

	"Blockchain_Project/network"
	"context"
	"encoding/hex"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"syscall"

	"github.com/ethereum/go-ethereum/common"
)

type Client struct{}

func (cli *Client) validateArgs() {
	fmt.Println(os.Args)
	if len(os.Args) < 2 {
		cli.printUsage()
		runtime.Goexit()
	}
}

func (cli *Client) printUsage() {
	fmt.Println("Usage:")
	fmt.Println("go run ./main.go initiatenode -mineraddr <HEX_ADDRESS> ")

}

func convertToAddr(mineraddr string) (common.Address, error) {
	var address common.Address

	// Ensure the hexadecimal address string starts with "0x"
	if !strings.HasPrefix(mineraddr, "0x") {
		return address, fmt.Errorf("invalid hexadecimal address: must start with '0x'")
	}

	// Trim the "0x" prefix
	hexAddr := mineraddr[2:]

	// Decode the hexadecimal string into bytes
	addrBytes, err := hex.DecodeString(hexAddr)
	if err != nil {
		return address, err
	}

	// Check if the byte slice has the correct length
	if len(addrBytes) != 20 {
		return address, fmt.Errorf("invalid address length: must be 20 bytes long")
	}
	fmt.Println(addrBytes)
	// Copy the bytes into the address array
	copy(address[:], addrBytes)

	return address, nil
}

var minerAddr common.Address

func (cli *Client) Run() {

	fmt.Println("Running the CLI command")
	cli.validateArgs()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	InitiateNodeFlag := flag.NewFlagSet("initializeNode", flag.ExitOnError)
	switch os.Args[1] {
	case "initiatenode":
		minerAddress := InitiateNodeFlag.String("mineraddr", "", "Send mining rewards to the address ")
		err := InitiateNodeFlag.Parse(os.Args[2:])
		if err != nil {
			cli.printUsage()
			panic(err)
		}

		addrBytes, err := convertToAddr(*minerAddress)
		if err != nil {
			cli.printUsage()
			panic(err)
		}
		fmt.Println(addrBytes)
		// if !account.ValidateAddress(addrBytes) {
		// 	fmt.Println("Invalid address")
		// 	return // Exit early if the address is invalid
		// }
		minerAddr = common.HexToAddress(*minerAddress)
		// If the address is valid, execute network.Run(ctx)

		done := Start(ctx)

		// Handle termination signals
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
		<-sigCh

		// Signal the network to stop
		cancel()
		<-done

		network.Run(ctx, minerAddr)
	default:
		cli.printUsage()
	}
}

func Start(ctx context.Context) <-chan struct{} {
	done := make(chan struct{})
	go func() {
		defer close(done)
		network.Run(ctx, minerAddr)
	}()
	return done
}
