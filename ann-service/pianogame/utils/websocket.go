package utils

import "github.com/gorilla/websocket"

// NewWSocketUpgrader websocket upgrader
func NewWSocketUpgrader(readBuffSize, writeBuffSize int) websocket.Upgrader {
	return websocket.Upgrader{
		ReadBufferSize:  readBuffSize,
		WriteBufferSize: writeBuffSize,
	}
}
