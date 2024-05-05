package network

import (
	"Blockchain_Project/blockchain"
	"Blockchain_Project/database"
	"context"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"

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
	fmt.Printf("the public key is %v \n ", pub)

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
	KnownHosturl := "http://10.1.153.234:8000/getKnownHosts"
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

func Run(ctx context.Context) {
	ctxt = ctx
	host, err := StartNewNode()
	if err != nil {
		panic(err)
	}

	GetMultiAddr()
	ConnectToPeers(host)

	host.SetStreamHandler("/Hello", func(s network.Stream) {
		// fmt.Println("Received stream from:", s.Conn().RemotePeer())
		msg, err := receiveMessage(s)
		if err != nil {
			panic(err)
		}
		fmt.Println(string(msg.Data))
		if msg.Want == uint(1) {
			SendPONG(ctx, host, s.Conn().RemotePeer())
		}

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
		fmt.Println("List of peer addresses:")
		addresses = append(addresses, hostPeerAddr)
		for _, addr := range addresses {
			fmt.Println(addr)
		}
		peerAddrList = addresses

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

			sendMessageWithCTX(ctx, host, peerAddrInfo.ID, msgPING)
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

func SendNewBlock(ctx context.Context, host host.Host, peerID peer.ID, block blockchain.Block) error {
	encodedBlock, err := rlp.EncodeToBytes(block)
	if err != nil {
		return fmt.Errorf("error occurred while encoding block: %v", err)
	}

	msgNewBlock := Message{
		ID:   rand.Uint64(),
		Code: uint(5),
		Want: uint(542),
		Data: encodedBlock,
	}
	sendMessageWithCTX(ctx, host, peerID, msgNewBlock)

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

func SendBlock(ctx context.Context, host host.Host, peerID peer.ID, blocks []blockchain.Block) error {
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
