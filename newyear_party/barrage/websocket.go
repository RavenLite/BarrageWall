package main

import (
	"encoding/json"
	_ "encoding/json"
	"fmt"
	"github.com/garyburd/redigo/redis"
	_ "github.com/garyburd/redigo/redis"
	"github.com/gorilla/websocket"
	"log"
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

// user
type User struct {
	Name string
	Image string
	StudentId string
	College string
	Gender string
}

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
	name, image, studentId, college, gender := vars["name"][0], vars["image"][0], vars["studentId"][0], vars["college"][0], vars["gender"][0]
	var user = User{name, image, studentId, college, gender}

	// persistence
	redisStorage(user)

	// keep connection
	var client Client
	client.name = studentId
	client.conn = wbsCon

	if !clients[client] {
		join <- client
	}

	defer func() {
		leave <- client
		client.conn.Close()
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
		msg.Name = name
		msg.Image = image
		msg.Message = string(msgStr)
		msg.Count = len(clients)
		message <- msg
	}
}

// persistent messages using redis
func redisStorage(user User) {
	// struct to json
	userJson, err := json.Marshal(user)
	if err != nil {
		fmt.Println("Transferring Error:", err)
	}

	// build connection
	conn,err := redis.Dial("tcp","jupyter.neuyan.com:6379", redis.DialDatabase(1), redis.DialPassword("weneudb2019"))
	if err != nil {
		fmt.Println("connect redis error :",err)
		return
	}
	defer conn.Close()

	// set value, add a new user json to the end of the list
	_, err = conn.Do("SET", user.StudentId, userJson)
	if err != nil {
		fmt.Println("redis set error:", err)
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
			fmt.Println("broadcaster-----------%s send message: %s\n", msg.Name, msg.Message)
			// Broadcase
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
			fmt.Println("broadcaster-----------%s join in the chat room\n", client.name)
			clients[client] = true

			var msg Message
			msg.Name = client.name
			msg.Message = fmt.Sprintf("%s join in, there are %d preson in room", client.name, len(clients))

			// comment: no need for our scene
			// message <- msg

		// user exit
		case client := <-leave:
			fmt.Println("broadcaster-----------%s leave the chat room\n", client.name)
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

// TODO: list all users from redis
func listUsers(w http.ResponseWriter, r *http.Request) {
	// build connection
	conn,err := redis.Dial("tcp","jupyter.neuyan.com:6379", redis.DialDatabase(1), redis.DialPassword("weneudb2019"))
	if err != nil {
		fmt.Println("connect redis error :",err)
		return
	}
	defer conn.Close()

	// get all keys
	keys, err := conn.Do("KEYS", "*")
	fmt.Println(keys)
	if err != nil {
		fmt.Println("redis KEYS error:", err)
	}

}

func main()  {
	go broadcaster()
	// upgrade http to websocket
	http.HandleFunc("/ws", wsHandler)
	http.HandleFunc("/user", listUsers)
	// monitor localhost:7777
	err := http.ListenAndServe("0.0.0.0:7777", nil)
	if err != nil {
		log.Fatal("ListenAndServe Monitoring Error", err.Error())
	}
}