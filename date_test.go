package main

import (
	"strings"
	"testing"
	"time"

	"github.com/boltdb/bolt"
	. "github.com/smartystreets/goconvey/convey"
)

var dateDb *bolt.DB

func init() {
	var err error

	// Open a test database
	dateDb, err = dbOpen("./testDate.db")
	if err != nil {
		panic("Failed to open test database")
	}
}

func TestNewDateFromCSV(t *testing.T) {
	g := Group{
		ID:       "testGroup",
		Name:     "testGroup",
		Location: locString,
	}
	saveGroup(dateDb, &g)

	Convey("Given a CSV row with too few columns", t, func() {
		row := []string{"testGroup", "Mon", "02:00", "1234"}

		Convey("date conversion should fail", func() {
			_, err := NewDateFromCSV(dateDb, row)
			So(err, ShouldNotBeNil)
		})
	})

	Convey("Given a CSV row of Monday, Feb 13 2016, 02:00 - 06:00", t, func() {
		row := []string{"testGroup", "2016-02-13", "02:00", "06:00", "1234"}

		Convey("The resulting Day should be congruent", func() {
			date, err := NewDateFromCSV(dateDb, row)
			So(err, ShouldBeNil)
			So(date.Group, ShouldEqual, "testGroup")
			So(date.Date.Unix(), ShouldEqual, time.Date(2016, 02, 13, 2, 0, 0, 0, loc).Unix())
			So(date.Time, ShouldEqual, 4*time.Hour)
		})
	})
}

func TestDateActiveAt(t *testing.T) {
	Convey("Given a Date of  Feb 13, 2016 with start time of 2:00 and duration 04:00 EST", t, func() {
		date := Date{
			Group:  "0",
			Target: "411",
			Date:   time.Date(2016, 02, 13, 2, 0, 0, 0, loc),
			Time:   4 * time.Hour,
		}

		Convey("When the test time is 3:00 on Feb 13, 2016 EST", func() {
			testTime := time.Date(2016, 02, 13, 3, 0, 0, 0, loc)
			Convey("Active should be true", func() {
				So(date.ActiveAt(testTime), ShouldBeTrue)
			})
		})

		Convey("When the test time is 3:00 on Feb 13, 2016 UTC", func() {
			testTime := time.Date(2016, 02, 13, 3, 0, 0, 0, time.UTC)
			Convey("Active should be false", func() {
				So(date.ActiveAt(testTime), ShouldBeFalse)
			})
		})

		Convey("When the test time is 7:00 on Feb 13, 2016 UTC", func() {
			testTime := time.Date(2016, 02, 13, 7, 0, 0, 0, time.UTC)
			Convey("Active should be true", func() {
				So(date.ActiveAt(testTime), ShouldBeTrue)
			})
		})

		Convey("When the test time is 7:00 on Feb 13, 2016", func() {
			testTime := time.Date(2016, 02, 13, 7, 0, 0, 0, loc)
			Convey("Active should be false", func() {
				So(date.ActiveAt(testTime), ShouldBeFalse)
			})
		})

		Convey("When the test time is 3:00 on Feb 11, 2016", func() {
			testTime := time.Date(2016, 02, 11, 3, 0, 0, 0, loc)
			Convey("Active should be false", func() {
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

func TestEncodeDate(t *testing.T) {
	Convey("Given a Date of  Feb 13, 2016 with start time of 2:00 and duration 04:00", t, func() {
		date := Date{
			Group:  "0",
			Target: "411",
			Date:   time.Date(2016, 02, 13, 2, 0, 0, 0, time.UTC),
			Time:   4 * time.Hour,
		}

		Convey("Encoding that date should succeed", func() {
			buf, err := encodeDate(&date)
			So(err, ShouldBeNil)

			Convey("Decoding that date should result in a date with the same target", func() {
				var d2 Date
				So(decodeDate(buf, &d2), ShouldBeNil)
				So(d2.Target, ShouldEqual, date.Target)
			})
		})
	})
}

func TestSaveRestore(t *testing.T) {
	Convey("Given a Date of  Feb 13, 2016 with start time of 2:00 and duration 04:00", t, func() {
		date := Date{
			Group:  "0",
			Target: "411",
			Date:   time.Date(2016, 02, 13, 2, 0, 0, 0, time.UTC),
			Time:   4 * time.Hour,
		}

		Convey("When starting with no Group 0 bucket", func() {
			dateDb.Update(func(tx *bolt.Tx) error {
				return tx.DeleteBucket([]byte(date.Group))
			})

			Convey("Saving the date should work", func() {
				err := dateDb.Update(func(tx *bolt.Tx) error {
					return date.Save(tx)
				})
				So(err, ShouldBeNil)

				Convey("Loading the date should work", func() {
					dateDb.View(func(tx *bolt.Tx) error {
						b := tx.Bucket([]byte(date.Group))
						Convey("The group bucket should exist", func() {
							So(b, ShouldNotBeNil)
						})
						Convey("The group bucket should have one sub-bucket", func() {
							So(b.Stats().BucketN, ShouldEqual, 2) // 1 (self) + 1 (other)
						})
						b = b.Bucket(datesBucket)
						Convey("The dates bucket should exist", func() {
							So(b, ShouldNotBeNil)
						})
						Convey("The dates bucket should have no sub-buckets", func() {
							So(b.Stats().BucketN, ShouldEqual, 1)
						})
						Convey("The dates bucket should have only one direct key", func() {
							So(b.Stats().KeyN, ShouldEqual, 1)
							var d2 Date
							data := b.Get(date.Key())
							So(data, ShouldNotBeEmpty)

							Convey("That key should contain a date", func() {
								err = decodeDate(data, &d2)
								So(err, ShouldBeNil)

								Convey("Which has an equivalent target to the original", func() {
									So(d2.Target, ShouldEqual, date.Target)
								})
							})
						})

						return nil
					})
				})
			})
		})
	})
}

func TestDateDump(t *testing.T) {
	groupID := "testDateDump"
	g := &Group{ID: groupID}

	list := []Date{
		Date{
			Group:  groupID,
			Target: "111",
			Date:   time.Date(2016, 2, 13, 2, 0, 0, 0, time.UTC),
			Time:   4 * time.Hour,
		},
		Date{
			Group:  groupID,
			Target: "112",
			Date:   time.Date(2016, 2, 14, 2, 0, 0, 0, time.UTC),
			Time:   5 * time.Hour,
		},
		Date{
			Group:  groupID,
			Target: "113",
			Date:   time.Date(2016, 2, 15, 2, 0, 0, 0, time.UTC),
			Time:   6 * time.Hour,
		},
	}

	dateDb.Update(func(tx *bolt.Tx) error {
		return tx.DeleteBucket([]byte(groupID))
	})

	err := dateDb.Update(func(tx *bolt.Tx) error {
		for _, d := range list {
			err := d.Save(tx)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		t.Skip("Failed to write test data to bucket for TestDateDump", err)
		return
	}

	ret, err := DatesForGroup(dateDb, g)
	if err != nil {
		t.Error("Failed to get dates for group", err)
		return
	}
	if len(ret) != len(list) {
		t.Error("Dumped list did not match input list")
		return
	}
}
