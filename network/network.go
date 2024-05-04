package network

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

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

type Message struct {
	ID   uint64      `json:"id"`
	Code int         `json:"code"`
	Want *int        `json:"want,omitempty"`
	Data interface{} `json:"data"`
}

func sendMessage(stream network.Stream, msg Message) {
	encoder := json.NewEncoder(stream)
	if err := encoder.Encode(msg); err != nil {
		fmt.Println("Error sending message:", err)
	}
}

func sendHelloMessage(ctx context.Context, host host.Host, peerID peer.ID, peerAddr multiaddr.Multiaddr, helloMessage Message) {
	fmt.Println("Sending Hello message...")

	stream, err := host.NewStream(ctx, peerID, "/Hello")
	if err != nil {
		fmt.Println("Error opening stream to peer:", err)
		return
	}
	defer stream.Close()

	sendMessage(stream, helloMessage)

	fmt.Println("Hello message sent successfully.")
}

func Run(ctx context.Context) {
	priv, pub, err := crypto.GenerateKeyPair(crypto.Ed25519, -1)
	if err != nil {
		panic(err)
	}
	fmt.Printf("the public key is %v \n ", pub)

	host2, err := libp2p.New(libp2p.Identity(priv),libp2p.ListenAddrStrings("/ip4/0.0.0.0/tcp/0"))

	if err != nil {
		panic(err)
	}

	fmt.Println("Addresses:", host2.Addrs())
	fmt.Println("ID:", host2.ID())
	fmt.Println("Peer_ADDR:", os.Getenv("PEER_ADDR"))
	hostAddr := host2.Addrs()[0].String()
	peerID := host2.ID()
	peerAddr := hostAddr + "/p2p/" + peerID.String()
	fmt.Println("Host_ADDR:", peerAddr)
	peerMA, err := multiaddr.NewMultiaddr(os.Getenv("PEER_ADDR"))
	if err != nil {
		panic(err)
	}
	peerAddrInfo, err := peer.AddrInfoFromP2pAddr(peerMA)
	if err != nil {
		panic(err)
	}

	if err := host2.Connect(context.Background(), *peerAddrInfo); err != nil {
		panic(err)
	}
	fmt.Println("Connected to", peerAddrInfo.String())

	host2.SetStreamHandler("/Hello", func(s network.Stream) {
		fmt.Println("Received stream from:", s.Conn().RemotePeer())
		decoder := json.NewDecoder(s)
		var msg Message
		if err := decoder.Decode(&msg); err != nil {
			fmt.Println("Error decoding message:", err)
			s.Close()
			return
		}
		fmt.Println("Received message:", msg)

		helloMessage := Message{
			ID:   0,   // Random unique identifier
			Code: 0,   // Message type for "Hello"
			Want: nil, // Nothing expected back
			Data: "Hello",
		}
		if msg.Want != nil {
			sendHelloMessage(ctx, host2, s.Conn().RemotePeer(), s.Conn().RemoteMultiaddr(), helloMessage)
		}
		// Do not close the stream here
	})

	// Send the initial Hello message after setting up the stream handler
	resCode := 1
	firstHelloMessage := Message{
		ID:   0,
		Code: 1,
		Want: &resCode, // expects a string Hello
		Data: "Hello",
	}
	sendHelloMessage(ctx, host2, peerAddrInfo.ID, peerMA, firstHelloMessage)

	// Handle termination signals
	<-ctx.Done()

}
