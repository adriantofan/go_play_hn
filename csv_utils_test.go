package main

import (
	"io"
	"reflect"
	"strings"
	"testing"
)

const sampleFileContent string = `
2015-08-01 00:01:23	http%3A%2F%2Fgoogle.com
2015-08-01 00:02:23	http%3A%2F%2Ffacebook.com
2015-08-01 00:02:23
something
`

func Test_parseRecord(t *testing.T) {
	type args struct {
		line []string
	}
	url1 := "http%3A%2F%2Fblog.thiagorodrigo.com.br%2Fcupom-desconto-natue"
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
			&record{1136214245000000000, &url1},
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

func Test_parseCSVFile(t *testing.T) {
	url1 := "http%3A%2F%2Fgoogle.com"
	url2 := "http%3A%2F%2Ffacebook.com"
	type args struct {
		f io.Reader
	}
	tests := []struct {
		name           string
		args           args
		wantRecords    []record
		wantErrorCount int
		wantLineCount  int
	}{
		{
			"parses a simple file and reports errors",
			args{strings.NewReader(sampleFileContent)},
			[]record{
				{1438387283000000000, &url1},
				{1438387343000000000, &url2},
			},
			2,
			4,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRecords, gotErrorCount, gotLineCount := parseCSVFile(tt.args.f)
			if !reflect.DeepEqual(gotRecords, tt.wantRecords) {
				t.Errorf("parseCSVFile() gotRecords = %v, want %v", gotRecords, tt.wantRecords)
			}
			if gotErrorCount != tt.wantErrorCount {
				t.Errorf("parseCSVFile() gotErrorCount = %v, want %v", gotErrorCount, tt.wantErrorCount)
			}
			if gotLineCount != tt.wantLineCount {
				t.Errorf("parseCSVFile() gotLineCount = %v, want %v", gotLineCount, tt.wantLineCount)
			}
		})
	}
}
