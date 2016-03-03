package main

import (
	"os"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestGroupCodec(t *testing.T) {
	Convey("Given a group", t, func() {
		loc, _ := time.LoadLocation("US/Eastern")
		tg := Group{
			ID:       "testGroupID",
			Name:     "testGroup",
			Location: loc.String(),
		}

		Convey("Encoding it should render a non-zero-length byte slice", func() {
			data, err := encodeGroup(&tg)
			So(len(data), ShouldNotEqual, 0)
			So(err, ShouldBeNil)

			Convey("Decoding that data should return the original group", func() {
				var g2 Group
				err := decodeGroup(data, &g2)
				So(err, ShouldBeNil)
				So(g2.ID, ShouldEqual, tg.ID)
				So(g2.Location, ShouldEqual, loc.String())
			})
		})
	})
}

func TestGroups(t *testing.T) {
	db, err := dbOpen("./groupTest.db")
	if err != nil {
		panic("Failed to open test database")
	}
	defer func() {
		db.Close()
		os.Remove("./groupTest.db")
	}()

	Convey("Given an empty group bucket", t, func() {
		Convey("allGroups should return an empty list", func() {
			list, err := allGroups(db)
			So(err, ShouldBeNil)
			So(list, ShouldBeEmpty)
		})
		Convey("Adding the group 'testGroup' should succeed", func() {
			loc, _ := time.LoadLocation("US/Eastern")
			tg := Group{
				ID:       "testGroupID",
				Name:     "testGroup",
				Location: loc.String(),
			}
			err := saveGroup(db, &tg)
			So(err, ShouldBeNil)
			Convey("Getting that group should succeed", func() {
				g, err := getGroup(db, tg.ID)
				So(err, ShouldBeNil)
				So(g.Name, ShouldEqual, tg.Name)
				So(g.Location, ShouldEqual, loc.String())
			})
			Convey("allGroups should return a list with a single group", func() {
				list, err := allGroups(db)
				So(err, ShouldBeNil)
				So(len(list), ShouldEqual, 1)
				So(list[0].ID, ShouldEqual, tg.ID)
			})
		})
	})
}
