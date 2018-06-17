package main

import "time"

// LogDateComponents returns [Year, Month, Day, Hour, Minute, Seccond] extracted from the date
func LogDateComponents(date time.Time) []int {
	r := [...]int{
		date.Year(),
		int(date.Month()),
		date.Day(),
		date.Hour(),
		date.Minute(),
		date.Second(),
	}
	return r[:]
}
