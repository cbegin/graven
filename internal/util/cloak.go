package util

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
)

var (
	key   = []byte("9riX.Jax2gvKy4%4{[H#Nd,E")
	nonce = []byte("04f6d37b804f")
)

// Simple obfuscation function to help avoid over-the-shoulder
// viewing of passwords and other data. Not intended to be
// cryptographically secure.
func Cloak(text string) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}

	ciphertext := aesgcm.Seal(nil, nonce, []byte(text), nil)
	return encodeBase64(ciphertext), nil
}

// Simple obfuscation function to help avoid over-the-shoulder
// viewing of passwords and other data. No intended to be
// cryptographically secure.
func Uncloak(text string) (string, error) {
	ciphertext, err := decodeBase64(text)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}

	plaintext, err := aesgcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		panic(err.Error())
	}

	return string(plaintext), nil
}

func encodeBase64(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}

func decodeBase64(s string) ([]byte, error) {
	data, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return nil, err
	}
	return data, nil
}
