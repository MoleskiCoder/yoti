package server

import (
	"encoding/hex"
	"encoding/json"
	"net/http"
	"strconv"
	"sync"
)

//////////////////

type handler func(w http.ResponseWriter, r *http.Request)

//////////////////

type Repository struct {
	sync.RWMutex
	encrypted map[uint64][]byte
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

	id := request.Id
	data := []byte(request.Data)

	key, err := GenerateKey()
	if err != nil {
		panic("Unable to generate AES256 key")
	}

	aesCrypt, err := NewAes(key)
	if err != nil {
		panic("Unable to create AES cypher")
	}

	encrypted := aesCrypt.Encrypt(data)

	repo.Lock()
	repo.encrypted[id] = encrypted
	repo.Unlock()

	response := &StoreResponse{
		Key: hex.EncodeToString(key),
	}

	return response
}

func (repo Repository) retrieve(key []byte, id uint64) *RetrieveResponse {

	repo.RLock()
	encrypted := repo.encrypted[id]
	repo.RUnlock()

	aesCrypt, err := NewAes(key)
	if err != nil {
		panic("Unable to create AES cypher")
	}

	decrypted := aesCrypt.Decrypt(encrypted)

	response := &RetrieveResponse{
		Data: hex.EncodeToString(decrypted),
	}
	return response
}

//////////////////

// Must be POST
func storeHandler(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)

	var parsed StoreRequest
	err := decoder.Decode(&parsed)
	if err != nil {
		panic("Unable to decode store request")
	}

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

	parsedId, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	decodedKey, err := hex.DecodeString(key)
	if err != nil {
		http.Error(w, "Invalid Key", http.StatusBadRequest)
		return
	}

	payload := repository.retrieve(decodedKey, parsedId)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(payload)
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
	}

	http.HandleFunc("/store", postOnly(storeHandler))
	http.HandleFunc("/retrieve", getOnly(retrieveHandler))
	http.ListenAndServe(":8080", nil)
}
