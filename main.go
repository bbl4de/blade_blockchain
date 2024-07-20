package main

import (
	"time"

	"github.com/bbl4de/blade_blockchain/network"
)

// Server - container
// Transport Layer => tcp, udp,
// Block
// Tx
// Keypair

func main() {
	// Transport for our local node, later we will add peers which are remote nodes - servers in the network
	trLocal := network.NewLocalTransport("LOCAL")
	// this is confusing, but we're making a FAKE remote transport, even though it's local
	trRemote := network.NewLocalTransport("REMOTE")

	// connect two transport layers 
	trLocal.Connect(trRemote) 
	trRemote.Connect(trLocal)

	// remote node will keep sending messages to our local node every 1 second
	go func() {
		
		for{
			trRemote.SendMessage(trLocal.Addr(), []byte("Hello world"))
			time.Sleep(1 * time.Second)
		}
	}()

	opts := network.ServerOpts{
		Transports: []network.Transport{trLocal},
	}	
	// then we configure and start our local node/server to listen for the messages
	s := network.NewServer(opts)
	s.Start()
}