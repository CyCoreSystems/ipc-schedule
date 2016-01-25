package main

import (
	"os"
	"path"

	"github.com/boltdb/bolt"
)

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
		if _, err = tx.CreateBucketIfNotExists([]byte("groups")); err != nil {
			return err
		}
		if _, err = tx.CreateBucketIfNotExists([]byte("days")); err != nil {
			return err
		}
		if _, err = tx.CreateBucketIfNotExists([]byte("dates")); err != nil {
			return err
		}
		return nil
	})

	return
}
