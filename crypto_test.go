package libsignal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRandIsRand(t *testing.T) {
	randint1 := randUint32()
	randint2 := randUint32()
	assert.NotEqual(t, randint1, randint2)
}

func TestMAC(t *testing.T) {
	key := make([]byte, 32)
	randBytes(key)
	msg := make([]byte, 100)
	randBytes(msg)
	macced := appendMAC(key, msg)
	assert.True(t, verifyMAC(key, macced[:100], macced[100:]))
}
func TestAES(t *testing.T) {
	key := make([]byte, 32)
	randBytes(key)
	msg := make([]byte, 100)
	randBytes(msg)
	encrypted, err := aesEncrypt(key, msg)
	assert.Nil(t, err)
	decrypted, err := aesDecrypt(key, encrypted)
	assert.Nil(t, err)
	assert.Equal(t, decrypted, msg)
}
