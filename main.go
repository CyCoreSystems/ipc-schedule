package main

import (
	"bytes"

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

	e.Run(":8080")
}

func instructions(ctx *echo.Context) error {
	instr, err := Asset("data/instructions.md")
	if err != nil {
		return err
	}
	return ctx.HTML(200, bytes.NewBuffer(blackfriday.MarkdownCommon(instr)).String())
}

func getSchedule(ctx *echo.Context) error {
	return nil
}

// getTarget returns the target for the present time
func getTarget(ctx *echo.Context) error {
	g, err := getGroup(ctx.Param("id"))
	if err != nil {
		return err
	}

	// See if we have an explicit date entry

}
