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
	Id, err := strconv.ParseUint(string(id), 10, 64)
	if err != nil {
		return nil, err
	}

	request := &server.StoreRequest{
		Id:   Id,
		Data: string(payload)}
	jsonRequest, err := json.Marshal(request)
	if err != nil {
		panic("Unable to create store request JSON")
	}

	httpRequest, err := http.NewRequest("POST", Url.String(), bytes.NewBuffer(jsonRequest))
	if err != nil {
		panic("Unable to create HTTP POST request for store")
	}

	httpRequest.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(httpRequest)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	decoder := json.NewDecoder(response.Body)

	var parsedResponse server.StoreResponse
	err = decoder.Decode(&parsedResponse)
	if err != nil {
		return nil, err
	}

	key, err := hex.DecodeString(parsedResponse.Key)
	if err != nil {
		return nil, err
	}

	return key, nil
}

func (c HttpClient) Retrieve(id, aesKey []byte) ([]byte, error) {

	Url := c.createUrl("/retrieve")
	parameters := url.Values{}
	parameters.Add("id", string(id))
	parameters.Add("key", hex.EncodeToString(aesKey))
	Url.RawQuery = parameters.Encode()

	httpRequest, err := http.NewRequest("GET", Url.String(), nil)
	if err != nil {
		panic("Unable to create HTTP GET request for retrieve")
	}

	client := &http.Client{}
	response, err := client.Do(httpRequest)
	if err != nil {
		return nil, err
	}

	decoder := json.NewDecoder(response.Body)

	var parsedResponse server.RetrieveResponse
	err = decoder.Decode(&parsedResponse)
	if err != nil {
		return nil, err
	}

	decodedData, err := hex.DecodeString(parsedResponse.Data)
	if err != nil {
		return nil, err
	}

	return decodedData, nil
}
