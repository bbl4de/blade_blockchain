package core

import (
	"testing"
	"time"

	"github.com/bbl4de/blade_blockchain/crypto"
	"github.com/bbl4de/blade_blockchain/types"
	"github.com/stretchr/testify/assert"
)


func randomBlock(height uint32) *Block {
	header := &Header{
		Version: 1,
		PrevBlockHash: types.RandomHash(), 
		Height: height,
		Timestamp: time.Now().UnixNano(),
	}

	tx := Transaction {
		Data: []byte("Hello, World!"),
		
	}

	return NewBlock(header, []Transaction{tx})
}

func TestSignBlock(t *testing.T) {
	privKey := crypto.GeneratePrivateKey()
	b := randomBlock(0)
	assert.Nil(t, b.Sign(privKey))
	assert.NotNil(t, b.Signature)
}

func TestVerifyBlock(t *testing.T) {
	privKey := crypto.GeneratePrivateKey()
	b := randomBlock(0)
	
	assert.Nil(t, b.Sign(privKey))
	assert.Nil(t, b.Verify())

	otherPrivKey := crypto.GeneratePrivateKey()
	b.Validator = otherPrivKey.PublicKey()
	assert.NotNil(t, b.Verify())

	// with different height the block is altered so it has to return an error
	b.Height = 100
	assert.NotNil(t, b.Verify())
}
