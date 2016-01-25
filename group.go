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
	ID       string         // group identifier
	Name     string         // group name
	Location *time.Location // Location / timezone
}

func getGroup(id string) (g Group, err error) {
	db.View(func(tx *bolt.Tx) error {
		return decodeGroup(tx.Bucket(groupBucket).Get([]byte(id)), &g)
	})
	return

}

func saveGroup(g Group) error {
	b, err := encodeGroup(&g)
	if err != nil {
		return err
	}
	return db.Update(func(tx *bolt.Tx) error {
		return tx.Bucket(groupBucket).Put([]byte(g.ID), b)
	})
}

func encodeGroup(g *Group) ([]byte, error) {
	var buf bytes.Buffer
	err := gob.NewEncoder(&buf).Encode(g)
	if err != nil {
		return nil, fmt.Errorf("Failed to encode group: %s", err.Error())
	}
	return buf.Bytes(), nil
}

func decodeGroup(data []byte, g *Group) error {
	return gob.NewDecoder(bytes.NewReader(data)).Decode(g)
}
