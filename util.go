package main

import "time"

// LogDateComponents returns [Year, Month, Day, Hour, Minute, Seccond] extracted from the date
func LogDateComponents(date time.Time) [6]int {
	r := [...]int{
		date.Year(),
		int(date.Month()),
		date.Day(),
		date.Hour(),
		date.Minute(),
		date.Second(),
	}
	return r
}

// ParseTime parse a time from string having the default format and returns a time.Time in UTC
func ParseTime(s string) time.Time {
	t, _ := time.Parse(config.DateFormat, s)
	return t.UTC()
}

// TimeRange given a formated datestring it tries to parse it and then returns an interval arround it
func TimeRange(dateString string) (start time.Time, end time.Time, ok bool) {
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
