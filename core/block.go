package core

import (
	"bytes"
	"encoding/gob"
	"fmt"

	"github.com/bbl4de/blade_blockchain/crypto"
	"github.com/bbl4de/blade_blockchain/types"
)

type Header struct {
	Version uint32
	DataHash types.Hash // hash of transaction data
	PrevBlockHash types.Hash
	Height uint32
	Timestamp int64

}


func (h *Header) Bytes() []byte {
buf := &bytes.Buffer{}
	enc := gob.NewEncoder(buf)

	enc.Encode(h)

	return buf.Bytes()

}
type Block struct {
	*Header
	Transactions []Transaction
	Validator crypto.PublicKey
	Signature *crypto.Signature
	// Catched version of the header hash
	hash types.Hash
}

func NewBlock(h *Header, txx []Transaction) *Block {
	return &Block{
		Header: h,
		Transactions: txx,
	}
}

func (b *Block) AddTransaction(tx *Transaction) {
	b.Transactions = append(b.Transactions, *tx)
}

func (b *Block) Sign(privKey crypto.PrivateKey) error {
	sig, err := privKey.Sign(b.Header.Bytes())
	if err != nil {	
		return err
	}

	b.Validator = privKey.PublicKey()	
	b.Signature = sig
	return nil
}

// Verify the signature of the block and every transaction this block contains
func (b *Block) Verify() error {	
	if b.Signature == nil {
		return fmt.Errorf("block has no signature")
	}
	if !b.Signature.Verify(b.Validator, b.Header.Bytes()) {
		return fmt.Errorf("block has invalid signature")
	}

	for _, tx := range b.Transactions {
		if err := tx.Verify(); err != nil {
			return err
		}
	}
	return nil
}

func (b *Block) Decode(dec Decoder[*Block]) error {
	return dec.Decode(b)

}

func (b *Block) Encode(enc Encoder[*Block]) error {
	return enc.Encode(b)

}

func (b *Block) Hash(hasher Hasher[*Header]) types.Hash {
	if b.hash.IsZero() {
		b.hash = hasher.Hash(b.Header)
	}	

	return b.hash
}	
