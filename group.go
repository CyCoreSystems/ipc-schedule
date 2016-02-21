package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"

	"github.com/boltdb/bolt"
)

var groupBucket = []byte("groups")

// Group describes a group's parameters
type Group struct {
	ID       string         // group identifier
	Name     string         // group name
	Location *time.Location // Location / timezone
}

// Key returns the BoltDB keyname for the group
func (g *Group) Key() []byte {
	return []byte(g.ID)
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

func getGroup(db *bolt.DB, id string) (*Group, error) {
	var g Group
	err := db.View(func(tx *bolt.Tx) error {
		data := tx.Bucket(groupBucket).Get([]byte(id))
		if len(data) == 0 {
			return ErrNotFound
		}
		return decodeGroup(data, &g)
	})
	return &g, err
}

func saveGroup(db *bolt.DB, g Group) error {
	b, err := encodeGroup(&g)
	if err != nil {
		return err
	}
	return db.Update(func(tx *bolt.Tx) error {
		return tx.Bucket(groupBucket).Put(g.Key(), b)
	})
}

func encodeGroup(g *Group) ([]byte, error) {
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(g)
	return buf.Bytes(), err
}

func decodeGroup(data []byte, g *Group) error {
	// FIXME: gob encoding isn't working here, for some reason
	return json.NewDecoder(bytes.NewReader(data)).Decode(g)
}
