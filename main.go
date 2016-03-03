package main

import (
	"encoding/csv"
	"flag"
	"io"
	"net/http"
	"time"

	"github.com/GeertJohan/go.rice"
	"github.com/boltdb/bolt"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/satori/go.uuid"
	"gopkg.in/inconshreveable/log15.v2"
)

var dbFile = "/var/db/ringfree/ipc.db"
var db *bolt.DB

// addr is the listen address
var addr string

// Log is the top-level logger
var Log log15.Logger

func init() {
	flag.StringVar(&addr, "addr", ":9000", "Address binding")
}

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
	e.Use(func(ctx *echo.Context) error {
		ctx.Set("db", db)
		return nil
	})

	// Attach handlers

	// Static content
	assetHandler := http.FileServer(rice.MustFindBox("public").HTTPBox())
	e.Get("/", func(c *echo.Context) error {
		assetHandler.ServeHTTP(c.Response().Writer(), c.Request())
		return nil
	})

	e.Get("/app/*", func(c *echo.Context) error {
		http.StripPrefix("/app/", assetHandler).
			ServeHTTP(c.Response().Writer(), c.Request())
		return nil
	})

	//e.Index("public/index.html")
	//e.Favicon("public/favicon.ico")

	//e.Static("/tags", "public/tags")
	//e.Static("/scripts", "public/scripts")

	// Data endpoints

	e.Get("/target/:id", getTarget)
	e.Get("/groups", getGroups)
	e.Post("/group", postGroup)
	e.Get("/group/:id", getGroupHandler)
	e.Delete("/group/:id", deleteGroupHandler)
	//e.Put("/group/:id", editGroup)

	// Import endpoints

	e.Post("/sched/import/days", fileHandler(importDays))
	e.Post("/sched/import/dates", fileHandler(importDates))

	e.Run(":8080")
}

func fileHandler(fn func(ctx *echo.Context, r io.Reader) error) func(ctx *echo.Context) error {
	return func(ctx *echo.Context) error {
		// Parse the attached file
		req := ctx.Request()

		var input io.Reader

		if h, ok := req.Header["Content-Type"]; ok {
			if h[0] == "text/csv" {
				i := req.Body
				defer i.Close()

				input = i
			} else {
				i, _, err := req.FormFile("file")
				if err != nil {
					return err
				}
				defer i.Close()

				input = i
			}
		}

		return fn(ctx, input)
	}
}

func importDates(ctx *echo.Context, file io.Reader) error {
	r := csv.NewReader(file)

	var count int
	for rec, err := r.Read(); err == nil; count++ {
		date, err := NewDateFromCSV(rec)
		if err != nil {
			if count > 0 {
				return err
			}
			continue // assume first row is header and skip
		}

		if err := date.Save(dbFromContext(ctx)); err != nil {
			return err
		}
	}

	return nil
}

func importDays(ctx *echo.Context, file io.Reader) error {
	r := csv.NewReader(file)

	var count int
	for rec, err := r.Read(); err == nil; count++ {
		// Convert the row to a Day
		day, err := NewDayFromCSVRow(rec)
		if err != nil {
			if count > 0 {
				return err
			}
			continue // assume first row is header and skip
		}

		// Add the location
		g, err := getGroup(dbFromContext(ctx), ctx.Param("group"))
		if err != nil {
			return err
		}
		day.Location = g.Location

		// Save the Day
		if err := day.Save(dbFromContext(ctx)); err != nil {
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
	g, err := getGroup(dbFromContext(ctx), ctx.Param("id"))
	if err != nil {
		return err
	}

	// See if we have an explicit date entry
	d := ActiveDate(dbFromContext(ctx), g, time.Now())
	if d != nil {
		return ctx.String(200, d.Target)
	}

	// Otherwise, use the day schedule
	d2 := ActiveDay(dbFromContext(ctx), g, time.Now())
	if d2 != nil {
		return ctx.String(200, d2.Target)
	}
	return ctx.String(404, "Not found")
}

func getGroups(ctx *echo.Context) error {
	list, err := allGroups(dbFromContext(ctx))
	if err != nil {
		return err
	}
	return ctx.JSON(200, list)
}

func postGroup(ctx *echo.Context) error {
	g := Group{
		ID:       ctx.Form("id"),
		Name:     ctx.Form("name"),
		Location: ctx.Form("timezone"),
	}
	if g.ID == "" {
		g.ID = uuid.NewV1().String()
	}
	return saveGroup(dbFromContext(ctx), &g)
}

func getGroupHandler(ctx *echo.Context) error {
	g, err := getGroup(dbFromContext(ctx), ctx.Param("id"))
	if err != nil {
		return err
	}
	return ctx.JSON(200, g)
}
func deleteGroupHandler(ctx *echo.Context) error {
	return deleteGroup(dbFromContext(ctx), ctx.Param("id"))
}
