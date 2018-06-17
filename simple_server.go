package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func statusHandler(w http.ResponseWriter, r *http.Request, database []record, lineCount int, errorCount int) {
	w.Header().Add("content-type", "application/json")
	statusMessage, _ := json.Marshal(struct {
		LineCount  int `json:"line_count"`
		ErrorCount int `json:"error_count"`
	}{lineCount, errorCount})
	w.Write(statusMessage)
}

// Handles query count requests as follows:
//
// * GET /1/queries/count/<DATE_PREFIX>: returns a JSON object specifying the number of distinct queries
// that have been done during a specific time range
// * Distinct queries done in 2015: GET /1/queries/count/2015: returns { count: 573697 }
// * Distinct queries done in Aug: GET /1/queries/count/2015-08: returns { count: 573697 }
// * Distinct queries done on Aug 3rd: GET /1/queries/count/2015-08-03: returns { count: 198117 }
// * Distinct queries done on Aug 1st between 00:04:00 and 00:04:59: GET /1/queries/count/2015-08-01 00:04:
// returns { count: 617 }
// Uses data loaded in to database. Uses urlPrefix to detect where the date starts in the url (if any)
// Delegates real work doing to getDistinctQueries

func queryCountHandler(w http.ResponseWriter, r *http.Request, database []record, urlPrefix string) {
	w.Header().Add("content-type", "application/json")
	urlCount := getDistinctQueries(database, urlPrefix, r.URL.Path)
	result, _ := json.Marshal(struct {
		Count int `json:"count"`
	}{urlCount})
	w.Write(result)
}

// GET /1/queries/popular/<DATE_PREFIX>?size=<SIZE>: returns a JSON object listing the top <SIZE> popular queries that have been done during a specific time range
// Top 3 popular queries done in 2015: GET /1/queries/popular/2015?size=3: returns

//     {
//       queries: [
//         { query: "http%3A%2F%2Fwww.getsidekick.com%2Fblog%2Fbody-language-advice", count: 6675 },
//         { query: "http%3A%2F%2Fwebboard.yenta4.com%2Ftopic%2F568045", count: 4652 },
//         { query: "http%3A%2F%2Fwebboard.yenta4.com%2Ftopic%2F379035%3Fsort%3D1", count: 3100 }
//       ]
//     }
// Top 5 popular queries done on Aug 2nd: GET /1/queries/popular/2015-08-02?size=5: returns
//     {
//       queries: [
//         { query: "http%3A%2F%2Fwww.getsidekick.com%2Fblog%2Fbody-language-advice", count: 2283 },
//         { query: "http%3A%2F%2Fwebboard.yenta4.com%2Ftopic%2F568045", count: 1943 },
//         { query: "http%3A%2F%2Fwebboard.yenta4.com%2Ftopic%2F379035%3Fsort%3D1", count: 1358 },
//         { query: "http%3A%2F%2Fjamonkey.com%2F50-organizing-ideas-for-every-room-in-your-house%2F", count: 890 },
//         { query: "http%3A%2F%2Fsharingis.cool%2F1000-musicians-played-foo-fighters-learn-to-fly-and-it-was-epic", count: 701 }
//       ]
//     }

func topNHandler(w http.ResponseWriter, r *http.Request, database []record, urlPrefix string) {
	// w.Header().Add("content-type", "application/json")
	// urlCount := getDistinctQueries(database, urlPrefix, r.URL.Path)
	// result, _ := json.Marshal(struct {
	// 	Queries []struct{ Query string }
	// }{
	// 	{{""}},
	// })
	// w.Write(result)
}

func main() {
	const queryCountURL string = "/1/queries/count/"
	const topNURL string = "/1/queries/count/"

	database, errorCount, lineCount := readData("hn_logs.tsv")
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		statusHandler(w, r, database, lineCount, errorCount)
	})
	http.HandleFunc(queryCountURL, func(w http.ResponseWriter, r *http.Request) {
		queryCountHandler(w, r, database, queryCountURL)
	})
	http.HandleFunc(queryCountURL, func(w http.ResponseWriter, r *http.Request) {
		topNHandler(w, r, database, topNURL)
	})
	log.Fatal(http.ListenAndServe(":8080", nil))
}
