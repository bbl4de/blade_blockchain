package network

type NetAddr string

// message to be sent to the transport layer
type RPC struct {
	From NetAddr
	Payload []byte
}

// Transport is a module on the server and the server needs to have access to all messages sent to transport layers - it can to that with Consume method
type Transport interface {
	Consume() <- chan RPC
	Connect(Transport) error
	SendMessage(NetAddr, []byte) error
	Addr() NetAddr
}