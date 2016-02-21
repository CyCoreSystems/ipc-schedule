package main

import (
	"strings"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
	"gopkg.in/inconshreveable/log15.v2"
)

var loc *time.Location

func init() {
	var err error
	loc, err = time.LoadLocation("US/Eastern")
	if err != nil {
		panic("Failed to load location")
	}

	Log = log15.New()

	// Open a test database
	db, err = dbOpen("./test.db")
	if err != nil {
		panic("Failed to open test database")
	}
}

func TestTimeSinceMidnight(t *testing.T) {
	Convey("Given a time (UTC)", t, func() {
		hours := 11
		minutes := 25
		t := time.Date(2016, 01, 24, hours, minutes, 12, 0, time.UTC)

		Convey("The time since midnight should be close to the current time", func() {
			dur := timeSinceMidnight(t)
			So(dur, ShouldAlmostEqual, time.Duration(hours)*time.Hour+time.Duration(minutes)*time.Minute)
		})

	})

	Convey("Given a time (local)", t, func() {
		hours := 11
		minutes := 25
		t := time.Date(2016, 01, 24, hours, minutes, 12, 0, loc)

		Convey("The time since midnight should be close to the current time", func() {
			dur := timeSinceMidnight(t)
			So(dur, ShouldAlmostEqual, time.Duration(hours)*time.Hour+time.Duration(minutes)*time.Minute)
		})
	})

	Convey("Given a time (local) after UTC has changed days", t, func() {
		hours := 01
		minutes := 25
		t := time.Date(2016, 01, 24, hours, minutes, 12, 0, loc)

		Convey("The time since midnight should be close to the current time", func() {
			dur := timeSinceMidnight(t)
			So(dur, ShouldAlmostEqual, time.Duration(hours)*time.Hour+time.Duration(minutes)*time.Minute)
		})
	})
}

func TestTimeOfLastMidnight(t *testing.T) {
	Convey("Given a time (UTC)", t, func() {
		hours := 11
		minutes := 25
		tm := time.Date(2016, 01, 24, hours, minutes, 12, 0, time.UTC)

		Convey("The time of the last midnight should be 00:00", func() {
			mid := timeOfLastMidnight(tm)
			So(mid.Hour(), ShouldEqual, 0)
			So(mid.Minute(), ShouldEqual, 0)
		})
	})

	Convey("Given a time (local)", t, func() {
		hours := 11
		minutes := 25
		tm := time.Date(2016, 01, 24, hours, minutes, 12, 0, loc)

		Convey("The time of the last midnight should be 00:00", func() {
			mid := timeOfLastMidnight(tm)
			So(mid.Hour(), ShouldEqual, 0)
			So(mid.Minute(), ShouldEqual, 0)
		})
	})

	Convey("Given a time (local) after UTC has changed days", t, func() {
		hours := 23
		minutes := 25
		tm := time.Date(2016, 01, 24, hours, minutes, 12, 0, loc)

		Convey("The time of the last midnight should be 00:00, not more than 24 hours ago", func() {
			mid := timeOfLastMidnight(tm)
			So(mid.Hour(), ShouldEqual, 0)
			So(mid.Minute(), ShouldEqual, 0)
			So(mid, ShouldHappenAfter, tm.Add(-24*time.Hour))
		})
	})
}

func TestTodayAt(t *testing.T) {
	Convey("Given a time (UTC)", t, func() {
		hours := 11
		minutes := 25
		tm := time.Date(2016, 01, 24, hours, minutes, 12, 0, time.UTC)

		Convey("and an offset of 2 hours", func() {
			dur := 2 * time.Hour

			Convey("The time at the offset today should be 02:00", func() {
				nt := todayAt(tm, dur)
				So(nt, ShouldHappenAfter, timeOfLastMidnight(tm))
				So(nt, ShouldHappenBefore, timeOfLastMidnight(tm.Add(24*time.Hour)))
				So(nt.Hour(), ShouldEqual, 2)
				So(nt.Minute(), ShouldEqual, 0)
			})
		})
	})
}

func TestDayTimes(t *testing.T) {
	Convey("Given a Day", t, func() {
		day := Day{
			Group:    "0",
			Target:   "411",
			Day:      time.Monday,
			Start:    2 * time.Hour, // 02:00
			Duration: 4 * time.Hour, // until 06:00
			Location: loc,
		}

		Convey("When times are derived for the Zero time", func() {
			var zeroTime time.Time
			start, stop := day.Times(zeroTime)
			So(zeroTime.IsZero(), ShouldBeTrue)
			So(start.IsZero(), ShouldBeTrue)
			So(stop.IsZero(), ShouldBeTrue)
		})

		Convey("When the times are derived for time.Now()", func() {
			start, stop := day.Times(time.Now())

			Convey("The start time should be realistic", func() {
				So(start, ShouldHappenAfter, time.Now().Add(-8*24*time.Hour))
			})

			Convey("The stop time should be realistic", func() {
				So(stop, ShouldHappenAfter, time.Now().Add(-8*24*time.Hour))
			})
		})
	})
}

func TestDayActiveAt(t *testing.T) {
	Convey("Given a Day of Monday with start time 02:00 and duration 04:00", t, func() {
		day := Day{
			Group:    "0",
			Target:   "411",
			Day:      time.Monday,
			Start:    2 * time.Hour, // 02:00
			Duration: 4 * time.Hour, // until 06:00
			Location: loc,
		}

		Convey("When the test time is 03:00 on a Monday", func() {
			testTime := time.Date(2016, 01, 25, 3, 0, 0, 0, loc)
			So(testTime.Weekday(), ShouldEqual, time.Monday)

			Convey("Active should be true", func() {
				So(day.ActiveAt(testTime), ShouldBeTrue)
			})
		})

		Convey("When the test time is 07:00 on a Monday", func() {
			testTime := time.Date(2016, 01, 25, 7, 0, 0, 0, loc)
			So(testTime.Weekday(), ShouldEqual, time.Monday)

			Convey("Active should be false", func() {
				So(day.ActiveAt(testTime), ShouldBeFalse)
			})
		})

		Convey("When the test time is 03:00 on a Tuesday", func() {
			testTime := time.Date(2016, 01, 26, 3, 0, 0, 0, loc)

			Convey("Active should be false", func() {
				So(day.ActiveAt(testTime), ShouldBeFalse)
			})
		})
	})
}

func TestDayToExternal(t *testing.T) {
	Convey("Given a Day with start time 02:00 and duration 04:00", t, func() {
		day := Day{
			Group:    "0",
			Target:   "411",
			Day:      time.Monday,
			Start:    2 * time.Hour, // 02:00
			Duration: 4 * time.Hour, // until 06:00
			Location: loc,
		}

		Convey("When the day is exported", func() {
			e := day.ToExternal(time.Now())

			Convey("The group should be the same", func() {
				So(e.Group, ShouldEqual, day.Group)
			})

			Convey("The target should be the same", func() {
				So(e.Target, ShouldEqual, day.Target)
			})

			Convey("The day should be a string", func() {
				So(e.Day, ShouldNotBeBlank)
			})

			Convey("The start time should be 02:00", func() {
				pieces := strings.Split(e.Start, ":")
				So(pieces, ShouldHaveLength, 2)
				So(pieces[0], ShouldEqual, "02")
				So(pieces[1], ShouldEqual, "00")
			})

			Convey("The stop time should be formatted as HH:mm", func() {
				pieces := strings.Split(e.Stop, ":")
				So(pieces, ShouldHaveLength, 2)
				So(pieces[0], ShouldNotBeBlank)
				So(pieces[1], ShouldNotBeBlank)
			})

		})
	})
}

func TestParseDay(t *testing.T) {
	Convey("Given 'M'", t, func() {
		dow := "M"
		parsed, err := parseDay(dow)
		So(err, ShouldBeNil)
		So(parsed, ShouldEqual, time.Monday)
	})
	Convey("Given 'TU'", t, func() {
		dow := "TU"
		parsed, err := parseDay(dow)
		So(err, ShouldBeNil)
		So(parsed, ShouldEqual, time.Tuesday)
	})
	Convey("Given 'wed'", t, func() {
		dow := "wed"
		parsed, err := parseDay(dow)
		So(err, ShouldBeNil)
		So(parsed, ShouldEqual, time.Wednesday)
	})
	Convey("Given 'Thurs'", t, func() {
		dow := "Thurs"
		parsed, err := parseDay(dow)
		So(err, ShouldBeNil)
		So(parsed, ShouldEqual, time.Thursday)
	})
	Convey("Given 'FRIDAY'", t, func() {
		dow := "FRIDAY"
		parsed, err := parseDay(dow)
		So(err, ShouldBeNil)
		So(parsed, ShouldEqual, time.Friday)
	})
	Convey("Given '6'", t, func() {
		dow := "6"
		parsed, err := parseDay(dow)
		So(err, ShouldBeNil)
		So(parsed, ShouldEqual, time.Saturday)
	})
	Convey("Given 'sunday'", t, func() {
		dow := "sunday"
		parsed, err := parseDay(dow)
		So(err, ShouldBeNil)
		So(parsed, ShouldEqual, time.Sunday)
	})
	Convey("Given 'snoopday'", t, func() {
		dow := "snoopday"
		_, err := parseDay(dow)
		So(err, ShouldNotBeNil)
	})
}

func TestParseTime(t *testing.T) {
	Convey("Given 0204", t, func() {
		ref := "0204"

		Convey("The time should fail to be parsed", func() {
			_, err := parseTime(ref)
			So(err, ShouldNotBeNil)
		})

	})
	Convey("Given NaN:NaN", t, func() {
		ref := "NaN:NaN"

		Convey("The time should fail to be parsed", func() {
			_, err := parseTime(ref)
			So(err, ShouldNotBeNil)
		})

	})
	Convey("Given 00:NaN", t, func() {
		ref := "00:NaN"

		Convey("The time should fail to be parsed", func() {
			_, err := parseTime(ref)
			So(err, ShouldNotBeNil)
		})

	})
	Convey("Given 02:04", t, func() {
		ref := "02:04"

		Convey("The parsed time should be 02h04m", func() {
			dur, err := parseTime(ref)
			So(err, ShouldBeNil)
			So(dur.Minutes(), ShouldAlmostEqual, 2*60+4, 0.5)
		})
	})
	Convey("Given 02:04:13", t, func() {
		ref := "02:04:13"

		Convey("The parsed time should almost be 02h04m", func() {
			dur, err := parseTime(ref)
			So(err, ShouldBeNil)
			So(dur.Minutes(), ShouldAlmostEqual, 2*60+4, 0.5)
		})

	})
	Convey("Given 02:04:am", t, func() {
		ref := "02:04:am"

		Convey("The parsed time should almost be 02h04m", func() {
			dur, err := parseTime(ref)
			So(err, ShouldBeNil)
			So(dur.Minutes(), ShouldAlmostEqual, 2*60+4, 0.5)
		})

	})
}

func TestNewDayFromCSV(t *testing.T) {
	Convey("Given a CSV row with too few columns", t, func() {
		row := []string{"testGroup", "Mon", "02:00", "1234"}

		Convey("day conversion should fail", func() {
			_, err := NewDayFromCSVRow(row)
			So(err, ShouldNotBeNil)
		})
	})

	Convey("Given a CSV row of Monday 02:00 - 06:00", t, func() {
		row := []string{"testGroup", "Mon", "02:00", "06:00", "1234"}

		Convey("The resulting Day should be congruent", func() {
			day, err := NewDayFromCSVRow(row)
			So(err, ShouldBeNil)
			So(day.Group, ShouldEqual, "testGroup")
			So(day.Day, ShouldEqual, time.Monday)
			So(day.Start, ShouldEqual, 2*time.Hour)
		})
	})
}
