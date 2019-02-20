package main

import (
	"strconv"
	"time"
)

// Given a unix timestamp in milliseconds, returns the first three letters of the month
func getMonth(unixTime int64) string {
	date := time.Unix(unixTime/1000, 0)
	month := date.Month()
	return month.String()[:3]
}

// Given a unix timestamp in milliseconds, returns the day of the month (1-[28-31])
func getDate(unixTime int64) string {
	date := time.Unix(unixTime/1000, 0)
	return strconv.Itoa(date.Day())
}

// Given a unix timestamp in milliseconds, returns the first three letters of the week day.
func getWeekDay(unixTime int64) string {
	date := time.Unix(unixTime/1000, 0)
	weekDay := date.Weekday()
	return weekDay.String()[:3]
}
