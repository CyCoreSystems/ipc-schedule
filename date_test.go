package main

import (
	"strings"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

/*
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
*/

func TestNewDateFromCSV(t *testing.T) {
	Convey("Given a CSV row with too few columns", t, func() {
		row := []string{"testGroup", "Mon", "02:00", "1234"}

		Convey("date conversion should fail", func() {
			_, err := NewDateFromCSV(row)
			So(err, ShouldNotBeNil)
		})
	})

	Convey("Given a CSV row of Monday, Feb 13 2016, 02:00 - 06:00", t, func() {
		row := []string{"testGroup", "2016-02-13", "02:00", "06:00", "1234"}

		Convey("The resulting Day should be congruent", func() {
			date, err := NewDateFromCSV(row)
			So(err, ShouldBeNil)
			So(date.Group, ShouldEqual, "testGroup")
			So(date.Date.Unix(), ShouldEqual, time.Date(2016, 02, 13, 2, 0, 0, 0, time.UTC).Unix())
			So(date.Time, ShouldEqual, 4*time.Hour)
		})
	})
}

func TestDateActiveAt(t *testing.T) {
	Convey("Given a Date of  Feb 13, 2016 with start time of 2:00 and duration 04:00", t, func() {
		date := Date{
			Group:  "0",
			Target: "411",
			Date:   time.Date(2016, 02, 13, 2, 0, 0, 0, time.UTC),
			Time:   4 * time.Hour,
		}

		Convey("When the test time is 3:00 on Feb 13, 2016", func() {
			testTime := time.Date(2016, 02, 13, 3, 0, 0, 0, time.UTC)
			Convey("Active should be true", func() {
				So(date.ActiveAt(testTime), ShouldBeTrue)
			})
		})

		Convey("When the test time is 7:00 on Feb 13, 2016", func() {
			testTime := time.Date(2016, 02, 13, 7, 0, 0, 0, time.UTC)
			Convey("Active should be true", func() {
				So(date.ActiveAt(testTime), ShouldBeFalse)
			})
		})

		Convey("When the test time is 3:00 on Feb 11, 2016", func() {
			testTime := time.Date(2016, 02, 11, 3, 0, 0, 0, time.UTC)
			Convey("Active should be true", func() {
				So(date.ActiveAt(testTime), ShouldBeFalse)
			})
		})
	})
}

func TestDateToExternal(t *testing.T) {
	Convey("Given a Date of  Feb 13, 2016 with start time of 2:00 and duration 04:00", t, func() {
		date := Date{
			Group:  "0",
			Target: "411",
			Date:   time.Date(2016, 02, 13, 2, 0, 0, 0, time.UTC),
			Time:   4 * time.Hour,
		}

		Convey("When the date is exported", func() {
			e := date.ToExternal()

			Convey("The group should be the same", func() {
				So(e.Group, ShouldEqual, date.Group)
			})

			Convey("The target should be the same", func() {
				So(e.Target, ShouldEqual, date.Target)
			})

			Convey("The date should be Feb 13, 2016", func() {
				So(e.Date, ShouldEqual, "2016-02-13")
			})

			Convey("The start time should be 2:00", func() {
				So(e.Target, ShouldEqual, date.Target)
				pieces := strings.Split(e.Start, ":")
				So(pieces, ShouldHaveLength, 2)
				So(pieces[0], ShouldEqual, "02")
				So(pieces[1], ShouldEqual, "00")
			})

			Convey("The stop time should be 6:00", func() {
				pieces := strings.Split(e.Stop, ":")
				So(pieces, ShouldHaveLength, 2)
				So(pieces[0], ShouldEqual, "06")
				So(pieces[1], ShouldEqual, "00")
			})
		})
	})
}
