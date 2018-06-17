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
