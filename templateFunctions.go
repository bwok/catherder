package main

import (
	"strconv"
	"time"
)

func getMonth(unixTime int64) string {
	date := time.Unix(unixTime/1000, 0)
	month := date.Month()
	return month.String()[:3]
}

func getDate(unixTime int64) string {
	date := time.Unix(unixTime/1000, 0)
	return strconv.Itoa(date.Day())
}

func getWeekDay(unixTime int64) string {
	date := time.Unix(unixTime/1000, 0)
	weekDay := date.Weekday()
	return weekDay.String()[:3]
}

