package main

import (
	"sort"
	"strings"
	"time"
)

func getDistinctQueries(database []record, urlPrefix string, requestPath string) int {
	dateString := strings.TrimPrefix(requestPath, urlPrefix)
	start, end, _ := timeRange(dateString)
	return getDistinct(database, start, end)
}

func getDistinct(database []record, start time.Time, end time.Time) int {
	infiniteInverval := false
	if end.Before(start) {
		start, end = end, start
	}
	if start.IsZero() && end.IsZero() {
		infiniteInverval = true
	}
	m := make(map[string]int)
	startIndex := sort.Search(len(database), func(i int) bool {
		return database[i].time.After(start) || database[i].time.Equal(start)
	})
	for i := startIndex; i < len(database); i++ {
		t := database[i].time
		if infiniteInverval || t == start || (t.After(start) && t.Before(end)) {
			key := database[i].url
			counter, found := m[key]
			if found {
				m[key] = counter + 1
			} else {
				m[key] = 1
			}
		} else {
			if t == end || t.After(end) {
				break
			}
		}
	}
	return len(m)
}

func timeRange(dateString string) (start time.Time, end time.Time, ok bool) {
	t, err := time.Parse("2006-01-02 15:04:05", dateString)
	if err == nil {
		return t, t.Add(time.Second), true
	}
	t, err = time.Parse("2006-01-02 15:04", dateString)
	if err == nil {
		return t, t.Add(time.Minute), true
	}
	t, err = time.Parse("2006-01-02 15", dateString)
	if err == nil {
		return t, t.Add(time.Hour), true
	}
	t, err = time.Parse("2006-01-02", dateString)
	if err == nil {
		return t, t.AddDate(0, 0, 1), true
	}
	t, err = time.Parse("2006-01", dateString)
	if err == nil {
		return t, t.AddDate(0, 1, 0), true
	}
	t, err = time.Parse("2006", dateString)
	if err == nil {
		return t, t.AddDate(1, 0, 0), true
	}
	return time.Time{}, time.Time{}, false
}
