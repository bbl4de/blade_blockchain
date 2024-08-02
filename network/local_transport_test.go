package network

import (
	"io/ioutil"
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
	b, err := ioutil.ReadAll(rpc.Payload) 
		assert.Nil(t,err)
		assert.Equal(t, b, msg)
	assert.Equal(t, rpc.From, tra.addr)	
}

func TestBroadcast(t *testing.T) {
	tra := NewLocalTransport("A")
	trb := NewLocalTransport("B")
	trc := NewLocalTransport("C")

	tra.Connect(trb) 
	tra.Connect(trc) 

	msg := []byte("hello world!")
	assert.Nil(t, tra.Broadcast(msg))

	rpcb := <- trb.Consume()
	b, err := ioutil.ReadAll(rpcb.Payload)
	assert.Nil(t,err)
	assert.Equal(t, b, msg)

	rpcc := <- trc.Consume()
	c, err := ioutil.ReadAll(rpcc.Payload)
	assert.Nil(t,err)
	assert.Equal(t, c, msg)

	
}