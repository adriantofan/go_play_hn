package main

import (
	"encoding/csv"
	"errors"
	"io"
	"log"
	"os"
	"sort"
	"time"
)

const dateFormat = "2006-01-02 15:04:05"

type record struct {
	time int64
	urlP *string
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
	lineCount = 0
	for {
		line, err := r.Read()
		if err == io.EOF {
			break
		}
		//FIXME: find an example which fails without EOF !!!
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
	return
}
