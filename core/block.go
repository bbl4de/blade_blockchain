package core

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io"

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

func (b *Block) Sign(privKey crypto.PrivateKey) error {
	sig, err := privKey.Sign(b.HeaderData())
	if err != nil {	
		return err
	}

	b.Validator = privKey.PublicKey()	
	b.Signature = sig
	return nil
}

func (b *Block) Verify() error {	
	if b.Signature == nil {
		return fmt.Errorf("block has no signature")
	}
	if !b.Signature.Verify(b.Validator, b.HeaderData()) {
		return fmt.Errorf("block has invalid signature")
	}
	return nil
}

func (b *Block) Decode(r io.Reader, dec Decoder[*Block]) error {
	return dec.Decode(r, b)

}

func (b *Block) Encode(w io.Writer, enc Encoder[*Block]) error {
	return enc.Encode(w, b)

}

func (b *Block) Hash(hasher Hasher[*Block]) types.Hash {
	if b.hash.IsZero() {
		b.hash = hasher.Hash(b)
	}	

	return b.hash
}	

func (b *Block) HeaderData() []byte {
	buf := &bytes.Buffer{} 
	enc := gob.NewEncoder(buf)

	enc.Encode(b.Header)

	return buf.Bytes()
}