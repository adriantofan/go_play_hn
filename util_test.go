package main

import (
	"reflect"
	"testing"
	"time"
)

func TestLogDateComponents(t *testing.T) {
	type args struct {
		date time.Time
	}
	tests := []struct {
		name string
		args args
		want [6]int
	}{
		{
			"a date",
			args{time.Date(2000, time.January, 2, 11, 12, 13, 0, time.UTC)},
			[6]int{2000, 1, 2, 11, 12, 13},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := LogDateComponents(tt.args.date); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LogDateComponents() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseTime(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want time.Time
	}{
		{
			"parsable",
			args{"2000-01-02 11:12:13"},
			time.Date(2000, time.January, 2, 11, 12, 13, 0, time.UTC),
		},
		{
			"unparsable",
			args{"-01-02 11:12:13"},
			time.Time{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ParseTime(tt.args.s); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseTime() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_TimeRange(t *testing.T) {
	type args struct {
		dateString string
	}
	tests := []struct {
		name      string
		args      args
		wantStart time.Time
		wantEnd   time.Time
		wantOk    bool
	}{
		{
			"empty",
			args{""},
			time.Time{},
			time.Time{},
			false,
		},
		{
			"secconds",
			args{"2006-01-02 15:04:05"},
			time.Date(2006, 01, 02, 15, 4, 5, 0, time.UTC),
			time.Date(2006, 01, 02, 15, 4, 6, 0, time.UTC),
			true,
		},
		{
			"minutes",
			args{"2006-01-02 15:04"},
			time.Date(2006, 01, 02, 15, 4, 0, 0, time.UTC),
			time.Date(2006, 01, 02, 15, 5, 0, 0, time.UTC),
			true,
		},
		{
			"hours",
			args{"2006-01-02 15"},
			time.Date(2006, 01, 02, 15, 0, 0, 0, time.UTC),
			time.Date(2006, 01, 02, 16, 0, 0, 0, time.UTC),
			true,
		},
		{
			"days",
			args{"2006-01-02"},
			time.Date(2006, 01, 02, 0, 0, 0, 0, time.UTC),
			time.Date(2006, 01, 03, 0, 0, 0, 0, time.UTC),
			true,
		},
		{
			"months",
			args{"2006-01"},
			time.Date(2006, 1, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2006, 2, 1, 0, 0, 0, 0, time.UTC),
			true,
		},
		{
			"years",
			args{"2006"},
			time.Date(2006, 1, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2007, 1, 1, 0, 0, 0, 0, time.UTC),
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotStart, gotEnd, gotOk := TimeRange(tt.args.dateString)
			if !reflect.DeepEqual(gotStart, tt.wantStart) {
				t.Errorf("TimeRange() gotStart = %v, want %v", gotStart, tt.wantStart)
			}
			if !reflect.DeepEqual(gotEnd, tt.wantEnd) {
				t.Errorf("TimeRange() gotEnd = %v, want %v", gotEnd, tt.wantEnd)
			}
			if gotOk != tt.wantOk {
				t.Errorf("TimeRange() gotOk = %v, want %v", gotOk, tt.wantOk)
			}
		})
	}
}
