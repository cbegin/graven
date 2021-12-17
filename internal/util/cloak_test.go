package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShouldEncryptAndDecryptText(t *testing.T) {
	plainText := "af242bb9ab97a6e44fc2bb5c3555f25b3b737826"
	cipherText, err := Cloak(plainText)
	assert.NoError(t, err)
	decryptedText, err := Uncloak(cipherText)
	assert.NoError(t, err)
	assert.Equal(t, plainText, decryptedText)
}
