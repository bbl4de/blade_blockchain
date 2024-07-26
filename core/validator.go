package core

import "fmt"

type Validator interface {
	ValidateBlock(b *Block) error
}

type BlockValidator struct {
	bc *Blockchain
}

func NewBlockValidator(bc *Blockchain) *BlockValidator {
	return &BlockValidator{
		bc: bc,
	}
}

// Validates the block and all it's transactions ( inside Verify )
func (v *BlockValidator) ValidateBlock(b *Block) error {
	// check if the block already exists
	if v.bc.HasBlock(b.Height) {
		return fmt.Errorf("block %d already exists and has a hash %s", b.Height, b.Hash(BlockHasher{}))
	}
	// check if the block is the next one in the chain
	if b.Height != v.bc.Height()+1  {
		return fmt.Errorf("block height %d is too high", b.Height)
	}
	// check if the hash of the previous block is correct
	prevHeader, err := v.bc.GetHeader(b.Height - 1)
	
	if err != nil {
		return err
	}

	hash := BlockHasher{}.Hash(prevHeader)
	if hash != b.PrevBlockHash {
		return fmt.Errorf("the hash of the previous block (%s) is incorrect", b.PrevBlockHash)
	}
	// check if the block is correctly signed
	if err := b.Verify(); err != nil {
		return err
	}

	return nil
}
