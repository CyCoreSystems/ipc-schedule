package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Day represents a template schedule for a day of the week
type Day struct {
	Group    string         // The group identifier
	Target   string         // The target number
	Day      time.Weekday   // Day or date
	Start    time.Duration  // Time from 00:00
	Duration time.Duration  // Length of shift
	Location *time.Location // Location for this schedule
}

// ToExternal exports a Day schedule to its external version
func (d *Day) ToExternal(now time.Time) *DayExternal {
	var e DayExternal
	start, stop := d.Times(now)
	e.Group = d.Group
	e.Target = d.Target
	e.Day = d.Day.String()
	e.Start = start.Format("15:04")
	e.Stop = stop.Format("15:04")
	return &e
}

// Times returns the start and stop time for
// the closest Day translation to now.
func (d *Day) Times(now time.Time) (closestStart time.Time, closestStop time.Time) {
	if now.IsZero() {
		Log.Error("Cannot find times without reference")
		return
	}

	if d.Location != nil {
		now = now.In(d.Location)
	}

	// find nearest start time of this slot
	if now.Weekday() == d.Day && todayAt(now, d.Start).Before(now) {
		// if the schedule is for today
		// and we are after that scheduled time,
		// then we have found the closest start
		closestStart = todayAt(now, d.Start)
	} else {
		// Step back a day at a time until we find
		// the most recent day that matches the day
		// of the week of our Day.
		for i := 1; i < 8; i++ {
			sTime := now.Add(-time.Duration(i) * 24 * time.Hour)
			if sTime.Weekday() == d.Day {
				closestStart = todayAt(sTime, d.Start)
				break
			}
		}
	}

	// closestStop is just the closests start plus
	// the duration
	closestStop = closestStart.Add(d.Duration)

	return
}

// timeSinceMidnight returns the difference in
// the given time from the previous midnight
func timeSinceMidnight(t time.Time) time.Duration {
	t = t.Truncate(time.Minute)
	return time.Duration(t.Hour())*time.Hour +
		time.Duration(t.Minute())*time.Minute
}

// timeOfLastMidnight returns the time of the most
// recent midnight
func timeOfLastMidnight(t time.Time) time.Time {
	t = t.Truncate(time.Minute)
	return t.Add(-timeSinceMidnight(t))
}

// todayAt returns the time for today
// at the given difference in time from midnight.
func todayAt(t time.Time, diff time.Duration) time.Time {
	base := timeOfLastMidnight(t)
	return base.Add(diff)
}

// ActiveAt says whether the given time is
// wihtin the schedule of this Day schedule.
func (d *Day) ActiveAt(t time.Time) bool {
	if d.Location != nil {
		t = t.In(d.Location)
	}
	start, stop := d.Times(t)
	return t.After(start) && t.Before(stop)
}

// NewDayFromCSV takes a slice of strings (from a CSV), and
// parses them into a Unit.
// Format:
//  `groupId, dayOfWeek, startTime, stopTime, cell/target`
func NewDayFromCSV(d []string) (*Day, error) {
	if len(d) != 5 {
		return nil, fmt.Errorf("CSV not in Group,Day,Time,Cell format")
	}

	e := DayExternal{
		Group:  d[0],
		Day:    d[1],
		Start:  d[2],
		Stop:   d[3],
		Target: d[4],
	}
	return e.ToDay()
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

// DayExternal represents a Day schedule unit
// suitable for import and export
type DayExternal struct {
	Group  string `json:"group"`  // Group ID
	Target string `json:"target"` // Target cell phone
	Day    string `json:"day"`    // Day of the Week
	Start  string `json:"start"`  // Start time
	Stop   string `json:"stop"`   // Stop time
}

// ToDay converts an exported day schedule to a
// proper Day schedule.
func (e *DayExternal) ToDay() (*Day, error) {
	var ret Day
	// 0: group
	ret.Group = e.Group
	if ret.Group == "" {
		return nil, fmt.Errorf("Group is mandatory")
	}
	g, err := getGroup(ret.Group)
	if err != nil {
		return nil, fmt.Errorf("Group does not exist: %s", err.Error())
	}
	// Store the group's location information to this schedule
	ret.Location = g.Location

	// 1: Day of week
	dow, err := parseDay(e.Day)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse day of the week: %s", err.Error())
	}
	ret.Day = dow

	// 2: Start time
	start, err := parseTime(e.Start)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse time range: %s", err.Error())
	}
	ret.Start = start

	// 3: Stop time
	stop, err := parseTime(e.Stop)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse time range: %s", err.Error())
	}
	diff := stop - start
	if diff < 0 {
		// If the time is negative, we wrapped to the next day
		diff += 24 * time.Hour
	}
	ret.Duration = diff

	// 3: Cell / Target
	ret.Target = e.Target
	if ret.Target == "" {
		return nil, fmt.Errorf("Target/Cell is mandatory")
	}

	return &ret, nil
}
