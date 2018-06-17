package main

import (
	"fmt"
	"log"
	"net/http"
)

var database []record

func statusHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("content-type", "application/json")
	fmt.Fprintf(w, "{'linecount':%d}", len(database))
}

func main() {
	// readData("/Users/adriantofan/devel/a_test/hn_logs_short.tsv")
	database, _, _ = readData("/Users/adriantofan/devel/a_test/hn_logs.tsv")
	http.HandleFunc("/", statusHandler)
	http.HandleFunc("/", statusHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
