package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"time"

	"github.com/boltdb/bolt"
)

var groupBucket = []byte("groups")

// Group describes a group's parameters
type Group struct {
	ID       string `json:"id"`       // group identifier
	Name     string `json:"name"`     // group name
	Location string `json:"timezone"` // Location / timezone
}

// Key returns the BoltDB keyname for the group
func (g *Group) Key() []byte {
	return []byte(g.ID)
}

// GetLocation gets the location object
func (g *Group) GetLocation() (*time.Location, error) {
	return time.LoadLocation(g.Location)
}

// ClearDays clears the day for the group
func (g *Group) ClearDays(tx *bolt.Tx) error {
	b, err := tx.CreateBucketIfNotExists([]byte(g.ID))
	if err != nil {
		return err
	}

	err = b.DeleteBucket(daysBucket)
	if err != nil && err != bolt.ErrBucketNotFound {
		return err
	}

	return nil
}

// ClearDates clears the date schedule for the group
func (g *Group) ClearDates(tx *bolt.Tx) error {
	b, err := tx.CreateBucketIfNotExists([]byte(g.ID))
	if err != nil {
		return err
	}

	err = b.DeleteBucket(datesBucket)
	if err != nil && err != bolt.ErrBucketNotFound {
		return err
	}

	return nil
}

// allGroups returns the list of all groups
func allGroups(db *bolt.DB) (list []*Group, err error) {
	list = []*Group{}
	db.View(func(tx *bolt.Tx) error {
		c := tx.Bucket(groupBucket).Cursor()
		if c == nil {
			return fmt.Errorf("No groups found")
		}
		for k, v := c.First(); k != nil; k, v = c.Next() {
			var g Group
			err = decodeGroup(v, &g)
			if err != nil {
				return err
			}
			list = append(list, &g)
		}
		return nil
	})
	return
}

func getGroup(db *bolt.DB, id string) (g *Group, err error) {
	err = db.View(func(tx *bolt.Tx) error {
		g, err = getGroupWithTx(tx, id)
		return err
	})
	return
}

func getGroupWithTx(tx *bolt.Tx, id string) (*Group, error) {
	var g Group
	data := tx.Bucket(groupBucket).Get([]byte(id))
	if len(data) == 0 {
		return &g, ErrNotFound
	}
	err := decodeGroup(data, &g)
	return &g, err
}

func saveGroup(db *bolt.DB, g *Group) error {
	b, err := encodeGroup(g)
	if err != nil {
		return err
	}
	return db.Update(func(tx *bolt.Tx) error {
		return tx.Bucket(groupBucket).Put(g.Key(), b)
	})
}

func deleteGroup(db *bolt.DB, id string) error {
	if id == "" {
		return fmt.Errorf("Cannot delete nothing")
	}
	g := &Group{ID: id}
	return db.Update(func(tx *bolt.Tx) error {
		return tx.Bucket(groupBucket).Delete(g.Key())
	})
}

func encodeGroup(g *Group) ([]byte, error) {
	var buf bytes.Buffer
	err := gob.NewEncoder(&buf).Encode(g)
	return buf.Bytes(), err
}

func decodeGroup(data []byte, g *Group) error {
	return gob.NewDecoder(bytes.NewReader(data)).Decode(g)
}
