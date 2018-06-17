package main

import (
	"reflect"
	"sync"
	"testing"
)

var (
	trie       Trie
	lineCount  int
	errorCount int
	loadOnce   sync.Once
)

func loadRealData() (Trie, int, int) {
	loadOnce.Do(func() {
		trie, lineCount, errorCount = readData("hn_logs.tsv")
		trie.ComputeSortedURLs()
	})
	return trie, lineCount, errorCount

}

func Test_readDataCount(t *testing.T) {
	_, lineCount, errorCount := loadRealData()
	if errorCount != 0 {
		t.Errorf("readData() errorCount = %v, want %v", errorCount, 0)
	}
	if lineCount != 1623420 {
		t.Errorf("readData() lineCount = %v, want %v", lineCount, 0)
	}
}
func Test_DistinctValidation(t *testing.T) {
	trie, _, _ := loadRealData()
	type args struct {
		t Trie
		c []int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			"distinct 2015",
			args{trie, []int{2015}},
			573697,
		},
		{
			"distinct 2015-08",
			args{trie, []int{2015, 8}},
			573697,
		},
		{
			"distinct 2015-08-03",
			args{trie, []int{2015, 8, 3}},
			198117,
		},
		{
			"distinct 2015-08-01 00:04",
			args{trie, []int{2015, 8, 1, 0, 4}},
			617,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Distinct(tt.args.t, tt.args.c); got != tt.want {
				t.Errorf("Distinct() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_TopNAtDate(t *testing.T) {
	trie, _, _ := loadRealData()
	type args struct {
		t Trie
		c []int
		n int
	}
	tests := []struct {
		name string
		args args
		want []urlCountPair
	}{
		{
			"topn 2015",
			args{trie, []int{2015}, 3},
			[]urlCountPair{
				{"http%3A%2F%2Fwww.getsidekick.com%2Fblog%2Fbody-language-advice", 6675},
				{"http%3A%2F%2Fwebboard.yenta4.com%2Ftopic%2F568045", 4652},
				{"http%3A%2F%2Fwebboard.yenta4.com%2Ftopic%2F379035%3Fsort%3D1", 3100},
			},
		},
		{
			"distinct 2015-08-02",
			args{trie, []int{2015, 8}, 5},
			[]urlCountPair{
				{"http%3A%2F%2Fwww.getsidekick.com%2Fblog%2Fbody-language-advice", 6675},
				{"http%3A%2F%2Fwebboard.yenta4.com%2Ftopic%2F568045", 4652},
				{"http%3A%2F%2Fwebboard.yenta4.com%2Ftopic%2F379035%3Fsort%3D1", 3100},
				{"http%3A%2F%2Fjamonkey.com%2F50-organizing-ideas-for-every-room-in-your-house%2F", 2800},
				{"chrome%3A%2F%2Fnewtab%2F", 2535},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := TopNAtDate(tt.args.t, tt.args.c, tt.args.n)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TopNAtDate() = %v, want %v", got, tt.want)
			}
			gotCount := len(got)
			if gotCount != tt.args.n {
				t.Errorf("TopNAtDate count = %v, want %v", gotCount, tt.args.n)
			}
		})
	}
}
