package core

import (
	"testing"

	"github.com/bbl4de/blade_blockchain/crypto"
	"github.com/stretchr/testify/assert"
)

func TestSignTransactio(t *testing.T) {
	privKey := crypto.GeneratePrivateKey()	
	tx := &Transaction {
		Data: []byte("Hello, World!"),
	}

	assert.Nil(t, tx.Sign(privKey))
	assert.NotNil(t, tx.Signature)
}

func TestVerifyTransaction(t *testing.T) {
	privKey := crypto.GeneratePrivateKey()
	tx := &Transaction {
		Data: []byte("Hello, World!"),
	}

	assert.Nil(t, tx.Sign(privKey))
	assert.Nil(t, tx.Verify())

	otherPrivKey := crypto.GeneratePrivateKey()
	// tamper the existing transaction with new public key - which did not sign that transaction
	tx.PublicKey = otherPrivKey.PublicKey()	
	// make sure it returns an error
	assert.NotNil(t, tx.Verify())
}