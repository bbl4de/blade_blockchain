package network

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func connectLayers() (*LocalTransport, *LocalTransport){
	tra := NewLocalTransport("A")
	trb := NewLocalTransport("B")

	tra.Connect(trb)
	trb.Connect(tra)

	return tra.(*LocalTransport), trb.(*LocalTransport)
}

func TestConnect(t *testing.T) {
	 tra,trb := connectLayers()
	 assert.Equal(t, tra.peers[trb.addr], trb)
	 assert.Equal(t, trb.peers[tra.addr], tra)
}

func TestSendMessage(t *testing.T) {
	tra,trb := connectLayers()

	msg := []byte("hello world!")
	assert.Nil(t, tra.SendMessage(trb.addr, msg))
	
	rpc := <- trb.Consume()
	assert.Equal(t, rpc.Payload, msg)
	assert.Equal(t, rpc.From, tra.addr)
	
}