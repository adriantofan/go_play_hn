package main

import (
	"fmt"
	"testing"
	"time"
)

func Test_getDistinctE2E(t *testing.T) {
	type args struct {
		database  []record
		startTime time.Time
		endTime   time.Time
	}
	database, errorCount, lineCount := readData("hn_logs.tsv")
	if errorCount != 0 {
		t.Errorf("readData() errorCount = %v, want %v", errorCount, 0)
	}
	if lineCount != 1623420 {
		t.Errorf("readData() lineCount = %v, want %v", lineCount, 0)
	}
	if lineCount != len(database) {
		t.Errorf("readData() lineCount should equal len(database)")
	}
	newArgs := func(dateStr string) args {
		start, end, _ := timeRange(dateStr)
		return args{database, start, end}
	}

	tests := []struct {
		name string
		args args
		want int
	}{
		{
			"distinct 2015",
			newArgs("2015"),
			573697,
		},
		{
			"distinct 2015-08",
			newArgs("2015-08"),
			573697,
		},
		{
			"distinct 2015-08-03",
			newArgs("2015-08-03"),
			198117,
		},
		{
			"distinct 2015-08-01 00:04",
			newArgs("2015-08-01 00:04"),
			617,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getDistinct(tt.args.database, tt.args.startTime, tt.args.endTime); got != tt.want {
				t.Errorf("getDistinct() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDistinctExamples(t *testing.T) {
	database, errorCount, lineCount := readData("hn_logs.tsv")
	if errorCount != 0 {
		t.Errorf("readData() errorCount = %v, want %v", errorCount, 0)
	}
	if lineCount != 1623420 {
		t.Errorf("readData() lineCount = %v, want %v", lineCount, 0)
	}
	if lineCount != len(database) {
		t.Errorf("readData() lineCount should equal len(database)")
	}
	start, end, _ := timeRange("2015")
	distinctCount := getDistinct(database, start, end)
	if distinctCount != 573697 {
		t.Errorf("distinct 2015 should return 573697")
	}
	start, end, _ = timeRange("2015-08")
	distinctCount = getDistinct(database, start, end)
	if distinctCount != 573697 {
		t.Errorf("distinct 2015-08 should return 573697")
	}
	start, end, _ = timeRange("2015-08")
	distinctCount = getDistinct(database, start, end)
	if distinctCount != 573697 {
		t.Errorf("distinct 2015-08 should return 573697")
	}
	fmt.Println(errorCount, lineCount)
}
