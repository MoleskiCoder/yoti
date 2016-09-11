package server

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
)

type Aes struct {
	key   string
	iv    []byte
	block cipher.Block
}

func NewAes(key string, iv []byte) *Aes {
	constructed := &Aes{
		key: key,
		iv:  iv,
	}
	constructed.block, _ = aes.NewCipher([]byte(constructed.key))
	return constructed
}

func (current *Aes) Encrypt(data []byte) []byte {
	input := current.Pad(data)
	encoded := make([]byte, len(input))
	encrypter := cipher.NewCBCEncrypter(current.block, current.iv)
	encrypter.CryptBlocks(encoded, input)
	return encoded
}

func (current *Aes) Decrypt(data []byte) []byte {
	decoded := make([]byte, len(data))
	decrypter := cipher.NewCBCDecrypter(current.block, current.iv)
	decrypter.CryptBlocks(decoded, data)
	return current.Trim(decoded)
}

func (current *Aes) Pad(data []byte) []byte {

	dataLength := len(data)

	padding := aes.BlockSize - (dataLength % aes.BlockSize)

	length := aes.BlockSize * ((dataLength / aes.BlockSize) + 1)
	padded := make([]byte, length)

	copy(padded, data)

	// Each padding byte is the amount of padded data used
	for i := 0; i < length; i++ {
		padded[i] = byte(padding)
	}

	return padded
}

func (current *Aes) Trim(data []byte) []byte {
	// Retrieve the amount of padding first
	padding := int(data[len(data)-1])
	return data[0 : len(data)-padding]
}

func GenerateIV() []byte {
	// By reading "aes.Blocksize" bytes from crypt/rand
	iv := make([]byte, aes.BlockSize)
	_, _ = io.ReadFull(rand.Reader, iv)
	return iv
}
