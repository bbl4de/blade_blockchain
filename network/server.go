package network

import (
	"bytes"
	"os"
	"time"

	"github.com/bbl4de/blade_blockchain/core"
	"github.com/bbl4de/blade_blockchain/crypto"
	"github.com/bbl4de/blade_blockchain/types"
	"github.com/go-kit/log"
)

var defaultBlockTime = 5 * time.Second
type ServerOpts struct{
	ID string
	Logger log.Logger
	RPCDecodeFunc RPCDecodeFunc
	RPCProcessor RPCProcessor
	Transports []Transport
	BlockTime time.Duration
	PrivateKey *crypto.PrivateKey
}

// Container containing every module
// This will hold the transaction mempool as well
type Server struct {
	ServerOpts
	memPool *TxPool
	chain *core.Blockchain
	isValidator bool
	rpcCh chan RPC
	quitCh chan struct{}
}

func NewServer(opts ServerOpts) (*Server, error) {
	if opts.BlockTime == time.Duration(0) {
		opts.BlockTime = defaultBlockTime
	}
	if opts.RPCDecodeFunc == nil {
		opts.RPCDecodeFunc = DefaultRPCDecodeFunc
	}

	if opts.Logger == nil {
		opts.Logger = log.NewLogfmtLogger(os.Stderr)
		opts.Logger = log.With(opts.Logger, "ID", opts.ID)
	}

	chain, err := core.NewBlockchain(genesisBlock())
	if err != nil {
		return nil, err
	}

	s := &Server {
		ServerOpts: opts,
		chain: chain,
		memPool: NewTxPool(),
		isValidator: opts.PrivateKey != nil,
		rpcCh: make(chan RPC),
		quitCh: make(chan struct{}, 1),
	}

	// if we don't have any processor from the server options, we assume the server is the default processor
	if s.RPCProcessor == nil {
		s.RPCProcessor = s
	}

	if s.isValidator {
		go s.validatorLoop()
	}

	return s, nil
}


func (s *Server) Start() {
	s.initTransports()

	free:
	for {
		// keep checking whether there is something to consumefv from the RPC channel
		// if there is print it out
		// if not RPC channel is it suposed to quit?
		// if not we need something else - default statement
		select {
		case rpc := <- s.rpcCh:
			// handle a message
			msg, err := s.RPCDecodeFunc(rpc)
			if err != nil {
				s.Logger.Log("error", err)
			}

			if err := s.RPCProcessor.ProcessMessage(msg); err != nil {
								s.Logger.Log("error", err)

			}
		case <- s.quitCh:
			break free
		
	}
}
	s.Logger.Log("msg","Server is shutting down")

}

func (s *Server) validatorLoop() {
	ticker := time.NewTicker(s.BlockTime)

	s.Logger.Log("msg", "Starting validator loop",
"blockTime", s.BlockTime,)

	for {
		<- ticker.C
		s.createNewBlock()
	}
}

func (s *Server) ProcessMessage(msg *DecodedMessage) error {
	switch t := msg.Data.(type) {
	case *core.Transaction:
			return s.processTransaction( t)
	}
	return nil
}

func (s *Server) broadcast(payload []byte) error {
	for _, tr := range s.Transports {
		if err := tr.Broadcast(payload); err != nil {
			return err
		}
	}

	return nil
}

func (s *Server) processTransaction(tx *core.Transaction) error {
	hash := tx.Hash(core.TxHasher{})

	if s.memPool.Has(hash) {
		return nil
	}
	
	if err := tx.Verify(); err != nil {
		return err 
	}

	tx.SetFirstSeen(time.Now().UnixNano())

	s.Logger.Log(
		"msg", "adding new tx to mempool", 
		"hash", hash, 
		"mempoolLength", s.memPool.Len(),
	)

	// broadcast tx to peers

	go s.broadcastTx(tx) 

	return s.memPool.Add(tx)

}

func (s *Server) broadcastTx(tx *core.Transaction) error {
	buf := &bytes.Buffer{}
	if err := tx.Encode(core.NewGobTxEncoder(buf)); err != nil {
		return err
	}

	msg:= NewMessage(MessageTypeTx, buf.Bytes())
	return s.broadcast(msg.Bytes())

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

func (s *Server) createNewBlock() error {
	currentHeader, err := s.chain.GetHeader(s.chain.Height())
	if err != nil {
		return err
	}

	// for now no transactions as we will have to define the amount of transactions allowed in one block
	block, err := core.NewBlockFromPrevHeader(currentHeader, nil)
	if err != nil {
		return err
	}
	
	if err := block.Sign(*s.PrivateKey); err != nil {
		return err
	}

	if err := s.chain.AddBlock(block); err != nil {
		return err
	}

	return nil
}

func genesisBlock() *core.Block {
	header := &core.Header{
		Version : 1, 
		DataHash : types.Hash{}, 
		Height: 0,
		Timestamp : time.Now().UnixNano(),
	}

	// TODO test the error
	b,_ :=  core.NewBlock(header, nil)
	return b
}