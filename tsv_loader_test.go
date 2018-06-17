package main

import (
	"path/filepath"
	"reflect"
	"runtime"
	"testing"
	"time"
)

func Benchmark_loadData(b *testing.B) {
	for i := 0; i < b.N; i++ {
		t, _, _ := readData("hn_logs.tsv")
		t.ComputeSortedURLs()
	}
}

func Benchmark_readData(b *testing.B) {
	for i := 0; i < b.N; i++ {
		readData("hn_logs.tsv")
	}
}
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
			if !gotTrie.Equal(tt.wantTrie) {
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
