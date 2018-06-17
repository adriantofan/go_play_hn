package main

import (
	"log"
)

var config = struct {
	DateFormat string
	logFatal   func(v ...interface{})
}{
	"2006-01-02 15:04:05",
	log.Fatal,
}

func main() {
	const queryCountURL string = "/1/queries/count/"
	const topNURL string = "/1/queries/count/"
	_, _, _ = readData("hn_logs.tsv")
	// logFatal(http.ListenAndServe(":8080", nil))
}
