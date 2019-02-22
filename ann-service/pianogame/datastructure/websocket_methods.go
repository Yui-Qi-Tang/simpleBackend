package datastructure

import (
	"simpleBackend/ann-service/pianogame/msg"

	"github.com/gorilla/websocket"
)

// GetID returns WebSocketUser's
func (w *WebSocketUser) GetID() string {
	return w.id
}

// GetConn returns WebSocketUser's
func (w *WebSocketUser) GetConn() *websocket.Conn {
	return w.wsconn
}

// SetID set user ID
func (w *WebSocketUser) SetID(id string) {
	w.id = id
}

// SetWsConn set websocket conn for user
func (w *WebSocketUser) SetWsConn(conn *websocket.Conn) {
	w.wsconn = conn
}

// Close close ws con
func (w *WebSocketUser) Close() {
	w.wsconn.Close()
}

// SendMsg send message by message type
func (w *WebSocketUser) SendMsg(message interface{}) {
	/*
		Message category:
			1. piano key
			2. welcom message
	*/
	// need to be check all type? must be survey a better solution
	switch v := message.(type) {
	case *msg.Welcome:
		w.wsconn.WriteJSON(v)
	case msg.Welcome:
		w.wsconn.WriteJSON(v)
	case *msg.PianoKey:
		// number of msg = n-1
		if v.From != w.id {
			w.wsconn.WriteJSON(v)
		}
	case msg.PianoKey:
		// number of msg = n-1
		if v.From != w.id {
			w.wsconn.WriteJSON(v)
		}
	default:
		w.wsconn.WriteJSON(&msg.Error{Text: "Unknow msg structure"})
	}
}
