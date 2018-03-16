package core

import (
	"crypto/cipher"
	"bytes"
	"fmt"
	"crypto/aes"
	"os"
)

type Chiper struct{
	Key []byte
}


func (c *Chiper) EncryptMessage(msg []byte) ([]byte, error){
	chiper, err := aes.NewCipher(c.Key)
	if err != nil {
		fmt.Printf("Error: NewCipher(%d bytes) = %s", len(c.Key), err)
		os.Exit(-1)
	}
	blockSize := chiper.BlockSize()
	msg = cPKCS5Padding(msg, blockSize)
	blockMode := cipher.NewCBCEncrypter(chiper, c.Key[:blockSize])
	crypted := make([]byte, len(msg))
	blockMode.CryptBlocks(crypted, msg)
	return crypted, nil
}

func (c *Chiper) DecryptMessage(crypted []byte) ([]byte, error) {
	chiper, err := aes.NewCipher(c.Key)
	if err != nil {
		fmt.Printf("Error: NewCipher(%d bytes) = %s", len(c.Key), err)
		os.Exit(-1)
	}
	blockSize := chiper.BlockSize()
	blockMode := cipher.NewCBCDecrypter(chiper, c.Key[:blockSize])
	origData := make([]byte, len(crypted))
	blockMode.CryptBlocks(origData, crypted)
	origData = cPKCS5UnPadding(origData)
	return origData, nil
}


func cPKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func cPKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}
