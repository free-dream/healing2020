package tools

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
)

func CFDDecrypter(keyStr string,ivStr string,data string) string{
	key := []byte(keyStr)
	ciphertext, _ := hex.DecodeString(data)

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	if len(ciphertext) < aes.BlockSize {
		panic("ciphertext too short")
	}
	//iv := ciphertext[:aes.BlockSize]
	//ciphertext = ciphertext[aes.BlockSize:]

    iv := []byte(ivStr)

	stream := cipher.NewCFBDecrypter(block, iv)

	// XORKeyStream can work in-place if the two arguments are the same.
	stream.XORKeyStream(ciphertext, ciphertext)

    return string(ciphertext)
}
