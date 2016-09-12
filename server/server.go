package server

import (
	"encoding/json"
	"fmt"
	"net/http"
)

//////////////////

type handler func(w http.ResponseWriter, r *http.Request)

//////////////////

type Repository struct {
	encrypted map[uint64][]byte
	iv        map[uint64][]byte
}

//////////////////

type StoreRequest struct {
	Id   uint64 `json:"id"`
	Data string `json:"data"`
}

type StoreResponse struct {
	Key string `json:"key"`
}

type RetrieveResponse struct {
	Data string `json:"data"`
}

//////////////////

func (repo Repository) store(request StoreRequest) *StoreResponse {

	key := "AES256Key-32Characters1234567890"

	id := request.Id
	data := request.Data

	iv := GenerateIV()
	aesCrypt, _ := NewAes(key, iv)

	encrypted := aesCrypt.Encrypt([]byte(data))
	repo.encrypted[id] = encrypted
	repo.iv[id] = iv

	response := &StoreResponse{
		Key: key,
	}

	return response
}

func (repo Repository) retrieve(key string, id uint64) *RetrieveResponse {

	encrypted := repo.encrypted[id]
	iv := repo.iv[id]

	aesCrypt, _ := NewAes(key, iv)
	decrypted := aesCrypt.Decrypt(encrypted)

	response := &RetrieveResponse{
		Data: string(decrypted),
	}
	return response
}

//////////////////

// Must be POST
func storeHandler(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)

	var parsed StoreRequest
	_ = decoder.Decode(&parsed)

	payload := repository.store(parsed)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(payload)
}

// Must be GET
func retrieveHandler(w http.ResponseWriter, r *http.Request) {

	url := r.URL
	query := url.Query()

	key := query.Get("key")
	key_problem := len(key) == 0
	if key_problem {
		http.Error(w, "Missing key", http.StatusBadRequest)
	}

	id := query.Get("id")
	id_problem := len(id) == 0
	if id_problem {
		http.Error(w, "Missing ID", http.StatusBadRequest)
	}

	if key_problem || id_problem {
		return
	}

	fmt.Fprintf(w, "retrieve: %s, method: %s", url.Path[1:], r.Method)
	fmt.Fprintf(w, "	ID: %s, Key: %s", id, key)
}

//////////////////

func postOnly(h handler) handler {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			h(w, r)
			return
		}
		http.Error(w, "post only", http.StatusMethodNotAllowed)
	}
}

func getOnly(h handler) handler {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			h(w, r)
			return
		}
		http.Error(w, "get only", http.StatusMethodNotAllowed)
	}
}

//////////////////

var repository *Repository

func Start() {

	repository = &Repository{
		encrypted: map[uint64][]byte{},
		iv:        map[uint64][]byte{},
	}

	http.HandleFunc("/store", postOnly(storeHandler))
	http.HandleFunc("/retrieve", getOnly(retrieveHandler))
	http.ListenAndServe(":8080", nil)
}
