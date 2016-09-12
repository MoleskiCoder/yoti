package main

import (
	"fmt"
	"strconv"

	"encoding/hex"

	"github.com/MoleskiCoder/yoti/client"
)

func Store(connection client.HttpClient, id uint64, payload []byte) []byte {

	key, problem := connection.Store([]byte(strconv.FormatUint(id, 10)), payload)

	if problem != nil {
		panic("Store error")
	}

	return key
}

func Retrieve(connection client.HttpClient, id uint64, key []byte) string {

	payload, problem := connection.Retrieve([]byte(strconv.FormatUint(id, 10)), key)

	if problem != nil {
		panic("Retrieve error")
	}

	return string(payload)

}

func main() {

	var connection = client.New("http", "localhost", 8080)

	var id uint64 = 1
	original := "The quick brown fox jumps over the lazy dog"

	key := Store(connection, id, []byte(original))
	fmt.Printf("Stored key: %s\n", hex.EncodeToString(key))

	retrieved := Retrieve(connection, id, key)
	fmt.Printf("Retrieved text: %s\n", retrieved)
}
