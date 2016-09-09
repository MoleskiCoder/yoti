package main

import (
    "fmt"
    "net/http"
)

func store_handler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "store: %s", r.URL.Path[1:])
}
func retrieve_handler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "retrieve: %s", r.URL.Path[1:])
}

func main() {
    http.HandleFunc("/store", store_handler)
    http.HandleFunc("/retrieve", retrieve_handler)
    http.ListenAndServe(":8080", nil)
}
