package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

var database []record

func statusHandler(w http.ResponseWriter, r *http.Request, lineCount int, errorCount int) {
	w.Header().Add("content-type", "application/json")
	statusMessage, _ := json.Marshal(struct {
		LineCount  int `json:"line_count"`
		ErrorCount int `json:"error_count"`
	}{lineCount, errorCount})
	fmt.Fprintf(w, string(statusMessage))
}

func main() {
	// readData("/Users/adriantofan/devel/a_test/hn_logs_short.tsv")
	var errorCount, lineCount int
	database, errorCount, lineCount = readData("hn_logs.tsv")
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		statusHandler(w, r, lineCount, errorCount)
	})
	log.Fatal(http.ListenAndServe(":8080", nil))
}
