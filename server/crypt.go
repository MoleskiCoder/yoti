package server

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
)

type Aes struct {
	key   string
	block cipher.Block
}

func NewAes(key string) (*Aes, error) {
	constructed := &Aes{
		key: key,
	}

	aesCipher, err := aes.NewCipher([]byte(constructed.key))
	if err != nil {
		return nil, err
	}

	constructed.block = aesCipher
	return constructed, nil
}

func (current *Aes) Encrypt(plain []byte) []byte {

	// Data to be encrypted *must* be padded to aes.BlockSize
	input := current.pad(plain)

	// Plonk the IV at the beginning of the package
	// i.e. packaged == iv + encoded
	packaged := make([]byte, len(input)+aes.BlockSize)
	iv := packaged[:aes.BlockSize]
	readRandom(iv)
	encoded := packaged[aes.BlockSize:]

	encrypter := cipher.NewCBCEncrypter(current.block, iv)
	encrypter.CryptBlocks(encoded, input)

	return packaged
}

func (current *Aes) Decrypt(packaged []byte) []byte {

	iv := packaged[:aes.BlockSize]
	encrypted := packaged[aes.BlockSize:]

	decoded := make([]byte, len(encrypted))
	decrypter := cipher.NewCBCDecrypter(current.block, iv)
	decrypter.CryptBlocks(decoded, encrypted)
	return current.trim(decoded)
}

func readRandom(iv []byte) (int, error) {
	return io.ReadFull(rand.Reader, iv)
}

func (current *Aes) pad(data []byte) []byte {

	dataLength := len(data)

	padding := aes.BlockSize - (dataLength % aes.BlockSize)

	length := aes.BlockSize * ((dataLength / aes.BlockSize) + 1)
	padded := make([]byte, length)

	// Each padding byte is the amount of padded data used
	for i := 0; i < length; i++ {
		padded[i] = byte(padding)
	}

	copy(padded, data)

	return padded
}

func (current *Aes) trim(data []byte) []byte {
	// Retrieve the amount of padding first
	padding := int(data[len(data)-1])
	return data[0 : len(data)-padding]
}
