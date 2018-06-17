package main

import "testing"

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

func Test_readData(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name           string
		args           args
		wantLineCount  int
		wantErrorCount int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, gotLineCount, gotErrorCount := readData(tt.args.path)
			if gotLineCount != tt.wantLineCount {
				t.Errorf("readData() gotLineCount = %v, want %v", gotLineCount, tt.wantLineCount)
			}
			if gotErrorCount != tt.wantErrorCount {
				t.Errorf("readData() gotErrorCount = %v, want %v", gotErrorCount, tt.wantErrorCount)
			}
		})
	}
}
