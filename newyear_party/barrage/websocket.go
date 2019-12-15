package main

import (
	"BarrageWall/newyear_party/barrage/db"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"io"
	"net/http"
)

// http upgrader
var (
	upgrader = websocket.Upgrader {
		// size of read buffer
		ReadBufferSize:1024,
		// size of write buffer
		WriteBufferSize:1024,
		// allow CrossOrigin
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	clients = make(map [Client] bool)
	join = make(chan Client, 10)
	leave = make(chan Client, 10)
	message = make(chan Message, 10)
)


// message
type Message struct {
	Name string
	Image string
	Message string
	Count int
}

// websocket client
type Client struct {
	conn *websocket.Conn
	name string
}

// handler for 'ws'
func wsHandler(w http.ResponseWriter, r *http.Request) {
	var (
		wbsCon *websocket.Conn
		err error
	)

	// finish http response, add the arg in http header
	if wbsCon, err = upgrader.Upgrade(w, r, nil); err != nil {
		return
	}

	// get params and package to a struct
	vars := r.URL.Query();
	fmt.Println(vars)
	userUuid := vars["s"][0]

	// persistence
	user,err  := db.GetUser(userUuid)
	if err != nil {
		w.WriteHeader(400)
		_, _ = io.WriteString(w, "NOT CORRECT SESSION")
		return
	}

	// keep connection
	var client Client
	client.name = user.StudentId
	client.conn = wbsCon

	if !clients[client] {
		join <- client
	}

	defer func() {
		leave <- client
		_ = client.conn.Close()
	}()

	for {
		// read new message
		_, msgStr, err := client.conn.ReadMessage()

		if err != nil {
			fmt.Print("Keeping Connection Error")
			break
		}

		// pour new message into Message Panel
		var msg Message
		msg.Name = user.Name
		msg.Image = user.Image
		msg.Message = string(msgStr)
		msg.Count = len(clients)
		message <- msg
	}
}

// a global broad caster
func broadcaster() {
	fmt.Println("Broad Caster Starts")
	for {
		// wait for panel's available status
		select {
		// new message
		case msg := <-message:
			fmt.Printf("broadcaster-----------%s send message: %s\n", msg.Name, msg.Message)
			// Broadcast
			for client := range clients {
				data, err := json.Marshal(msg)
				if err != nil {
					return
				}
				if client.conn.WriteMessage(websocket.TextMessage, data) != nil {
				}
			}

		// new user join
		case client := <-join:
			fmt.Printf("broadcaster-----------%s join in the chat room\n", client.name)
			clients[client] = true

			var msg Message
			msg.Name = client.name
			msg.Message = fmt.Sprintf("%s join in, there are %d preson in room", client.name, len(clients))

			// comment: no need for our scene
			// message <- msg

		// user exit
		case client := <-leave:
			fmt.Printf("broadcaster-----------%s leave the chat room\n", client.name)
			if !clients[client] {
				break
			}

			delete(clients, client)

			var msg Message
			msg.Name = client.name
			msg.Message = fmt.Sprintf("%s leave, there are %d preson in room", client.name, len(clients))

			// comment: no need for our scene
			// message <- msg
		}
	}
}