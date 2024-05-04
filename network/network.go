package network

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/multiformats/go-multiaddr"
)

// func main() {

//		// Run(ctx)
//	}
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
	fmt.Printf("The public key is %v \n", pub)

	host, err := libp2p.New(libp2p.Identity(priv), libp2p.ListenAddrStrings("/ip4/0.0.0.0/tcp/0"))
	if err != nil {
		panic(err)
	}
	defer host.Close()

	fmt.Printf("The host ID is %s \n", host.ID())
	fmt.Printf("The host address is %s \n", host.Addrs()[0])

	hostAddr := host.Addrs()[0].String()
	peerID := host.ID()
	peerAddr := hostAddr + "/p2p/" + peerID.String()
	fmt.Printf("The peer address is %s \n", peerAddr)

	host.SetStreamHandler("/Hello", func(s network.Stream) {
		fmt.Println("Received stream from:", s.Conn().RemotePeer())
		decoder := json.NewDecoder(s)
		var msg Message
		if err := decoder.Decode(&msg); err != nil {
			fmt.Println("Error decoding message:", err)
			s.Close()
			return
		}
		helloMessage := Message{
			ID:   0,   // Random unique identifier
			Code: 0,   // Message type for "Hello"
			Want: nil, // Nothing expected back
			Data: "Hello from peer 1!",
		}
		fmt.Println("Received message:", msg.Data)
		if msg.Want != nil {
			sendHelloMessage(ctx, host, s.Conn().RemotePeer(), s.Conn().RemoteMultiaddr(), helloMessage)
		}
		s.Close()
	})

	// Define the callback function
	TimerWithCallback := func() {
		peers := host.Network().Peers()

		fmt.Println("Callback function called")
		fmt.Println("Connected peers:")
		for _, peer := range peers {
			addrs := host.Network().Peerstore().Addrs(peer)
			fmt.Printf("Peer ID: %s, Addresses: %v\n", peer, addrs)
		}
	}

	// Start the ticker to execute the callback function every 2 seconds
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	go func() {
		for {
			select {
			case <-ticker.C:
				TimerWithCallback()
			case <-ctx.Done():
				return
			}
		}
	}()

	// Wait for the program to be terminated
<- ctx.Done()
}
