package main

import (
	"fmt"
	"reflect"
	"testing"
	"time"
)

func Test_getDistinct(t *testing.T) {
	type args struct {
		database  []record
		startTime time.Time
		endTime   time.Time
	}
	testDb := []record{
		{pt("2006-01-02 15:04:05"), "urlA"},
		{pt("2006-01-02 15:04:05"), "urlA"},
		{pt("2006-01-02 15:04:05"), "urlB"},
		{pt("2006-01-03 15:04:05"), "urlC"},
		{pt("2006-01-04 15:04:05"), "urlD"},
	}

	tests := []struct {
		name string
		args args
		want int
	}{
		{
			"no overlap on the left ",
			args{testDb, pt("2004-01-01 01:01:00"), pt("2005-01-01 01:01:00")},
			0,
		},
		{
			"start outide and end inside",
			args{testDb, pt("2004-01-01 01:01:00"), pt("2006-01-02 16:01:00")},
			2,
		},
		{
			"start inside and end inside",
			args{testDb, pt("2006-01-02 15:04:05"), pt("2006-01-02 16:01:00")},
			2,
		},
		{
			"start inside and end outside",
			args{testDb, pt("2006-01-02 16:01:00"), pt("2007-01-01 01:01:00")},
			2,
		},
		{
			"no overlap on the right",
			args{testDb, pt("2007-01-01 01:01:00"), pt("2008-01-01 01:01:00")},
			0,
		},
		{
			"all",
			args{testDb, time.Time{}, time.Time{}},
			4,
		},
		{
			"all - start / end inversed",
			args{testDb, pt("5000-01-01 01:01:00"), time.Time{}},
			4,
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

func Test_getDistinctQueries(t *testing.T) {
	type args struct {
		database    []record
		urlPrefix   string
		requestPath string
	}

	testDb := []record{
		{pt("2006-01-02 15:04:05"), "urlA"},
		{pt("2006-01-02 15:04:05"), "urlA"},
		{pt("2006-01-02 15:04:05"), "urlB"},
		{pt("2006-01-03 15:04:05"), "urlC"},
		{pt("2006-01-04 15:04:05"), "urlD"},
	}

	tests := []struct {
		name string
		args args
		want int
	}{
		{
			"full dataset query count",
			args{
				testDb,
				"",
				"",
			},
			4,
		},
		{
			"at the begining",
			args{
				testDb,
				"",
				"2006-01-02 15:04:05",
			},
			2,
		},
		{
			"in the middle",
			args{
				testDb,
				"",
				"2006-01-03",
			},
			1,
		},
		{
			"over the end",
			args{
				testDb,
				"",
				"2006-01",
			},
			4,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getDistinctQueries(tt.args.database, tt.args.urlPrefix, tt.args.requestPath); got != tt.want {
				t.Errorf("getDistinctQueries() = %v, want %v", got, tt.want)
			}
		})
	}
}

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

// 573697       3	 483784677 ns/op	61922493 B/op	   19113 allocs/op
// 573697       3	 486602819 ns/op	61933309 B/op	   19165 allocs/op  6.245s
// 573697       2	 534955281 ns/op	61948908 B/op	   19240 allocs/op  3.757s

// 573697       3	 435576216 ns/op	61925890 B/op	   19130 allocs/op  5.992
// 573697       3	 412219861 ns/op	61946690 B/op	   19230 allocs/op  5.831
// 573697       3	 444749906 ns/op	61934834 B/op	   19173 allocs/op

// macbook pro
// 573697       3	 435760932 ns/op	61948077 B/op	   19236 allocs/op
// 573697       3	 456520327 ns/op	61919720 B/op	   19100 allocs/op
// 573697       3	 454580880 ns/op	61916392 B/op	   19084 allocs/op

// macbook pro with date.Time
// 573697       3	 339206517 ns/op	61934418 B/op	   19171 allocs/op
// 573697       3	 429861690 ns/op	61913410 B/op	   19070 allocs/op
// 573697       2	 539505379 ns/op	61925312 B/op	   19128 allocs/op
// 573697       3	 425834306 ns/op	61923949 B/op	   19120 allocs/op
// 573697       3	 459390444 ns/op	61897818 B/op	   18995 allocs/op

func BenchmarkGetDistinctQueries(b *testing.B) {
	database, _, _ := readData("hn_logs.tsv")
	var count int
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		count = getDistinctQueries(database, "", "")
	}
	fmt.Print(count)
}
