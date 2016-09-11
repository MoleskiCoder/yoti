package main

import (
	"fmt"
	"strconv"

	"github.com/MoleskiCoder/yoti/client"
)

func Store(id uint64, payload string) string {

	var crypt client.HttpClient

	key, problem := crypt.Store([]byte(strconv.FormatUint(id, 10)), []byte(payload))

	if problem != nil {
		panic("Store error")
	}

	return string(key)
}

func Retrieve(id uint64, key string) string {

	var crypt client.HttpClient

	payload, problem := crypt.Retrieve([]byte(strconv.FormatUint(id, 10)), []byte(key))

	if problem != nil {
		panic("Retrieve error")
	}

	return string(payload)

}

func main() {

	var id uint64 = 1
	original := "The quick brown fox jumps over the lazy dog"

	key := Store(id, original)
	fmt.Printf("Stored key: %s\n", key)

	retrieved := Retrieve(id, key)
	fmt.Printf("Retrieved text: %s\n", retrieved)
}
