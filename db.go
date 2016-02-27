package main

import (
	"errors"
	"os"
	"path"

	"github.com/boltdb/bolt"
	"github.com/labstack/echo"
)

// ErrNotFound indicates that the entity
// was not found (in the database).
var ErrNotFound = errors.New("Not Found")

func dbOpen(f string) (handle *bolt.DB, err error) {
	// Ensure database path exists
	if err = os.MkdirAll(path.Dir(f), 0770); err != nil {
		return
	}

	// Create or Open the database file
	if handle, err = bolt.Open(f, 0660, nil); err != nil {
		return
	}

	// Make sure the buckets exist
	handle.Update(func(tx *bolt.Tx) error {
		if _, err = tx.CreateBucketIfNotExists(groupBucket); err != nil {
			return err
		}
		if _, err = tx.CreateBucketIfNotExists(daysBucket); err != nil {
			return err
		}
		if _, err = tx.CreateBucketIfNotExists(datesBucket); err != nil {
			return err
		}
		return nil
	})

	return
}

// dbFromContext returns the database pointer from
// the echo.Context.
func dbFromContext(ctx *echo.Context) *bolt.DB {
	return ctx.Get("db").(*bolt.DB)
}
