package main

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"
)

func Test_distinctQueryCountHandler(t *testing.T) {
	savedComputeDistinctQueryCount := config.computeDistinctQueryCount
	requestPath := "/1/queries/count/2222-01-01"
	config.computeDistinctQueryCount = func(trie Trie, path string) int {
		if trie != config.trie {
			t.Errorf("distinctQueryCountHandler should get global trie %v got %v", config.trie, trie)
		}
		if path != requestPath {
			t.Errorf("distinctQueryCountHandler should get request path %v got %v", requestPath, path)
		}
		return 10
	}
	req, err := http.NewRequest("GET", requestPath, nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(distinctQueryCountHandler)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	expected := `{"count":10}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
	config.computeDistinctQueryCount = savedComputeDistinctQueryCount
}

func TestComputeDistinctQueryCount(t *testing.T) {
	trie := MakeTrie()
	trie.AddLog(ParseTime("2015-08-01 00:03:43"), "google.com")
	trie.AddLog(ParseTime("2015-08-01 00:03:42"), "facebook.com")
	trie.AddLog(ParseTime("2015-08-01 00:03:41"), "google.com")
	trie.AddLog(ParseTime("2015-09-01 00:03:40"), "go.co")
	trie.ComputeSortedURLs()
	type args struct {
		trie Trie
		path string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			"parses components and returns valid result",
			args{
				trie,
				"/1/queries/count/2015-08-01 00:03",
			},
			2,
		},
		{
			"returns 0 when the date is un-parsable",
			args{
				trie,
				"/1/queries/count/garbage",
			},
			0,
		},
		{
			"returns 0 when there is no date",
			args{
				trie,
				"/1/queries/count/",
			},
			0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ComputeDistinctQueryCount(tt.args.trie, tt.args.path); got != tt.want {
				t.Errorf("ComputeDistinctQueryCount() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_topNHandler(t *testing.T) {
	savedComputeTopNQueries := config.computeTopNQueries
	requestPath := topNURL + "2222-01-01"
	requestParams := url.Values{"size": {"1"}}
	config.computeTopNQueries = func(trie Trie, path string, params url.Values) []QueryCountPair {
		if trie != config.trie {
			t.Errorf("topNQueriesCountHandler should get global trie %v got %v", config.trie, trie)
		}
		if path != requestPath {
			t.Errorf("topNQueriesCountHandler should get request path %v got %v", requestPath, path)
		}
		if !reflect.DeepEqual(params, requestParams) {
			t.Errorf("topNQueriesCountHandler should get request params %v got %v", requestParams, params)
		}
		return []QueryCountPair{
			QueryCountPair{"one", 100}}
	}
	req, err := http.NewRequest("GET", requestPath+"?size=1", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(topNHandler)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	expected := `{"queries":[{"query":"one","count":100}]}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
	config.computeTopNQueries = savedComputeTopNQueries
}

func TestComputeTopNQueries(t *testing.T) {
	trie := MakeTrie()
	n := trie.rootNode.getOrMake(2015)
	n.sortedUrls = &[]QueryCountPair{
		QueryCountPair{"one", 100},
		QueryCountPair{"two", 90},
		QueryCountPair{"three", 80},
		QueryCountPair{"four", 70},
		QueryCountPair{"five", 60},
		QueryCountPair{"six", 50},
	}
	type args struct {
		trie   Trie
		path   string
		params url.Values
	}
	tests := []struct {
		name string
		args args
		want []QueryCountPair
	}{
		{
			"happy path",
			args{
				trie,
				topNURL + "2015",
				url.Values{},
			},
			[]QueryCountPair{
				QueryCountPair{"one", 100},
				QueryCountPair{"two", 90},
				QueryCountPair{"three", 80},
				QueryCountPair{"four", 70},
				QueryCountPair{"five", 60},
			},
		},
		{
			"more than size",
			args{
				trie,
				topNURL + "2015",
				url.Values{"size": {"10"}},
			},
			[]QueryCountPair{
				QueryCountPair{"one", 100},
				QueryCountPair{"two", 90},
				QueryCountPair{"three", 80},
				QueryCountPair{"four", 70},
				QueryCountPair{"five", 60},
				QueryCountPair{"six", 50},
			},
		},
		{
			"less than size",
			args{
				trie,
				topNURL + "2015",
				url.Values{"size": {"2"}},
			},
			[]QueryCountPair{
				QueryCountPair{"one", 100},
				QueryCountPair{"two", 90},
			},
		},
		{
			"no date",
			args{
				trie,
				topNURL,
				url.Values{"size": {"2"}},
			},
			[]QueryCountPair{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ComputeTopNQueries(tt.args.trie, tt.args.path, tt.args.params); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ComputeTopNQueries() = %v, want %v", got, tt.want)
			}
		})
	}
}
