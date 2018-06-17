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
			"un-parsable",
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

func Test_LogDateComponentsFromString(t *testing.T) {
	type args struct {
		dateString string
	}
	tests := []struct {
		name           string
		args           args
		wantComponents []int
	}{
		{
			"empty",
			args{""},
			[]int{},
		},
		{
			"seconds",
			args{"2006-01-02 15:04:05"},
			[]int{2006, 01, 02, 15, 4, 5},
		},
		{
			"minutes",
			args{"2006-01-02 15:04"},
			[]int{2006, 01, 02, 15, 4},
		},
		{
			"hours",
			args{"2006-01-02 15"},
			[]int{2006, 01, 02, 15},
		},
		{
			"days",
			args{"2006-01-02"},
			[]int{2006, 01, 02},
		},
		{
			"months",
			args{"2006-01"},
			[]int{2006, 1},
		},
		{
			"years",
			args{"2006"},
			[]int{2006},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotComponents := LogDateComponentsFromString(tt.args.dateString)
			if !reflect.DeepEqual(gotComponents, tt.wantComponents) {
				t.Errorf("TimeRange() LogDateComponentsFromString = %v, want %v", gotComponents, tt.wantComponents)
			}
		})
	}
}
