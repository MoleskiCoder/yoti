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

	connection := client.New("http", "localhost", 8080)

	keys := map[uint64][]byte{}

	{
		{
			var id uint64 = 1
			original := "The quick brown fox jumps over the lazy dog"

			key := Store(connection, id, []byte(original))
			keys[id] = key
			fmt.Printf("Id: %d, Stored key: %s\n", id, hex.EncodeToString(key))
		}

		{
			var id uint64 = 2
			original := "Pack my box with five dozen liquor jugs"

			key := Store(connection, id, []byte(original))
			keys[id] = key
			fmt.Printf("Id: %d, Stored key: %s\n", id, hex.EncodeToString(key))
		}
	}

	{
		{
			var id uint64 = 1
			retrieved := Retrieve(connection, id, keys[id])
			fmt.Printf("Id: %d, Retrieved text: %s\n", id, retrieved)
		}
		{
			var id uint64 = 2
			retrieved := Retrieve(connection, id, keys[id])
			fmt.Printf("Id: %d, Retrieved text: %s\n", id, retrieved)
		}
	}
}
