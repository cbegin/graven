package util

import (
	"encoding/base64"
	"crypto/aes"
	"crypto/cipher"
)

var (
	key = []byte("9riX.Jax2gvKy4%4{[H#Nd,E")
	iv = []byte("<'RnpW4E3.L:/Ax*")
)

// Simple obfuscation function to help avoid over-the-shoulder
// viewing of passwords and other data. No intended to be
// cryptographically secure.
func Cloak(text string) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	plaintext := []byte(text)
	cfb := cipher.NewCFBEncrypter(block, iv)
	ciphertext := make([]byte, len(plaintext))
	cfb.XORKeyStream(ciphertext, plaintext)
	return encodeBase64(ciphertext), nil
}

// Simple obfuscation function to help avoid over-the-shoulder
// viewing of passwords and other data. No intended to be
// cryptographically secure.
func Uncloak(text string) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	ciphertext, err := decodeBase64(text)
	if err != nil {
		return "", err
	}
	cfb := cipher.NewCFBEncrypter(block, iv)
	plaintext := make([]byte, len(ciphertext))
	cfb.XORKeyStream(plaintext, ciphertext)
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
