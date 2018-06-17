package main

import (
	"encoding/csv"
	"errors"
	"io"
	"log"
	"os"
	"time"
)

// reads, filters and sorts data
func readData(path string) (trie Trie, lineCount int, errorCount int) {
	f, err := os.Open(path)
	if err != nil {
		config.logFatal(err)
		return Trie{}, 0, 0
	}
	defer f.Close()
	trie = MakeTrie()
	errorCount, lineCount = parseTSVFile(f, func(t time.Time, url string) {
		trie.AddLog(t, url)
	})
	log.Printf("Loaded and sorted %d lines with %d errors from file %s\n", lineCount, errorCount, path)
	return
}

func parseRecord(line []string) (t time.Time, url string, err error) {
	if len(line) != 2 {
		err = errors.New("TSV line must contain two elements")
		return
	}
	t, err = time.Parse(config.dateFormat, line[0])
	if err != nil {
		return
	}
	url = line[1]
	err = nil
	return
}

func parseTSVFile(f io.Reader, handler func(time.Time, string)) (errorCount int, lineCount int) {
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
			date, url, err := parseRecord(line)
			if err != nil {
				errorCount++
				log.Println(err)
			} else {
				handler(date, url)
			}
		}
	}
}
