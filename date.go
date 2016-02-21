package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"time"

	"github.com/boltdb/bolt"
)

var datesBucket = []byte("dates")

// Date represents a schedule for a specific date
type Date struct {
	Group  string        // The group identifier
	Target string        // The target number
	Date   time.Time     // Start timestamp of event
	Time   time.Duration // Duration of event
}

// TimeToDateKey returns the BoltDB "date" bucket key for the
// given time.
func TimeToDateKey(t time.Time) []byte {
	return []byte(t.String())
}

// DateRangeFor returns the BoltDB "date" bucket keys for
// the start and stop range filters to find Dates which
// may be applicable for the given time.
// We return -48hours and current.
func DateRangeFor(t time.Time) (from, to []byte) {
	start := t.Add(-48 * time.Hour)
	return TimeToDateKey(start), TimeToDateKey(t)
}

// ActiveDate returns the currently-active Date in the schedule.
func ActiveDate(db *bolt.DB, g *Group, t time.Time) *Date {
	var d Date
	var err error

	from, to := DateRangeFor(t)

	err = db.View(func(tx *bolt.Tx) error {
		c := tx.Bucket(g.Key()).Bucket(datesBucket).Cursor()
		for k, v := c.Seek(from); k != nil && bytes.Compare(k, to) <= 0; k, v = c.Next() {
			err = decodeDate(v, &d)
			if d.Group != g.ID {
				continue
			}
			if err != nil {
				fmt.Println("Failed to decode date", v, err)
				continue
			}
			if d.ActiveAt(t) {
				// Found
				return nil
			}
		}
		return fmt.Errorf("No active date found")
	})
	if err != nil {
		return nil
	}
	return &d
}

// Key returns the BoltDB key for this date
func (d *Date) Key() []byte {
	return []byte(d.Date.String())
}

// Save stores the Date in the database
func (d *Date) Save(db *bolt.DB) error {
	return db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(d.Group))
		if err != nil {
			return err
		}
		b, err = b.CreateBucketIfNotExists(datesBucket)
		if err != nil {
			return err
		}

		return b.Put(d.Key(), encodeDate(d))
	})
}

// ToExternal exports a Date schedule to its external version
func (d *Date) ToExternal() *DateExternal {
	var e DateExternal
	e.Group = d.Group
	e.Target = d.Target
	e.Date = d.Date.Format("2006-01-02")
	e.Start = d.Date.Format("15:04")
	e.Stop = d.Date.Add(d.Time).Format("15:04")
	return &e
}

// ActiveAt says whether the given time is
// within the schedule of this Date schedule
func (d *Date) ActiveAt(t time.Time) bool {
	return t.After(d.Date) && t.Before(d.Date.Add(d.Time))
}

// NewDateFromCSV takes a slice of strings (from a CSV), and
// parses them into a Unit.
// Format:
//	 `groupId`, `date`, `startTime`, `stopTime`, `cell/target`
func NewDateFromCSV(d []string) (*Date, error) {
	if len(d) != 5 {
		return nil, fmt.Errorf("CSV not in Group,Date,Time,Cell format")
	}

	e := DateExternal{
		Group:  d[0],
		Date:   d[1],
		Start:  d[2],
		Stop:   d[3],
		Target: d[4],
	}

	return e.ToDate()
}

// DateExternal represents a Date schedule unit
// suitable for import and export
type DateExternal struct {
	Group  string `json:"group"`  // Group ID
	Target string `json:"target"` // Target
	Date   string `json:"date"`   // Start date of event (YYYY-MM-DD)
	Start  string `json:"start"`  // Start time (HH:MM)
	Stop   string `json:"stop"`   // Stop Time  (HH:MM)
}

// ToDate converts an exported date schedule to a proper Date schedule
func (e *DateExternal) ToDate() (*Date, error) {
	var ret Date
	// 0: group
	ret.Group = e.Group
	if ret.Group == "" {
		return nil, fmt.Errorf("Group is mandatory")
	}

	// 1: Date (yyyy-mm-dddd)

	date, err := parseDate(e.Date)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse date: %s", err.Error())
	}
	ret.Date = date

	// 2: Start time
	start, err := parseTime(e.Start)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse time range: %s", err.Error())
	}

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

	ret.Date = ret.Date.Add(start)
	ret.Time = diff

	ret.Target = e.Target
	if ret.Target == "" {
		return nil, fmt.Errorf("Target/Cell is mandetory")
	}

	return &ret, nil

}

func encodeDate(d *Date) []byte {
	var buf bytes.Buffer
	gob.NewEncoder(&buf).Encode(d)
	return buf.Bytes()
}

func decodeDate(data []byte, d *Date) error {
	return gob.NewDecoder(bytes.NewReader(data)).Decode(d)
}
