package main

import (
	"fmt"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"
	"time"
)

func pt(s string) time.Time {
	t, _ := time.Parse("2006-01-02 15:04:05", s)
	return t.UTC()
}

func strP(s string) *string {
	return &s
}

// BenchmarkReadData-8   	       1	1118722055 ns/op	358121792 B/op	 4870361 allocs/op
// BenchmarkReadData-8   	       1	1119643516 ns/op	358123712 B/op	 4870363 allocs/op 1.145
// BenchmarkReadData-8   	       1	1167624803 ns/op	358123648 B/op	 4870362 allocs/op 1.196
// BenchmarkReadData-8   	       1	1195088345 ns/op	358125696 B/op	 4870366 allocs/op 1.222

// BenchmarkReadData-8   	       1	1121361311 ns/op	462216960 B/op	 4870364 allocs/op 1.145s
// BenchmarkReadData-8   	       1	1174937320 ns/op	462216768 B/op	 4870361 allocs/op 1.200
// BenchmarkReadData-8   	       1	1181765164 ns/op	462218816 B/op	 4870365 allocs/op 1.209
// BenchmarkReadData-8   	       1	1208232719 ns/op	462218880 B/op	 4870366 allocs/op 1.235

// macbookpro
// BenchmarkReadData-8   	       1	1351746437 ns/op	462216832 B/op	 4870362 allocs/op
// BenchmarkReadData-8   	       1	1335289881 ns/op	462216640 B/op	 4870359 allocs/op
// BenchmarkReadData-8   	       1	1267709804 ns/op	462216640 B/op	 4870359 allocs/op
// BenchmarkReadData-8   	       1	1339727746 ns/op	462216832 B/op	 4870362 allocs/op

// BenchmarkReadData-8   	       1	1621703140 ns/op	665249584 B/op	 4870366 allocs/op
// BenchmarkReadData-8   	       1	1739551230 ns/op	665247472 B/op	 4870361 allocs/op
// BenchmarkReadData-8   	       1	1657654081 ns/op	665247568 B/op	 4870362 allocs/op
// BenchmarkReadData-8   	       1	1689174244 ns/op	665247472 B/op	 4870361 allocs/op
// BenchmarkReadData-8   	       1	1667111939 ns/op	665247472 B/op	 4870361 allocs/op

func BenchmarkReadData(b *testing.B) {
	for i := 0; i < b.N; i++ {
		database, errorCount, lineCount := readData("hn_logs.tsv")
		fmt.Print(len(database), errorCount, lineCount)
	}
}

func Test_parseRecord(t *testing.T) {
	type args struct {
		line []string
	}
	tests := []struct {
		name    string
		args    args
		wantR   *record
		wantErr bool
	}{
		{
			"ignores empty line",
			args{[]string{""}},
			nil,
			true,
		},
		{
			"ignores date only line line",
			args{[]string{"2015-08-01 00:01:16"}},
			nil,
			true,
		},
		{
			"decodes a line",
			args{[]string{"2006-01-02 15:04:05", "http%3A%2F%2Fblog.thiagorodrigo.com.br%2Fcupom-desconto-natue"}},
			&record{pt("2006-01-02 15:04:05"), "http%3A%2F%2Fblog.thiagorodrigo.com.br%2Fcupom-desconto-natue"},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotR, err := parseRecord(tt.args.line)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseRecord() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotR, tt.wantR) {
				t.Errorf("parseRecord() = %v, want %v", gotR, tt.wantR)
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
		wantR          []record
		wantErrorCount int
		wantLineCount  int
	}{
		{
			"parses a simple file and reports errors",
			args{filepath.Dir(filename) + "/test_small.tsv"},
			[]record{
				{pt("2015-08-01 00:03:40"), "go.co"},
				{pt("2015-08-01 00:03:41"), "google.com"},
				{pt("2015-08-01 00:03:42"), "facebook.com"},
				{pt("2015-08-01 00:03:43"), "google.com"},
			},
			2,
			6,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotR, gotErrorCount, gotLineCount := readData(tt.args.path)
			if !reflect.DeepEqual(gotR, tt.wantR) {
				t.Errorf("readData() gotR = %v, want %v", gotR, tt.wantR)
			}
			if gotErrorCount != tt.wantErrorCount {
				t.Errorf("readData() gotErrorCount = %v, want %v", gotErrorCount, tt.wantErrorCount)
			}
			if gotLineCount != tt.wantLineCount {
				t.Errorf("readData() gotLineCount = %v, want %v", gotLineCount, tt.wantLineCount)
			}
		})
	}
}
