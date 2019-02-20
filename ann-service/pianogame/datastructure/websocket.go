package datastructure

import "github.com/gorilla/websocket"

// WebSocketUser store websocket user
type WebSocketUser struct {
	id     string
	wsconn *websocket.Conn
}
