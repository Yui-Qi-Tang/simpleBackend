package http2server

import (
	"time"
	"math/rand"
	"github.com/gorilla/websocket"
)

type user struct {
	id int
	wsconn *websocket.Conn
}

func generateUserId() int {
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	return r1.Intn(1000)
}