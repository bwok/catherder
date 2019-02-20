package main

import (
	"math"
	"testing"
)

// TODO all tests here are wrong for any time zone except NZ. Fix.

func TestGetMonth(t *testing.T) {
	var input = []struct {
		date     int64
		expected string
	}{
		{0, "Jan"},
		{-9999999999, "Sep"},
		{math.MaxInt64, "Aug"},
		{math.MinInt64, "May"},
	}

	for _, test := range input {
		retVal := getMonth(test.date)

		if retVal != test.expected {
			t.Errorf(`getMonth(%d) = %q, want: %q`, test.date, retVal, test.expected)
		}
	}
}

func TestGetDate(t *testing.T) {
	var input = []struct {
		date     int64
		expected string
	}{
		{0, "1"},
		{-9999999999, "7"},
		{math.MaxInt64, "17"},
		{math.MinInt64, "17"},
	}

	for _, test := range input {
		retVal := getDate(test.date)

		if retVal != test.expected {
			t.Errorf(`getDate(%d) = %q, want: %q`, test.date, retVal, test.expected)
		}
	}
}

func TestGetWeekDay(t *testing.T) {
	var input = []struct {
		date     int64
		expected string
	}{
		{0, "Thu"},
		{-9999999999, "Sun"},
		{math.MaxInt64, "Sun"},
		{math.MinInt64, "Mon"},
	}

	for _, test := range input {
		retVal := getWeekDay(test.date)

		if retVal != test.expected {
			t.Errorf(`getWeekDay(%d) = %q, want: %q`, test.date, retVal, test.expected)
		}
	}
}
