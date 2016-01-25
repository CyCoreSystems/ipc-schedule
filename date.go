package main

import "time"

// Date represents a schedule for a specific date
type Date struct {
	Group  string    // The group identifier
	Target string    // The target number
	Date   time.Time // date
	Time   int64     // Time, in minutes
}
