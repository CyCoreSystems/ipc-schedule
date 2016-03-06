package main

import "time"

var loc *time.Location
var locString = "US/Eastern"

func init() {
	var err error
	loc, err = time.LoadLocation(locString)
	if err != nil {
		panic("Failed to load location")
	}
}
