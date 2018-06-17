package main

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"
)

func strP(s string) *string {
	return &s
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
			&record{1136214245000000000, strP("http%3A%2F%2Fblog.thiagorodrigo.com.br%2Fcupom-desconto-natue")},
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
func Test_Open(t *testing.T) {
	_, srcFilename, _, _ := runtime.Caller(0)
	filename := filepath.Dir(srcFilename) + ""
	fmt.Println(filename)
	f, err := os.Open(filename)
	if err != nil {
		fmt.Println("fail", err)
	} else {
		fmt.Println("success")
	}
	bytes := make([]byte, 100)
	_, err = f.Read(bytes)
	if err != nil {
		fmt.Println("fail read", err)
	} else {
		fmt.Println("read success")
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
				{1438387420000000000, strP("go.co")},
				{1438387421000000000, strP("google.com")},
				{1438387422000000000, strP("facebook.com")},
				{1438387423000000000, strP("google.com")},
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
