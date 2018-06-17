package main

import (
	"strings"
	"time"
)

func getDistinctQueries(database []record, urlPrefix string, requestPath string) int {
	dateString := strings.TrimPrefix(requestPath, urlPrefix)
	m := make(map[string]int)
	start, end, ok := timeRange(dateString)
	var startTs, endTs int64 = 0, 0xffffffff
	if ok {
		startTs = start.UnixNano()
		endTs = end.UnixNano()
	}

	for i := 0; i < len(database); i++ {
		t := database[i].time
		if !ok || (t >= startTs && t < endTs) {
			key := *database[i].urlP
			counter, found := m[key]
			if found {
				m[key] = counter + 1
			} else {
				m[key] = 1
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
