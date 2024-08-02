package main

import (
	"bytes"
	"log"
	"math/rand"
	"strconv"
	"time"

	"github.com/bbl4de/blade_blockchain/core"
	"github.com/bbl4de/blade_blockchain/crypto"
	"github.com/bbl4de/blade_blockchain/network"
	"github.com/sirupsen/logrus"
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
			if err := sendTransaction(trRemote, trLocal.Addr()); err != nil {
				logrus.Error(err)
			}
			time.Sleep(1 * time.Second)
		}
	}()

	privKey := crypto.GeneratePrivateKey()

	opts := network.ServerOpts{
		PrivateKey: &privKey,
		ID: "LOCAL",
		Transports: []network.Transport{trLocal},
	}	
	// then we configure and start our local node/server to listen for the messages
	s, err := network.NewServer(opts)
	if err != nil {
		log.Fatal(err)
	}

	s.Start()
}

func sendTransaction(tr network.Transport, to network.NetAddr) error {
	// send a transaction to the network
	privKey := crypto.GeneratePrivateKey()
	data := []byte(strconv.FormatInt(int64(rand.Intn(10000)), 10))
	tx := core.NewTransaction(data)
	tx.Sign(privKey)
	buf := &bytes.Buffer{}
	if err := tx.Encode(core.NewGobTxEncoder(buf)); err != nil {
		return err
	}

	msg := network.NewMessage(network.MessageTypeTx, buf.Bytes())

	return tr.SendMessage(to, msg.Bytes())
}