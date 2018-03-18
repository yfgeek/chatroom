package core

import (
	"crypto/cipher"
	"bytes"
	"fmt"
	"math/rand"
	"crypto/aes"
	"os"
	"encoding/base64"
	"strings"
)

const PasswordLength = 32

type Password [PasswordLength]byte

type Cipher struct{
	Key []byte
}

func RandPassword() Password {
	intArr := rand.Perm(PasswordLength)
	password := Password{}
	for i, v := range intArr {
		password[i] = byte(v)
		if i == v {
			return RandPassword()
		}
	}
	return password
}


func (c *Cipher) EncryptMessage(msg []byte) ([]byte, error){
	block, err := aes.NewCipher(c.Key)
	if err != nil {
		fmt.Printf("Error: NewCipher(%d bytes) = %s", len(c.Key), err)
		os.Exit(-1)
	}
	blockSize := block.BlockSize()
	msg = cPKCS5Padding(msg, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, c.Key[:blockSize])
	crypted := make([]byte, len(msg))
	blockMode.CryptBlocks(crypted, msg)
	return crypted, nil
}

func (c *Cipher) DecryptMessage(crypted []byte) ([]byte, error) {
	block, err := aes.NewCipher(c.Key)
	if err != nil {
		fmt.Printf("Error: NewCipher(%d bytes) = %s", len(c.Key), err)
		os.Exit(-1)
	}
	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, c.Key[:blockSize])
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

func (password *Password) String() string {
	return base64.StdEncoding.EncodeToString(password[:])
}


func ParsePassword(passwordString string) (Password, error) {
	bs, err := base64.StdEncoding.DecodeString(strings.TrimSpace(passwordString))
	if err != nil || len(bs) != PasswordLength {
		os.Exit(-1)
	}
	password := Password{}
	copy(password[:], bs)
	bs = nil
	return password, nil
}
