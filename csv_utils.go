package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"time"
)

const dateFormat = "2006-01-02 15:04:05"

// BUG(adrian) is time64 good enough !? . urlP can be just url ... because it is a struct with a reference
type record struct {
	time int64
	urlP *string
}

func (r record) String() string {
	date := time.Unix(0, r.time)
	dateStr := date.UTC().Format(dateFormat)
	return fmt.Sprintf("{%d, %s, %s}", r.time, dateStr, *r.urlP)
}

type ByTime []record

func (a ByTime) Len() int           { return len(a) }
func (a ByTime) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByTime) Less(i, j int) bool { return a[i].time < a[j].time }

// reads, filters and sorts data
func readData(path string) (r []record, errorCount int, lineCount int) {
	r = nil
	f, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer f.Close()
	r, errorCount, lineCount = parseCSVFile(f)
	sort.Sort(ByTime(r))
	log.Printf("Loaded and sorted %d lines with %d errors from file %s\n", lineCount, errorCount, path)
	return
}

func parseRecord(line []string) (r *record, err error) {
	if len(line) != 2 {
		err = errors.New("TSV line must contain two elements")
		r = nil
		return
	}
	t, err := time.Parse(dateFormat, line[0])
	if err != nil {
		r = nil
		return
	}
	r = &record{t.UnixNano(), &line[1]}
	err = nil
	return
}

func parseCSVFile(f io.Reader) (records []record, errorCount int, lineCount int) {
	records = make([]record, 0)
	errorCount = 0
	r := csv.NewReader(f)
	r.Comma = '\t'
	r.FieldsPerRecord = 2
	lineCount = 0
	for {
		line, err := r.Read()
		if err != nil {
			if err == io.EOF {
				return
			}
			if _, ok := err.(*csv.ParseError); !ok {
				panic(err)
			}

		}
		lineCount++
		if err != nil {
			errorCount++
			log.Println(err)

		} else {
			r, err := parseRecord(line)
			if err != nil {
				errorCount++
				log.Println(err)
			} else {
				records = append(records, *r)
			}
		}
	}
}
