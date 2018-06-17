package main

import (
	"encoding/csv"
	"errors"
	"io"
	"log"
	"os"
	"time"
)

const dateFormat = "2006-01-02 15:04:05"

type record struct {
	time int64
	urlP *string
}

// open the file and use a scanner to read line by line
func readData() {
	f, err := os.Open("/Users/adriantofan/devel/a_test/hn_logs_short.tsv")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
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
