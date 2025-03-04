package core

import (
	"testing"

	"github.com/bbl4de/blade_blockchain/types"
	"github.com/stretchr/testify/assert"
)

func TestAddBlock(t *testing.T) {
	bc :=newBlockchainWithGenesis(t)

	lenBlocks:=1000
	for i:=0; i < lenBlocks; i++ {
		b := randomBlock(t, uint32(i+1), getPrevBlockHash(t,bc,uint32(i+1)))
		err := bc.AddBlock(b)
		assert.Nil(t, err)
		assert.True(t, bc.HasBlock(uint32(i+1)))
	}

	assert.Equal(t, bc.Height(), uint32(lenBlocks))
	assert.Equal(t, len(bc.headers), lenBlocks+1) // genesis block

	assert.NotNil(t, bc.AddBlock(randomBlock(t,uint32(420), types.Hash{})))
}

func TestNewBlockchain(t *testing.T) {
	bc :=newBlockchainWithGenesis(t)

	assert.NotNil(t, bc.validator)
	assert.Equal(t, bc.Height(), uint32(0)) // genesis block saves us from the underflow
}

func TestGetHeader(t *testing.T) {
	bc :=newBlockchainWithGenesis(t)


	lenBlocks:=1000
	for i:=0; i < lenBlocks; i++ {
		b := randomBlock(t, uint32(i+1), getPrevBlockHash(t,bc,uint32(i+1)))
		assert.Nil(t, bc.AddBlock(b))
		header, err := bc.GetHeader(uint32(i+1))
		assert.Nil(t, err)
		assert.Equal(t,header,b.Header)
	}
	
}

func TestHashBlock(t *testing.T) {
	bc :=newBlockchainWithGenesis(t)
	
	assert.True(t, bc.HasBlock(0))
	assert.False(t, bc.HasBlock(1))
	assert.False(t, bc.HasBlock(1000))
}

func TestAddBlockTooHigh(t *testing.T) {
	bc:=newBlockchainWithGenesis(t)
	assert.Nil(t, bc.AddBlock(randomBlock(t,1, getPrevBlockHash(t,bc,uint32(1)))))
	assert.NotNil(t, bc.AddBlock(randomBlock(t,3,types.Hash{})))
}


func newBlockchainWithGenesis(t *testing.T) *Blockchain {
	bc, err := NewBlockchain(randomBlock(t,0, types.Hash{}))
	assert.Nil(t, err)

	return bc
}

func getPrevBlockHash(t *testing.T, bc *Blockchain, height uint32) types.Hash {
	prevHeader, err := bc.GetHeader(height - 1)
	assert.Nil(t, err)

	return BlockHasher{}.Hash(prevHeader)
}