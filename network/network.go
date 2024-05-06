package network

import (
	"Blockchain_Project/blockchain"
	"Blockchain_Project/database"
	"Blockchain_Project/transaction"
	"Blockchain_Project/txpool"
	"Blockchain_Project/utils"
	"context"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/multiformats/go-multiaddr"
)

// func main() {

// 	// Run(ctx)
// }

var peerIdList []string
var peerAddrList []string
var ctxt context.Context
var hostPeerAddr string
var globalHost host.Host
var localLatestBlock uint64
var minerAdderss common.Address

var tp *txpool.TransactionPool

type Message struct {
	ID   uint64 `json:"id"`
	Code uint   `json:"code"`
	Want uint   `json:"want,omitempty"`
	Data []byte `json:"data"`
}

func sendMessage(stream network.Stream, msg Message) {
	// Serialize the message
	encodedMsg, err := SerializeMessage(msg)
	if err != nil {
		fmt.Println("Error serializing message:", err)
		return
	}

	// Write the serialized message to the stream
	_, err = stream.Write(encodedMsg)
	if err != nil {
		fmt.Println("Error sending message:", err)
	}
}

// SerializeMessage serializes the Message struct into a byte slice
func SerializeMessage(msg Message) ([]byte, error) {
	// Encode the message struct
	encodedMsg, err := rlp.EncodeToBytes(msg)
	if err != nil {
		return nil, err
	}
	return encodedMsg, nil
}

func DeserializeMessage(encodedMsg []byte, msg *Message) error {
	err := rlp.DecodeBytes(encodedMsg, msg)
	if err != nil {
		return err
	}
	return nil

}

func receiveMessage(stream network.Stream) (Message, error) {
	var msg Message

	// Read the bytes from the stream
	buf := make([]byte, 1024) // Adjust the buffer size as needed
	n, err := stream.Read(buf)
	if err != nil {
		return msg, err
	}

	// Deserialize the message using your custom deserialization function
	err = DeserializeMessage(buf[:n], &msg)
	if err != nil {
		return msg, err
	}

	return msg, nil
}

func sendMessageWithCTX(ctx context.Context, host host.Host, peerID peer.ID, helloMessage Message) {
	stream, err := host.NewStream(ctx, peerID, "/Hello")
	if err != nil {
		fmt.Println("Error opening stream to peer:", err)
		return
	}
	defer stream.Close()
	sendMessage(stream, helloMessage)
}

func StartNewNode() (host.Host, error) {
	priv, pub, err := crypto.GenerateKeyPair(crypto.Ed25519, -1)
	if err != nil {
		return nil, err
	}
	pub = pub
	// fmt.Printf("the public key is %v \n ", pub)

	host2, err := libp2p.New(libp2p.Identity(priv), libp2p.ListenAddrStrings("/ip4/0.0.0.0/tcp/0"))

	if err != nil {
		return nil, err
	}

	fmt.Println("Addresses:", host2.Addrs())
	fmt.Println("ID:", host2.ID())
	// fmt.Println("Peer_ADDR:", os.Getenv("PEER_ADDR"))
	hostAddr := host2.Addrs()[0].String()
	peerID := host2.ID()
	peerAddr := hostAddr + "/p2p/" + peerID.String()
	fmt.Println("Host_ADDR:", peerAddr)
	hostPeerAddr = peerAddr
	return host2, nil
}

func ConnectToPeers(host host.Host) {
	for _, addr := range peerAddrList {
		if addr != "" {

			fmt.Println(addr)
			peerMA, err := multiaddr.NewMultiaddr(addr)
			if err != nil {
				panic(err)
			}
			peerAddrInfo, err := peer.AddrInfoFromP2pAddr(peerMA)
			if err != nil {
				panic(err)
			}
			if err := host.Connect(context.Background(), *peerAddrInfo); err != nil {
				panic(err)
			}
			fmt.Println("Connected to", peerAddrInfo.String())

		}
	}

}
func GetMultiAddr() { // []multiaddr.Multiaddr,[]peer.AddrInfo
	KnownHosturl := "http://10.1.153.234:8009/getKnownHosts"
	resp, err := http.Get(KnownHosturl)
	if err != nil {
		log.Fatalf("Error making GET request: %v", err)
	}
	defer resp.Body.Close()

	// Check if the response status code is OK
	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Unexpected status code: %v", resp.StatusCode)
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response body: %v", err)
	}

	// Convert the byte slice to a string
	addrString := string(body)

	// Split the string into individual addresses based on newline delimiter
	addrList := strings.Split(addrString, "\n")

	// Print the list of addresses
	fmt.Println("List of addresses:")
	for _, addr := range addrList {
		fmt.Println(addr)
	}
	peerAddrList = addrList

}

func Run(ctx context.Context, minerAddr common.Address) {
	ctxt = ctx
	minerAdderss = minerAddr
	locallatestblock, err := database.GetCurrentHeight()
	if err != nil {
		panic(err)
	}
	localLatestBlock = locallatestblock
	host, err := StartNewNode()
	if err != nil {
		panic(err)
	}
	globalHost = host

	GetMultiAddr()
	ConnectToPeers(host)

	host.SetStreamHandler("/Hello", func(s network.Stream) {
		msg, err := receiveMessage(s)
		if err != nil {
			panic(err)
		}
		// fmt.Println(string(msg.Data))
		if msg.Want == uint(1) {
			SendPONG(ctx, host, s.Conn().RemotePeer())
		}
		if msg.Code == uint(4) {
			RecvedTransaction, err := transaction.DeserializeTransaction(msg.Data)
			if err != nil {
				panic(err)
			}
			fmt.Println("This transaction is from the peer", RecvedTransaction)
			tp.AddTransactionToTxPool(&RecvedTransaction)
		}
		if msg.Code == uint(5) {
			// var RecvedBlock blockchain.Block
			// err := rlp.DecodeBytes(msg.Data, &RecvedBlock)
			// if err != nil {
			// 	panic(err)
			// }

			// fmt.Println("This blockchain is from the peer", msg.Data)
		}
		s.Close()

	})

	SendPING(ctx, host)

	GetUpdatedPeerList := func() {
		var addresses []string

		for _, conn := range host.Network().Conns() {
			// Extract peer addresses from the connection
			remoteAddr := conn.RemoteMultiaddr().String()

			// Extract Peer ID
			remotePeerID := conn.RemotePeer()

			// Construct the multiaddress with Peer ID
			remoteAddrWithPeerID := fmt.Sprintf("%s/p2p/%s", remoteAddr, remotePeerID)

			// Add the peer addresses to the list
			addresses = append(addresses, remoteAddrWithPeerID)
		}
		// fmt.Println("List of peer addresses:")
		addresses = append(addresses, hostPeerAddr)
		// for _, addr := range addresses {
		// 	// fmt.Println(addr)
		// }
		peerAddrList = addresses
		time.Sleep(time.Second * 5)
		if true {
			err := SendNewBlock()
			if err != nil {
				panic(err)
			}
		}

	}

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	go func() {
		for {
			select {
			case <-ticker.C:
				GetUpdatedPeerList()
			}
		}
	}()
	// Handle termination signals
	<-ctx.Done()

}

func SendPING(ctx context.Context, host host.Host) {

	msgPING := Message{
		ID:   rand.Uint64(),
		Code: uint(0),
		Want: uint(1),
		Data: []byte("PING"),
	}
	for _, addr := range peerAddrList {
		if addr != "" {

			peerMA, err := multiaddr.NewMultiaddr(addr)
			if err != nil {
				panic(err)
			}
			peerAddrInfo, err := peer.AddrInfoFromP2pAddr(peerMA)
			if err != nil {
				panic(err)
			}
			if hostPeerAddr != addr {
				sendMessageWithCTX(ctx, host, peerAddrInfo.ID, msgPING)
			}
		}
	}

}

func SendPONG(ctx context.Context, host host.Host, peerID peer.ID) {
	msgPONG := Message{
		ID:   rand.Uint64(),
		Code: uint(1),
		Want: uint(542),
		Data: []byte("PONG"),
	}
	sendMessageWithCTX(ctx, host, peerID, msgPONG)
}

func SendNewBlock() error {

	minedBlock := CreateBlocks(minerAdderss)

	// fmt.Println("Mined block: ", minedBlock)

	encodedBlock, err := rlp.EncodeToBytes(minedBlock)
	if err != nil {
		return fmt.Errorf("error occurred while encoding block: %v", err)
	}

	msgNewBlock := Message{
		ID:   rand.Uint64(),
		Code: uint(5),
		Want: uint(10),
		Data: encodedBlock,
	}
	for _, addr := range peerAddrList {
		if addr != "" {

			peerMA, err := multiaddr.NewMultiaddr(addr)
			if err != nil {
				panic(err)
			}
			peerAddrInfo, err := peer.AddrInfoFromP2pAddr(peerMA)
			if err != nil {
				panic(err)
			}
			if hostPeerAddr != addr {
				sendMessageWithCTX(ctxt, globalHost, peerAddrInfo.ID, msgNewBlock)
			}
		}
	}

	return nil
}

func SendGetBlock(ctx context.Context, host host.Host, peerID peer.ID, blockNumbers []int) error {
	encodedBlockNumbers, err := rlp.EncodeToBytes(blockNumbers)
	if err != nil {
		return fmt.Errorf("error occurred while encoding block numbers: %v", err)
	}

	msgGetBlock := Message{
		ID:   rand.Uint64(),
		Code: uint(6),
		Want: uint(7),
		Data: encodedBlockNumbers,
	}
	sendMessageWithCTX(ctx, host, peerID, msgGetBlock)

	return nil
}

func SendBlocks(ctx context.Context, host host.Host, peerID peer.ID, blocks []blockchain.Block) error {
	blocksToSend := make([]byte, len(blocks))
	for i := range blocks {
		block, err := database.RetrieveBlockHash(uint64(i))
		if err != nil {
			return fmt.Errorf("error occurred while retrieving block: %v", err)
		}
		serializedBlock, err := rlp.EncodeToBytes(block)
		blocksToSend = append(blocksToSend, serializedBlock...)
	}
	msgBlock := Message{
		ID:   rand.Uint64(),
		Code: uint(7),
		Want: uint(542),
		Data: blocksToSend,
	}
	sendMessageWithCTX(ctx, host, peerID, msgBlock)
	return nil
}

func SendTransaction(signedTransaction transaction.SignedTransaction) {
	fmt.Println(signedTransaction)
	encodedTransaction, err := transaction.SerializeTransaction(signedTransaction)

	if err != nil {
		panic(err)
	}
	msgTransaction := Message{
		ID:   rand.Uint64(),
		Code: uint(4),
		Want: uint(542),
		Data: encodedTransaction,
	}

	for _, addr := range peerAddrList {
		if addr != "" {

			peerMA, err := multiaddr.NewMultiaddr(addr)
			if err != nil {
				panic(err)
			}
			peerAddrInfo, err := peer.AddrInfoFromP2pAddr(peerMA)
			if err != nil {
				panic(err)
			}
			if hostPeerAddr != addr {
				sendMessageWithCTX(ctxt, globalHost, peerAddrInfo.ID, msgTransaction)
			}
		}
	}

}

func SendGetLatestBlock(ctx context.Context, host host.Host, peerID peer.ID) {
	msgGetLatestBlock := Message{
		ID:   rand.Uint64(),
		Code: uint(8),
		Want: uint(9),
	}
	sendMessageWithCTX(ctx, host, peerID, msgGetLatestBlock)
}

func SendGetLatestBlockResponse(ctx context.Context, host host.Host, peerID peer.ID) error {
	latestBlockNumber, err := database.GetCurrentHeight()
	if err != nil {
		panic(err)
	}
	encodedLatestBlockNumber, err := rlp.EncodeToBytes(latestBlockNumber)
	if err != nil {
		return fmt.Errorf("error occurred while encoding latest block number: %v", err)
	}

	msgGetLatestBlockResponse := Message{
		ID:   rand.Uint64(),
		Code: uint(9),
		Want: uint(542),
		Data: encodedLatestBlockNumber,
	}
	sendMessageWithCTX(ctx, host, peerID, msgGetLatestBlockResponse)

	return nil
}

func GetPeerAddrs() []string {
	return peerAddrList
}

func GetTxPool(p *txpool.TransactionPool) {
	tp = p
}
func GetMinerAddr() common.Address {
	return minerAdderss
}

// func GetMinedBlock() *blockchain.Block {
// 	// Create a new block

// 	block := &blockchain.Block{
// 		Transactions: make([]*transaction.SignedTransaction, 0),
// 	}

// 	// Maximum number of transactions to add to the block
// 	maxTransactions := 10
// 	var pickedTransactions []transaction.SignedTransaction
// 	// Iterate through transactions in the transaction pool
// 	for i := 0; i < len(tp.Transactions) && i < maxTransactions; i++ {
// 		if len(tp.Transactions) <= 0 {
// 			break
// 		}
// 		tx := tp.Transactions[i]

// 		// Add the transaction to the block
// 		block.Transactions = append(block.Transactions, tx)
// 		pickedTransactions := append(pickedTransactions, *tx)
// 		pickedTransactions = pickedTransactions

// 		// Remove the transaction from the pool
// 		tp.Transactions = append(tp.Transactions[:i], tp.Transactions[i+1:]...)

// 		// Update the index to handle the removed transaction
// 		i--

// 	}

// stateRoot := utils.StateRoot()
// transactionRoot := utils.CalculateTransactionsRoot(pickedTransactions)
// parentHash, err := database.GetLastBlockHash()
// parentHashBytes := database.RlpHash(parentHash)
// currentHeight, err := database.GetCurrentHeight()
// if err != nil {
// 	panic(err)
// }
// // if err != nil {
// // 	panic(err)
// // }
// timestamp := time.Now().Unix()
// minerAddr := minerAddr
// blockHeader := &blockchain.Header{
// 	ParentHash:       parentHashBytes,
// 	Miner:            minerAddr,
// 	StateRoot:        stateRoot,
// 	TransactionsRoot: transactionRoot,
// 	Number:           currentHeight,
// 	Timestamp:        uint64(timestamp),
// }
// extradata := blockchain.SignHeader(*blockHeader)
// blockHeader.ExtraData = extradata[:]
// block.Header = blockHeader
// return block
// }

func CreateBlocks(miner common.Address) blockchain.Block {
	// blockNumber := uint64(0)
	// for {
	// time.Sleep(2 * time.Second)
	if len(tp.Transactions) == 0 {
		emptyBlock := blockchain.Block{}
		return emptyBlock
	}
	var transactions []*transaction.SignedTransaction
	if len(tp.Transactions) > 3 {
		transactions = tp.Transactions[:3]    // taking first 10 transactions
		tp.Transactions = tp.Transactions[3:] // removing first 10 transactions from pool
	} else {
		transactions = tp.Transactions
		tp.Transactions = nil
	}

	var signedTransactions []transaction.SignedTransaction

	// Iterate over each transaction in the transactions slice
	for _, tx := range transactions {
		// Append each transaction to the signedTransactions slice
		signedTransactions = append(signedTransactions, *tx)
	}
	// fmt.Println("Signed Transaction", signedTransactions)
	// Create the block header
	stateRoot := utils.StateRoot()
	transactionRoot := utils.CalculateTransactionsRoot(signedTransactions)
	parentHash, err := database.GetLastBlockHash()
	parentHashBytes := database.RlpHash(parentHash)
	currentHeight, err := database.GetCurrentHeight()
	if err != nil {
		panic(err)
	}
	// if err != nil {
	// 	panic(err)
	// }
	timestamp := time.Now().Unix()
	minerAddr := minerAdderss
	blockHeader := &blockchain.Header{
		ParentHash:       parentHashBytes,
		Miner:            minerAddr,
		StateRoot:        stateRoot,
		TransactionsRoot: transactionRoot,
		Number:           currentHeight,
		Timestamp:        uint64(timestamp),
	}
	extradata := blockchain.SignHeader(*blockHeader)
	blockHeader.ExtraData = extradata[:]
	// return block
	block := blockchain.Block{
		Header:       blockHeader,
		Transactions: transactions,
	}

	// fmt.Println("Mineer:", minerAdderss)
	// fmt.Println("Header ParentHash:", block.Header.ParentHash)
	// fmt.Println("Header Miner:", block.Header.Miner)
	// fmt.Println("Header StateRoot:", block.Header.StateRoot)
	// fmt.Println("Header StateRoot:", block.Transactions[0])
	// Print other fields similarly
	// }
	return block
}
