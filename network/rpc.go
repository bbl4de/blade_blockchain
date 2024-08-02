package network

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io"

	"github.com/bbl4de/blade_blockchain/core"
	"github.com/sirupsen/logrus"
)

// Making an enum in golang for the MessageType
type MessageType byte

const (
	MessageTypeTx MessageType = 0x1
	MessageTypeBlock // this will automatically increment by 1
)

// message to be sent to the transport layer
type RPC struct {
	From NetAddr
	Payload io.Reader
}

type Message struct {
	Header MessageType
	Data []byte
}

func NewMessage(t MessageType, data []byte) *Message {
	return &Message {
		Header: t, 
		Data: data,
	}
}

func (msg *Message) Bytes() []byte {
	buf := &bytes.Buffer{}
	gob.NewEncoder(buf).Encode(msg)
	return buf.Bytes()
}

type DecodedMessage struct {
	From NetAddr
	Data any 
}

type RPCDecodeFunc func(RPC) (*DecodedMessage, error)

func DefaultRPCDecodeFunc(rpc RPC) (*DecodedMessage, error) {
	msg := Message{}
	if err := gob.NewDecoder(rpc.Payload).Decode(&msg); err != nil {
		return nil, fmt.Errorf("failed to decode message from this %s: %s",rpc.From, err)
	}
	logrus.WithFields(logrus.Fields{
		"from":rpc.From,
		"type": msg.Header, 
	}).Debug("new message incoming")

	switch msg.Header {
	// for now we only implement this type of message, meaning that we get the transaction either from a node that broadcast it's transaction to us
	case MessageTypeTx:
		tx := new(core.Transaction)
		if err := tx.Decode(core.NewGobTxDecoder(bytes.NewReader(msg.Data))); err != nil {
			return nil, err
		}

		return &DecodedMessage{From: rpc.From, Data: tx}, nil

	default: 
		return nil, fmt.Errorf("invalid message header %x", msg.Header)
	}
}



// processor will take the decoded messages from the handler and process them
type RPCProcessor interface {
	ProcessMessage(*DecodedMessage) error 
}

