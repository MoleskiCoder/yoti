package client

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"

	"bytes"
	"net/http"

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

func (c HttpClient) CreateUrl(resource string) url.URL {

	var Url url.URL
	Url.Scheme = c.scheme
	Url.Host = c.hostname + ":" + strconv.FormatInt(int64(c.port), 10)

	Url.Path += resource

	return Url
}

func (c HttpClient) Store(id, payload []byte) (aesKey []byte, err error) {

	Url := c.CreateUrl("/store")
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
	err = decoder.Decode(&parsedResponse)

	return []byte(parsedResponse.Key), nil
}

func (c HttpClient) Retrieve(id, aesKey []byte) (payload []byte, err error) {

	Url := c.CreateUrl("/retrieve")
	parameters := url.Values{}
	parameters.Add("id", string(id))
	parameters.Add("key", string(aesKey))
	Url.RawQuery = parameters.Encode()

	fmt.Printf("** Encoded URL is %q\n", Url.String())

	var result []byte
	payload = result

	var problem error
	err = problem

	return
}
