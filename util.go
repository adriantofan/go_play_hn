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

// LogDateComponentsFromString parses a date string to it's components
func LogDateComponentsFromString(dateString string) []int {
	t, err := time.Parse("2006-01-02 15:04:05", dateString)
	if err == nil {
		cs := LogDateComponents(t)
		return cs[:]
	}
	t, err = time.Parse("2006-01-02 15:04", dateString)
	if err == nil {
		cs := LogDateComponents(t)
		return cs[:5]
	}
	t, err = time.Parse("2006-01-02 15", dateString)
	if err == nil {
		cs := LogDateComponents(t)
		return cs[:4]
	}
	t, err = time.Parse("2006-01-02", dateString)
	if err == nil {
		cs := LogDateComponents(t)
		return cs[:3]
	}
	t, err = time.Parse("2006-01", dateString)
	if err == nil {
		cs := LogDateComponents(t)
		return cs[:2]
	}
	t, err = time.Parse("2006", dateString)
	if err == nil {
		cs := LogDateComponents(t)
		return cs[:1]
	}
	return []int{}
}
