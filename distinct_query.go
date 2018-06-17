package main

import (
	"sort"
	"strings"
	"time"
)

func getDistinctQueries(database []record, urlPrefix string, requestPath string) int {
	dateString := strings.TrimPrefix(requestPath, urlPrefix)
	start, end, ok := timeRange(dateString)
	var startTs, endTs int64 = 0, 0xffffffff
	if ok {
		startTs = start.UnixNano()
		endTs = end.UnixNano()
	}
	return getDistinct(database, startTs, endTs)
}

func getDistinct(database []record, startTs int64, endTs int64) int {
	infiniteInverval := false
	if endTs < startTs {
		startTs, endTs = endTs, startTs
	}
	if startTs == 0 && endTs == 0xffffffff {
		infiniteInverval = true
	}
	m := make(map[string]int)
	start := sort.Search(len(database), func(i int) bool {
		return database[i].time >= startTs
	})
	for i := start; i < len(database); i++ {
		t := database[i].time
		if infiniteInverval || (t >= startTs && t < endTs) {
			key := database[i].url
			counter, found := m[key]
			if found {
				m[key] = counter + 1
			} else {
				m[key] = 1
			}
		} else {
			if t >= endTs {
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
