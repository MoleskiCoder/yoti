package server

import (
	"encoding/json"
	"fmt"
	"net/http"
)

//////////////////

type handler func(w http.ResponseWriter, r *http.Request)

type StoreRequest struct {
	ID   int    `json:"id"`
	Data string `json:"data"`
}

type StoreResponse struct {
	Key string `json:"key"`
}

type RetrieveResponse struct {
	Data string `json:"data"`
}

//////////////////

// Must be POST
func StoreHandler(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)

	var parsed StoreRequest
	error := decoder.Decode(&parsed)
	if error != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	fmt.Fprintf(w, "store: %s, method: %s", r.URL.Path[1:], r.Method)
	fmt.Fprintf(w, "	ID: %d, Data: %s", parsed.ID, parsed.Data)
}

// Must be GET
func RetrieveHandler(w http.ResponseWriter, r *http.Request) {

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

func PostOnly(h handler) handler {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			h(w, r)
			return
		}
		http.Error(w, "post only", http.StatusMethodNotAllowed)
	}
}

func GetOnly(h handler) handler {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			h(w, r)
			return
		}
		http.Error(w, "get only", http.StatusMethodNotAllowed)
	}
}

//////////////////

func Start() {
	http.HandleFunc("/store", PostOnly(StoreHandler))
	http.HandleFunc("/retrieve", GetOnly(RetrieveHandler))
	http.ListenAndServe(":8080", nil)
}
