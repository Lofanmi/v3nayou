package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
)

// GenEncrypt 加密
func GenEncrypt(value, key []byte) (string, error) {
	// random iv
	iv := make([]byte, aes.BlockSize)
	rand.Read(iv)

	// 128 / 192 / 256
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// PKCS#7
	value = PKCS7Padding(value, aes.BlockSize)

	// CBC Mode
	blockMode := cipher.NewCBCEncrypter(block, iv)
	encrypted := make([]byte, len(value))
	blockMode.CryptBlocks(encrypted, value)

	// HMAC sha256
	mac := hmac.New(sha256.New, key)
	mac.Write(iv)
	mac.Write(encrypted)

	// MAC(32) + IV(16) + PAYLOAD
	payload := append(mac.Sum(nil), iv...)
	payload = append(payload, encrypted...)

	// base64 Encode
	s := base64.RawURLEncoding.EncodeToString(payload)
	// s = strings.Replace(s, "=", "", -1)
	return s, nil
}

// GenDecrypt 解密
func GenDecrypt(encrypted string, key []byte) ([]byte, error) {
	// base64 Decode
	payload, err := base64.RawURLEncoding.DecodeString(encrypted)
	if err != nil {
		return nil, err
	}

	if len(payload) < 32+16 {
		return nil, errors.New("Invalid payload")
	}

	// MAC(32) + IV(16) + PAYLOAD
	messageMac := payload[:32]
	iv := payload[32:48]
	value := payload[48:]

	// Check MAC
	mac := hmac.New(sha256.New, key)
	mac.Write(payload[32:])
	expectedMac := mac.Sum(nil)
	if !hmac.Equal(messageMac, expectedMac) {
		return nil, nil
	}

	// 128 / 192 / 256
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// CBC Mode
	blockMode := cipher.NewCBCDecrypter(block, iv)
	decrypted := make([]byte, len(value))
	blockMode.CryptBlocks(decrypted, value)

	// PKCS#7
	decrypted = PKCS7UnPadding(decrypted)

	return decrypted, nil
}

// PKCS7Padding PKCS7
func PKCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

// PKCS7UnPadding PKCS7
func PKCS7UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}
