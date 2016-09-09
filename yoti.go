package main

import (
	"fmt"
	"net/http"
)

//////////////////

type handler func(w http.ResponseWriter, r *http.Request)

//////////////////

// Should be POST
func StoreHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "store: %s, method: %s", r.URL.Path[1:], r.Method)
}

// Should be GET
func RetrieveHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "retrieve: %s, method: %s", r.URL.Path[1:], r.Method)
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

func main() {
	http.HandleFunc("/store", PostOnly(StoreHandler))
	http.HandleFunc("/retrieve", GetOnly(RetrieveHandler))
	http.ListenAndServe(":8080", nil)
}
