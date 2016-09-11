package client

type HttpClient struct {
	hostname string
	port     int
}

func (c HttpClient) Store(id, payload []byte) (aesKey []byte, err error) {

	var key []byte
	aesKey = key

	var problem error
	err = problem

	return
}

func (c HttpClient) Retrieve(id, aesKey []byte) (payload []byte, err error) {

	var result []byte
	payload = result

	var problem error
	err = problem

	return
}
