package client

import (
	"encoding/json"
	"net/url"
	"strconv"

	"bytes"
	"net/http"

	"encoding/hex"

	"github.com/MoleskiCoder/yoti/server"
)

type HttpClient struct {
	scheme   string
	hostname string
	port     int
}

func New(scheme string, hostname string, port int) HttpClient {
	var connection HttpClient
	connection.scheme = scheme
	connection.hostname = hostname
	connection.port = port
	return connection
}

func (c HttpClient) createUrl(resource string) url.URL {

	var Url url.URL
	Url.Scheme = c.scheme
	Url.Host = c.hostname + ":" + strconv.FormatInt(int64(c.port), 10)

	Url.Path += resource

	return Url
}

func (c HttpClient) Store(id, payload []byte) ([]byte, error) {

	Url := c.createUrl("/store")
	Id, _ := strconv.ParseUint(string(id), 10, 64)

	request := &server.StoreRequest{
		Id:   Id,
		Data: string(payload)}
	jsonRequest, _ := json.Marshal(request)

	httpRequest, _ := http.NewRequest("POST", Url.String(), bytes.NewBuffer(jsonRequest))
	httpRequest.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, _ := client.Do(httpRequest)
	defer response.Body.Close()

	decoder := json.NewDecoder(response.Body)

	var parsedResponse server.StoreResponse
	_ = decoder.Decode(&parsedResponse)

	key, _ := hex.DecodeString(parsedResponse.Key)
	return key, nil
}

func (c HttpClient) Retrieve(id, aesKey []byte) ([]byte, error) {

	Url := c.createUrl("/retrieve")
	parameters := url.Values{}
	parameters.Add("id", string(id))
	parameters.Add("key", hex.EncodeToString(aesKey))
	Url.RawQuery = parameters.Encode()

	httpRequest, _ := http.NewRequest("GET", Url.String(), nil)
	client := &http.Client{}
	response, _ := client.Do(httpRequest)

	decoder := json.NewDecoder(response.Body)

	var parsedResponse server.RetrieveResponse
	_ = decoder.Decode(&parsedResponse)

	decodedData, _ := hex.DecodeString(parsedResponse.Data)
	return decodedData, nil
}
