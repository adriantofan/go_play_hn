package main

var config = struct {
	DateFormat string
}{
	"2006-01-02 15:04:05",
}

func main() {
	const queryCountURL string = "/1/queries/count/"
	const topNURL string = "/1/queries/count/"
	_, _, _ = readData("hn_logs.tsv")
	// log.Fatal(http.ListenAndServe(":8080", nil))
}
