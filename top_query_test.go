package main

import (
	"fmt"
	"reflect"
	"sort"
	"testing"
)

type ByCount []record

func (a ByCount) Len() int           { return len(a) }
func (a ByCount) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByCount) Less(i, j int) bool { return a[i].count < a[j].count }

func Benchmark_agregate1(b *testing.B) {

	database, _, _ := readData("hn_logs.tsv")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		by_seccond := agregate(database, second)
		by_minute := agregate(by_seccond, minute)
		by_hour := agregate(by_minute, hour)
		by_day := agregate(by_hour, day)
		by_month := agregate(by_day, month)
		by_year := agregate(by_month, year)
		n := len(by_year)
		sort.Sort(ByCount(by_year))
		fmt.Println(n, by_year[(n-6):])
	}
}

func Test_agregate1(t *testing.T) {
	database, _, _ := readData("hn_logs.tsv")
	by_seccond := agregate(database, second)
	by_minute := agregate(by_seccond, minute)
	by_hour := agregate(by_minute, hour)
	by_day := agregate(by_hour, day)
	by_month := agregate(by_day, month)
	by_year := agregate(by_month, year)
	n := len(by_year)
	sort.Sort(ByCount(by_year))
	fmt.Println(n, by_year[(n-6):])
}

func Test_agregate(t *testing.T) {
	type args struct {
		d []record
		c dateComponent
	}
	testDb := []record{
		{pt("2006-01-02 15:04:04"), "url0", 1},
		{pt("2006-01-02 15:04:05"), "urlA", 1},
		{pt("2006-01-02 15:04:05"), "urlA", 1},
		{pt("2006-01-02 15:04:05"), "urlB", 1},
		{pt("2006-01-03 15:04:05"), "urlC", 1},
		{pt("2006-01-04 15:04:05"), "urlA", 1},
	}
	tests := []struct {
		name string
		args args
		want []record
	}{
		{
			"by seccond",
			args{testDb, second},
			[]record{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := agregate(tt.args.d, tt.args.c); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("agregate() = %v, want %v", got, tt.want)
			}
		})
	}
}
