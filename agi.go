package main

import (
	"net"

	"github.com/CyCoreSystems/agi"
	"github.com/boltdb/bolt"
)

func fastAGI(db *bolt.DB) {
	l, err := net.Listen("tcp", agiaddr)
	if err != nil {
		panic("Cannot listen on FastAGI address " + agiaddr)
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			Log.Error("Failed to accept AGI connection", "error", err)
			continue
		}

		Log.Debug("New AGI connection", "address", conn.RemoteAddr)
		go handleAGI(db, conn)
	}
}

func handleAGI(db *bolt.DB, c net.Conn) {
	defer c.Close()

	a := agi.New(c, c)

	exten, err := a.Get("EXTEN")
	if err != nil {
		Log.Error("Failed to get EXTEN variable from AGI", "error", err)
		return
	}

	Log.Debug("Loading target for AGI", "group", exten)
	t := getTarget(db, exten)
	err = a.Set("IPC_TARGET", t)
	if err != nil {
		Log.Error("Failed to set IPC-TARGET on AGI", "target", t, "error", err)
	}

	Log.Debug("Directed call for group to target", "group", exten, "target", t)
	return
}
