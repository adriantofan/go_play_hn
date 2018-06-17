package main

import (
	"path/filepath"
	"reflect"
	"runtime"
	"testing"
	"time"
)

// Benchmark_loadData-8   	       1	4111254434 ns/op	1207746736 B/op	 3511168 allocs/op
// Benchmark_loadData-8   	       1	4011234279 ns/op	1207763984 B/op	 3511256 allocs/op

// After not keeping years cumulative:

// Benchmark_loadData-8   	       1	3627594538 ns/op	1074935344 B/op	 3491303 allocs/op
// Benchmark_loadData-8   	       1	3679311550 ns/op	1074951360 B/op	 3491366 allocs/op

func Benchmark_loadData(b *testing.B) {
	for i := 0; i < b.N; i++ {
		t, _, _ := readData("hn_logs.tsv")
		t.ComputeSortedURLs()
	}
}

// Benchmark_readData-8   	       1	3421675559 ns/op	695082704 B/op	 3446286 allocs/op
// Benchmark_readData-8   	       1	3530647423 ns/op	695029824 B/op	 3446027 allocs/op

// After not keeping years cumulative:

//  Benchmark_readData-8   	       1	2954735302 ns/op	633087712 B/op	 3426824 allocs/op
// Benchmark_readData-8   	       1	3154488192 ns/op	633076816 B/op	 3426775 allocs/op

func Benchmark_readData(b *testing.B) {
	for i := 0; i < b.N; i++ {
		readData("hn_logs.tsv")
	}
}

// 2000	    559765 ns/op	  408466 B/op	    4786 allocs/op
// 3000	    470646 ns/op	  323041 B/op	    4776 allocs/op

// After not keeping years cumulative (doesen't seem very precise):

// 2000	    544730 ns/op	  373097 B/op	    4786 allocs/op
// 2000	    575863 ns/op	  373097 B/op	    4786 allocs/op

func Benchmark_computeTops(b *testing.B) {
	t, _, _ := readData("hn_logs.tsv")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		t.ComputeSortedURLs()
	}
}

func Test_parseRecord(t *testing.T) {
	type args struct {
		line []string
	}
	tests := []struct {
		name    string
		args    args
		wantT   time.Time
		wantURL string
		wantErr bool
	}{
		{
			"ignores empty line",
			args{[]string{""}},
			time.Time{},
			"",
			true,
		},
		{
			"ignores unparsable date",
			args{[]string{"2015--01 00:01:16", "url"}},
			time.Time{},
			"",
			true,
		},
		{
			"ignores date only line line",
			args{[]string{"2015-08-01 00:01:16"}},
			time.Time{},
			"",
			true,
		},
		{
			"decodes a line",
			args{[]string{"2006-01-02 15:04:05", "http%3A%2F%2Fblog.thiagorodrigo.com.br%2Fcupom-desconto-natue"}},
			ParseTime("2006-01-02 15:04:05"),
			"http%3A%2F%2Fblog.thiagorodrigo.com.br%2Fcupom-desconto-natue",
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotT, gotURL, err := parseRecord(tt.args.line)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseRecord() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotT, tt.wantT) {
				t.Errorf("parseRecord() gotT = %v, want %v", gotT, tt.wantT)
			}
			if gotURL != tt.wantURL {
				t.Errorf("parseRecord() gotUrl = %v, want %v", gotURL, tt.wantURL)
			}
		})
	}
}

func Test_readData(t *testing.T) {
	type args struct {
		path string
	}
	_, filename, _, _ := runtime.Caller(0)
	tests := []struct {
		name           string
		args           args
		wantTrie       Trie
		wantLineCount  int
		wantErrorCount int
		wantFatal      bool
	}{
		{
			"non existent file",
			args{"mising"},
			Trie{},
			0,
			0,
			true,
		},
		{
			"parses a simple file and reports errors",
			args{filepath.Dir(filename) + "/test_small.tsv"},
			makeTestSmallTrie(),
			6,
			2,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			oldLogFatal := config.logFatal
			gotFatal := false
			config.logFatal = func(v ...interface{}) {
				gotFatal = true
			}
			gotTrie, gotLineCount, gotErrorCount := readData(tt.args.path)
			if !reflect.DeepEqual(gotTrie, tt.wantTrie) {
				t.Errorf("readData() gotTrie = %v, want %v", gotTrie, tt.wantTrie)
			}
			if gotLineCount != tt.wantLineCount {
				t.Errorf("readData() gotLineCount = %v, want %v", gotLineCount, tt.wantLineCount)
			}
			if gotErrorCount != tt.wantErrorCount {
				t.Errorf("readData() gotErrorCount = %v, want %v", gotErrorCount, tt.wantErrorCount)
			}
			if gotFatal != tt.wantFatal {
				t.Errorf("readData() gotFatal = %v, want %v", gotFatal, tt.wantFatal)
			}
			config.logFatal = oldLogFatal
		})
	}
}

// makes a trie that looks like test small
func makeTestSmallTrie() Trie {
	t := MakeTrie()
	t.AddLog(ParseTime("2015-08-01 00:03:43"), "google.com")
	t.AddLog(ParseTime("2015-08-01 00:03:42"), "facebook.com")
	t.AddLog(ParseTime("2015-08-01 00:03:41"), "google.com")
	t.AddLog(ParseTime("2015-08-01 00:03:40"), "go.co")
	return t
}
