package main

import (
	"fmt"
	"strconv"

	"github.com/MoleskiCoder/yoti/client"
)

func Store(id int, payload string) string {

	var crypt client.Client

	key, problem := crypt.Store([]byte(strconv.FormatInt(id, 10)), []byte(payload))

	if problem != nil {
		panic("Store error")
	}

	return string(key)
}

func Retrieve(id int, key string) string {
	var crypt client.Client

	payload, problem := crypt.Retrieve([]byte(strconv.FormatInt(id, 10)), []byte(key))

	if problem != nil {
		panic("Retrieve error")
	}

	return string(payload)

}

func main() {

	id := 1
	original := "The quick brown fox jumps over the lazy dog"

	key := Store(id, original)
	fmt.Printf("Stored key: %s\n", key)

	retrieved := Retrieve(id, key)
	fmt.Printf("Retrieved text: %s\n", retrieved)
}
