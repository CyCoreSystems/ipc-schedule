package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"sort"
	"time"

	"github.com/boltdb/bolt"
)

// MinDayKey is the BoltDB key of the minimum day
var MinDayKey = []byte("0:0000")

// MaxDayKey is the BoldDB key of the maximum day
var MaxDayKey = []byte("6:1440")

// DaysBucket is the name of the Days bucket in
// BoltDB
var DaysBucket = []byte("days")

// Day represents a template schedule for a day of the week.
// Days are stored in _local_ time, for the associated group.
type Day struct {
	Group    string        // The group identifier
	Target   string        // The target number
	Day      time.Weekday  // Day or date
	Start    time.Duration // Time from 00:00
	Duration time.Duration // Length of shift

	Location string // Location for this schedule
}

// GetLocation gets the location attached to the day
func (d *Day) GetLocation() (*time.Location, error) {
	return time.LoadLocation(d.Location)
}

// TimeToDayKey returns the BoltDB "day" bucket key for the current time
func TimeToDayKey(t time.Time) []byte {
	minutes := t.Hour()*60 + t.Minute()
	return []byte(fmt.Sprintf("%d:%02.f", t.Day(), float64(minutes)))
}

// DayRangeFor returns the BoltDB "day" bucket keys for
// the start and stop range filters to find Days which
// may be applicable for the given time.
// We return -48hours and current.
func DayRangeFor(t time.Time) (from, to []byte) {
	start := t.Add(-48 * time.Hour)
	return TimeToDayKey(start), TimeToDayKey(t)
}

// ActiveDay returns the currently-active Day in the
// schedule.
func ActiveDay(db *bolt.DB, g *Group, t time.Time) *Day {
	var err error
	var d Day

	if g == nil {
		return nil
	}

	loc, err := g.GetLocation()
	if err != nil {
		return nil
	}

	// Day schedules are stored in local time, so convert
	from, to := DayRangeFor(t.In(loc))

	// generate the match func
	matchFunc := func(a, b []byte) func(tx *bolt.Tx) error {
		return func(tx *bolt.Tx) error {
			c := tx.Bucket(g.Key()).Bucket(DaysBucket).Cursor()
			for k, v := c.Seek(a); k != nil && bytes.Compare(k, b) <= 0; k, v = c.Next() {
				err = decodeDay(v, &d)
				if d.Group != g.ID {
					continue
				}
				if err != nil {
					fmt.Println("Failed to decode day", v, err)
					continue
				}
				if d.ActiveAt(t) {
					return nil
				}
			}
			return fmt.Errorf("No active day found")
		}
	}

	// Walk through the database until we find the first
	// active Day

	// If day of `from` is higher than day of `to`,
	// we have to split our search into two pieces.
	if sort.StringsAreSorted([]string{string(from[0]), string(to[0])}) {
		err = db.View(matchFunc(from, to))
	} else {
		// Search start time to end of week
		err = db.View(matchFunc(from, MaxDayKey))
		if err != nil {
			// Search start of week to end time
			err = db.View(matchFunc(MinDayKey, to))
		}
	}
	if err != nil {
		return nil
	}
	return &d
}

// Key returns the BoltDB key for this day
func (d *Day) Key() []byte {
	return []byte(fmt.Sprintf("%d:%02.f", d.Day, d.Start.Minutes()))
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

	if loc, err := d.GetLocation(); loc != nil && err == nil {
		now = now.In(loc)
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

// Save stores the Day in the database
func (d *Day) Save(db *bolt.DB) error {
	return db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(d.Group))
		if err != nil {
			return err
		}
		b, err = b.CreateBucketIfNotExists(DaysBucket)
		if err != nil {
			return err
		}
		data, err := encodeDay(d)
		if err != nil {
			return err
		}
		return b.Put(d.Key(), data)
	})
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
	if loc, err := d.GetLocation(); loc != nil && err == nil {
		t = t.In(loc)
	}
	start, stop := d.Times(t)
	return t.After(start) && t.Before(stop)
}

// NewDayFromCSVRow takes a slice of strings (from a CSV), and
// parses them into a Unit.
// Format:
//  `groupId, dayOfWeek, startTime, stopTime, cell/target`
func NewDayFromCSVRow(d []string) (*Day, error) {
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

func encodeDay(d *Day) ([]byte, error) {
	var buf bytes.Buffer
	err := gob.NewEncoder(&buf).Encode(d)
	return buf.Bytes(), err
}

func decodeDay(data []byte, d *Day) error {
	return gob.NewDecoder(bytes.NewReader(data)).Decode(d)
}
