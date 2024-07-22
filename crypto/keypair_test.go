package crypto

import (
	"testing"

	"github.com/stretchr/testify/assert"
)
func TestKeypair_Sign_Verify_Valid(t *testing.T) {
 	privKey := GeneratePrivateKey()
	pubKey := privKey.PublicKey()
	msg := []byte("hello, world!")

	sig, err := privKey.Sign(msg)
	assert.Nil(t, err)

	assert.True(t, sig.Verify(pubKey, msg))
}

func TestKeypair_Sign_Verify_Fail(t *testing.T) {
 	privKey := GeneratePrivateKey()
	pubKey := privKey.PublicKey()
	msg := []byte("hello, world!")

	sig, err := privKey.Sign(msg)
	assert.Nil(t, err)

	otherPubKey := GeneratePrivateKey().PublicKey()

	assert.False(t, sig.Verify(otherPubKey, msg))
	assert.False(t, sig.Verify(pubKey, []byte("not hello, world?")))
}