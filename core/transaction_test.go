package core

import (
	"bytes"
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
	tx.From = otherPrivKey.PublicKey()	
	// make sure it returns an error
	assert.NotNil(t, tx.Verify())
}

func TestTxEncodeDecode(t *testing.T) {
	tx := randomTxWithSignature(t)
	buf := &bytes.Buffer{}
	assert.Nil(t, tx.Encode(NewGobTxEncoder(buf)))

	txDecoded := new(Transaction)
	assert.Nil(t, txDecoded.Decode(NewGobTxDecoder(buf)))
	assert.Equal(t, tx, txDecoded)
}

func randomTxWithSignature(t *testing.T) *Transaction {
	privKey := crypto.GeneratePrivateKey()
	tx := &Transaction {
		Data: []byte("Hello, World!"),
	}

	assert.Nil(t, tx.Sign(privKey))
	return tx
}

