package main

import (
	"reflect"
	"testing"
	"time"
)

func Test_timeRange(t *testing.T) {
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
			gotStart, gotEnd, gotOk := timeRange(tt.args.dateString)
			if !reflect.DeepEqual(gotStart, tt.wantStart) {
				t.Errorf("timeRange() gotStart = %v, want %v", gotStart, tt.wantStart)
			}
			if !reflect.DeepEqual(gotEnd, tt.wantEnd) {
				t.Errorf("timeRange() gotEnd = %v, want %v", gotEnd, tt.wantEnd)
			}
			if gotOk != tt.wantOk {
				t.Errorf("timeRange() gotOk = %v, want %v", gotOk, tt.wantOk)
			}
		})
	}
}
