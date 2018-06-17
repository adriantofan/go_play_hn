package main

import (
	"fmt"
	"log"
	"net/http"
)

func statusHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("content-type", "application/json")
	fmt.Fprintf(w, "{}")
}

func main() {
	readData()
	return
	http.HandleFunc("/", statusHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
