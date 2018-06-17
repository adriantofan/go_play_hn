package main

import (
	"sort"
	"time"
)

type dateComponent int

const (
	second dateComponent = iota
	minute               = iota
	hour                 = iota
	day                  = iota
	month                = iota
	year                 = iota
)

// assumes d is sorted
func agregate(d []record, c dateComponent) []record {
	if len(d) <= 1 {
		return d
	}
	acc := []record{}
	current := make(map[string]uint)
	low, high := bounds(d[0].time, c)
	current[d[0].url] = d[0].count
	for i := 1; i < len(d); i++ {
		if d[i].time == low || d[i].time.Before(high) && d[i].time.After(low) {
			key := d[i].url
			counter, found := current[key]
			if found {
				current[key] = counter + d[i].count
			} else {
				current[key] = d[i].count
			}
		} else {
			currentAcc := []record{}
			for key, value := range current {
				currentAcc = append(currentAcc, record{low, key, value})
			}
			sort.Sort(ByTime(currentAcc))
			acc = append(acc, currentAcc...)
			current = make(map[string]uint)
			low, high = bounds(d[i].time, c)
			current[d[i].url] = d[i].count
		}
	}
	currentAcc := []record{}
	for key, value := range current {
		currentAcc = append(currentAcc, record{low, key, value})
	}
	sort.Sort(ByTime(currentAcc))
	acc = append(acc, currentAcc...)
	return acc
}

func bounds(d time.Time, c dateComponent) (low time.Time, high time.Time) {
	switch c {
	case second:
		return d, d.Add(time.Second)
	case minute:
		h := d.Add(time.Minute)
		return time.Date(d.Year(), d.Month(), d.Day(), d.Hour(), d.Minute(), 0, 0, d.Location()), time.Date(h.Year(), h.Month(), h.Day(), h.Hour(), h.Minute(), 0, 0, d.Location())
	case hour:
		h := d.Add(time.Hour)
		return time.Date(d.Year(), d.Month(), d.Day(), d.Hour(), 0, 0, 0, d.Location()), time.Date(h.Year(), h.Month(), h.Day(), h.Hour(), 0, 0, 0, d.Location())
	case day:
		h := d.AddDate(0, 0, 1)
		return time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, d.Location()), time.Date(h.Year(), h.Month(), h.Day(), 0, 0, 0, 0, d.Location())
	case month:
		h := d.AddDate(0, 1, 0)
		return time.Date(d.Year(), d.Month(), 1, 0, 0, 0, 0, d.Location()), time.Date(h.Year(), h.Month(), 1, 0, 0, 0, 0, d.Location())
	case year:
		h := d.AddDate(0, 1, 0)
		return time.Date(d.Year(), 1, 1, 0, 0, 0, 0, d.Location()), time.Date(h.Year(), 1, 1, 0, 0, 0, 0, d.Location())
	}
	return d, d
}
