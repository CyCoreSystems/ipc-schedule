package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// various parse functions

// Parse the provided YYYY-MM-DD date
func parseDate(d string) (time.Time, error) {
	return time.Parse("2006-01-02", d)
}

// Parse the provided day as a time.Weekday
func parseDay(d string) (w time.Weekday, err error) {
	switch strings.ToLower(d) {
	case "1", "m", "mo", "mon", "monday":
		w = time.Monday
	case "2", "t", "tu", "tue", "tuesday":
		w = time.Tuesday
	case "3", "w", "we", "wed", "wednesday":
		w = time.Wednesday
	case "4", "h", "th", "thu", "thur", "thurs", "thursday":
		w = time.Thursday
	case "5", "f", "fr", "fri", "friday":
		w = time.Friday
	case "6", "s", "sa", "sat", "saturday":
		w = time.Saturday
	case "0", "7", "u", "su", "sun", "sunday":
		w = time.Sunday
	default:
		err = fmt.Errorf("Failed to parse %s as weekday", d)
	}
	return
}

// parseTime returns the time.Duration from midnight
// of the given HH:mm-formatted time
func parseTime(src string) (dur time.Duration, err error) {
	pieces := strings.Split(src, ":")
	if len(pieces) < 2 {
		err = fmt.Errorf("Time must be formatted as HH:mm")
		return
	}

	h, err := strconv.ParseInt(pieces[0], 10, 64)
	if err != nil {
		return
	}
	dur += time.Duration(h) * time.Hour

	m, err := strconv.ParseInt(pieces[1], 10, 64)
	if err != nil {
		return
	}
	dur += time.Duration(m) * time.Minute

	if len(pieces) == 3 {
		s, err := strconv.ParseInt(pieces[2], 10, 64)
		if err != nil {
			Log.Warn("Seconds supplied but were unparseable; ignoring")
			err = nil
		} else {
			dur += time.Duration(s) * time.Second
		}
	}

	return
}
