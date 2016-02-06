package main

import (
	"bytes"
	"encoding/csv"
	"io"

	"github.com/boltdb/bolt"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/russross/blackfriday"
	"gopkg.in/inconshreveable/log15.v2"
)

var dbFile = "/var/db/ringfree/ipc.db"
var db *bolt.DB

// Log is the top-level logger
var Log log15.Logger

func main() {
	// Create a logger
	Log = log15.New()

	// Open the database
	db, err := dbOpen(dbFile)
	if err != nil {
		Log.Crit("Failed to open schedule database", "error", err)
		return
	}
	defer db.Close()

	// Create Echo web server
	e := echo.New()

	// Attach middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Attach handlers
	e.Get("/", instructions)
	e.Get("/target/:id", getTarget)
	//e.Post("/group", addGroup)

	e.Post("/sched/import/days", fileHandler(importDays))
	e.Post("/sched/import/dates", fileHandler(importDates))

	e.Run(":8080")
}

func instructions(ctx *echo.Context) error {
	instr, err := Asset("data/instructions.md")
	if err != nil {
		return err
	}
	return ctx.HTML(200, bytes.NewBuffer(blackfriday.MarkdownCommon(instr)).String())
}

func fileHandler(fn func(ctx *echo.Context, r io.Reader) error) func(ctx *echo.Context) error {
	return func(ctx *echo.Context) error {
		req := ctx.Request()

		var err error
		var input io.ReadCloser

		if h, ok := req.Header["Content-Type"]; ok {
			if h[0] == "text/csv" {
				input = req.Body
			} else {
				input, _, err = req.FormFile("file")
				if err != nil {
					return err
				}
			}
		}

		defer input.Close()

		return fn(ctx, input)
	}
}

func importDates(ctx *echo.Context, file io.Reader) error {
	r := csv.NewReader(file)

	var first bool = true
	var dates []Date
	for {
		rec, err := r.Read()
		if err != nil && err != io.EOF {
			return err
		}
		if err == io.EOF {
			break
		}

		date, err := NewDateFromCSV(rec)
		if err != nil && first {
			first = false
			continue // assumem first row is header and skip
		}

		first = false
		if err != nil {
			return err
		}

		dates = append(dates, date)
	}

	for _, date := range dates {
		if err := date.Save(); err != nil {
			return err
		}
	}

	return nil
}

func importDays(ctx *echo.Context, file io.Reader) error {
	r := csv.NewReader(file)

	var first bool = true
	var days []Day
	for {
		rec, err := r.Read()
		if err != nil && err != io.EOF {
			return err
		}
		if err == io.EOF {
			break
		}

		day, err := NewDayFromCSV(rec)
		if err != nil && first {
			first = false
			continue // assume first row is header and skip
		}

		first = false

		if err != nil {
			return err
		}

		days = append(days, day)
	}

	for _, day := range days {
		if err := day.Save(); err != nil {
			return err
		}
	}

	return nil
}

func getSchedule(ctx *echo.Context) error {
	return nil
}

// getTarget returns the target for the present time
func getTarget(ctx *echo.Context) error {
	_, err := getGroup(ctx.Param("id"))
	if err != nil {
		return err
	}

	// See if we have an explicit date entry

	return nil
}
