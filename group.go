package main

import (
	"bytes"
	"encoding/gob"
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

func getGroup(id string) (g Group, err error) {
	db.View(func(tx *bolt.Tx) error {
		return decodeGroup(tx.Bucket(groupBucket).Get([]byte(id)), &g)
	})
	return

}

func saveGroup(g Group) error {
	b := encodeGroup(&g)
	return db.Update(func(tx *bolt.Tx) error {
		return tx.Bucket(groupBucket).Put(g.Key(), b)
	})
}

func encodeGroup(g *Group) []byte {
	var buf bytes.Buffer
	gob.NewEncoder(&buf).Encode(g)
	return buf.Bytes()
}

func decodeGroup(data []byte, g *Group) error {
	return gob.NewDecoder(bytes.NewReader(data)).Decode(g)
}
