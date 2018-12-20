package http2server

import (
	"fmt"
	"log"
	"net/http" // "io/ioutil"

	"github.com/gorilla/websocket"
)

const indexHTML = `<html>
<head>
	<title>Hello World</title>
	<script src="/static/app.js"></script>
	<link rel="stylesheet" href="/static/style.css"">
</head>
<body>
Hello, gopher!
</body>
</html>
`

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var clients = make(map[*user]bool)

var flag = make(chan bool)

// CreateServer : Set server and start
func CreateServer(config *ServerConfig) {
	log.Printf("start server at %s\n", config.Port)
	// Set static file
	// filter   "prefix" of the request and leave
	//            request                          "prefix"                      file here
	http.Handle("/chatroom/", http.StripPrefix("/chatroom/", http.FileServer(http.Dir(config.Static))))
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(config.Static))))
	// functions
	StartService()
	log.Fatal(
		http.ListenAndServeTLS(config.Port, config.Crt, config.Key, nil),
	)
	//http.ListenAndServe(config.Port, nil)

}

// StartService : set router for the server
func StartService() {
	http.HandleFunc("/", index)
	http.HandleFunc("/test", test)
	http.HandleFunc("/socket/handler", webSocketHandle)
}

func test(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello world!")
}

func index(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	fmt.Println(r.UserAgent())

	pusher, ok := w.(http.Pusher)
	if ok {
		// Push is supported. Try pushing rather than
		// waiting for the browser request these static assets.
		if err := pusher.Push("/static/app.js", nil); err != nil {
			log.Printf("Failed to push: %v", err)
		}
		if err := pusher.Push("/static/style.css", nil); err != nil {
			log.Printf("Failed to push: %v", err)
		}
		if err := pusher.Push("/static/test.js", nil); err != nil {
			log.Printf("Failed to push: %v", err)
		}
	}
	fmt.Fprintf(w, indexHTML)
}

func webSocketHandle(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Origin") != "https://"+r.Host {
		http.Error(w, "Origin not allowed", 403)
		return
	}

	pusher, ok := w.(http.Pusher)
	fmt.Println("pusher status: ", ok)
	if ok {
		options := &http.PushOptions{
			Header: http.Header{
				"Accept-Encoding": r.Header["Accept-Encoding"],
			},
		}
		fmt.Println("Push ok!")
		// Push is supported. Try pushing rather than
		// waiting for the browser request these static assets.
		if err := pusher.Push("/static/app.js", nil); err != nil {
			log.Printf("Failed to push: %v", err)
		}
		if err := pusher.Push("/static/style.css", nil); err != nil {
			log.Printf("Failed to push: %v", err)
		}
		if err := pusher.Push("/static/test.js", options); err != nil {
			log.Printf("Failed to push: %v", err)
		}
	}
	// establish web socket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	// allocate user id
	newUserID := generateUserId()
	for k := range clients {
		if k.id == newUserID {
			newUserID = generateUserId()
			continue
		}
	}

	sendFirstJoinMsg(conn, newUserID)
	// add new user
	newUser := &user{wsconn: conn, id: newUserID}
	clients[newUser] = true
	// go echo(conn)
	go chatHandle(newUser)
}

func sendFirstJoinMsg(conn *websocket.Conn, guessID int) {
	welcome := &msg{Text: "Hello!!Wellcome join us!!", MyId: guessID, To: nil, From: nil}
	conn.WriteJSON(welcome)
}

func chatHandle(chater *user) {
	for {
		m := msg{} // custom msg
		err := chater.wsconn.ReadJSON(&m)
		if err != nil {
			fmt.Println("Error reading json.", err)
			chater.wsconn.Close()
			delete(clients, chater)
			fmt.Println(clients)
			flag <- false
		}

		fmt.Printf("Got message: %#v\n", m)
		// board cast msg
		for k := range clients {
			if k.id != chater.id {
				m.From = chater.id
				m.MyId = nil

				err := k.wsconn.WriteJSON(m)
				if err != nil {
					fmt.Printf("send failed!")
				}
			}
		}

	}
}
