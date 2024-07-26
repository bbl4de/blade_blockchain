package network

import (
	"fmt"
	"time"

	"github.com/bbl4de/blade_blockchain/core"
	"github.com/bbl4de/blade_blockchain/crypto"
	"github.com/sirupsen/logrus"
)

var defaultBlockTime = 5 * time.Second
type ServerOpts struct{
	Transports []Transport
	BlockTime time.Duration
	PrivateKey *crypto.PrivateKey
}

// Container containing every module
// This will hold the transaction mempool as well
type Server struct {
	ServerOpts
	blockTime time.Duration
	memPool *TxPool
	isValidator bool
	rpcCh chan RPC
	quitCh chan struct{}
}

func NewServer(opts ServerOpts) *Server {
	if opts.BlockTime == time.Duration(0) {
		opts.BlockTime = defaultBlockTime
	}
	return &Server {
		ServerOpts: opts,
		blockTime: opts.BlockTime,
		memPool: NewTxPool(),
		isValidator: opts.PrivateKey != nil,
		rpcCh: make(chan RPC),
		quitCh: make(chan struct{}, 1),
	}
}

func (s *Server) Start() {
	s.initTransports()

	ticker := time.NewTicker(s.blockTime)

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
			// consensus logic
			if s.isValidator{	
				s.createNewBlock()	
			}
		}
	}
	fmt.Println("Server shutdown")
}

func (s *Server) handleTransaction(tx *core.Transaction) error {
	if err := tx.Verify(); err != nil {
		return err 
	}

	hash := tx.Hash(core.TxHasher{})

	if s.memPool.Has(hash) {
		logrus.WithFields(logrus.Fields{
		"hash": hash,
	}).Info("transaction already in the mempool")
		return nil
	}

	logrus.WithFields(logrus.Fields{
		"hash": hash,
	}).Info("adding a new tx to the mempool")

	return s.memPool.Add(tx)
}

func (s *Server) createNewBlock() error {
	fmt.Println("creating a new block")
	return nil
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