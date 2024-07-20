package network

import (
	"fmt"
	"time"
)

type ServerOpts struct{
	Transports []Transport
}

// Container containing every module
// This will hold the transaction mempool as well
type Server struct {
	ServerOpts
	rpcCh chan RPC
	quitCh chan struct{}
}

func NewServer(opts ServerOpts) *Server {
	return &Server {
		ServerOpts: opts,
		rpcCh: make(chan RPC),
		quitCh: make(chan struct{}, 1),
	}
}

func (s *Server) Start() {
	s.initTransports()

	ticker := time.NewTicker(5 * time.Second)

	free:
	for {
		// keep checking whether there is something to consume from the RPC channel
		// if there is print it out
		// if not RPC channel is it suposed to quit?
		// if not we need something else - default statement
		select {
		case rpc := <- s.rpcCh:
			// handle a message
			fmt.Printf("%+v\n", rpc)
		case <- s.quitCh:
			break free
		case <- ticker.C:
			fmt.Println("Do stuff every x seconds")
		}
	}
	fmt.Println("Server shutdown")
}

func (s *Server) initTransports() {
	// For each Transport Layer in the server we will spin up a new go routine and keep reading the channels - keep Consuming their message channels
	for _, tr := range s.Transports {
		go func(tr Transport) {
			for rpc := range tr.Consume() {
				// there is nothing that tells us whether the message from this go routine is thread safe right now 
				// TODO  we will need some kind of a handler in the Start() function
				s.rpcCh <- rpc
			}
		} (tr)
	}
}