package main

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const queryCountURL string = "/1/queries/count/"
const topNURL string = "/1/queries/popular/"

var config = struct {
	dateFormat                string
	logFatal                  func(v ...interface{})
	trie                      Trie
	logCount                  int
	errorCount                int
	computeDistinctQueryCount func(Trie, string) int
	computeTopNQueries        func(trie Trie, path string, params url.Values) []QueryCountPair
}{
	"2006-01-02 15:04:05", // configuration
	log.Fatal,             // production configuration
	MakeTrie(),            // data from file
	0,                     // data from file
	0,                     // data from file
	ComputeDistinctQueryCount, // production configuration
	ComputeTopNQueries,        // production configuration
}

// Returns a formatted JSON answering the distinct query. Uses config.distinctQueryCountHandler to do the hard lifting
func distinctQueryCountHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("content-type", "application/json")
	urlCount := config.computeDistinctQueryCount(config.trie, r.URL.Path)
	result, _ := json.Marshal(struct {
		Count int `json:"count"`
	}{urlCount})
	w.Write(result)
}

// Returns a formatted JSON answering the top n query. Uses config.computeTopNQueries to do the hard lifting
func topNHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("content-type", "application/json")
	results := config.computeTopNQueries(config.trie, r.URL.Path, r.URL.Query())
	result, _ := json.Marshal(struct {
		Queries []QueryCountPair `json:"queries"`
	}{
		results,
	})
	w.Write(result)
}

// Registers the handlers and starts the server. Loads and processes data asynchronously
func main() {
	go func() {
		startTime := time.Now()
		log.Println("Loading data")
		config.trie, config.logCount, config.errorCount = readData("hn_logs.tsv")
		config.trie.ComputeSortedURLs()
		log.Println("Data loaded in ", time.Now().Sub(startTime).Seconds(), " seconds")
	}()

	http.HandleFunc(queryCountURL, func(w http.ResponseWriter, r *http.Request) {
		distinctQueryCountHandler(w, r)
	})

	http.HandleFunc(topNURL, func(w http.ResponseWriter, r *http.Request) {
		topNHandler(w, r)
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		statusHandler(w, r)
	})
	log.Println("Server on http://localhost:8080")
	config.logFatal(http.ListenAndServe(":8080", nil))
}

// ComputeTopNQueries computes distinct query counts for urls such as /1/queries/count/2015-08-03
func ComputeTopNQueries(trie Trie, path string, params url.Values) []QueryCountPair {
	dateString := strings.TrimPrefix(path, topNURL)
	dateComponents := LogDateComponentsFromString(dateString)
	if len(dateComponents) == 0 {
		return make([]QueryCountPair, 0)
	}
	count := 5
	countStrs, foundCount := params["size"]
	if foundCount && len(countStrs) > 0 {
		if parsedCount, parsed := strconv.ParseInt(countStrs[0], 10, 64); parsed != nil {
			count = int(parsedCount)
		}
	}
	return TopNAtDate(trie, dateComponents, count)
}

// ComputeDistinctQueryCount computes distinct query counts for urls such as /1/queries/count/2015-08-03
func ComputeDistinctQueryCount(trie Trie, path string) int {
	var urlCount int
	dateString := strings.TrimPrefix(path, queryCountURL)
	dateComponents := LogDateComponentsFromString(dateString)
	if len(dateComponents) > 0 {
		urlCount = Distinct(trie, dateComponents)
	}
	return urlCount
}

func statusHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("content-type", "application/json")
	if config.logCount == 0 {
		w.Header().Add("Refresh", "1")
	}
	statusMessage, _ := json.Marshal(struct {
		LineCount  int `json:"line_count"`
		ErrorCount int `json:"error_count"`
	}{config.logCount, config.errorCount})
	w.Write(statusMessage)
}
