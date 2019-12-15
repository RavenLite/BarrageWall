package main

import (
	_ "encoding/json"
	"fmt"
	"github.com/garyburd/redigo/redis"
	_ "github.com/garyburd/redigo/redis"
	"log"
	"net/http"
)


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
	http.HandleFunc("/test-ws", wsHandler)
	http.HandleFunc("/login", QQLogin)
	http.HandleFunc("/ws", wsHandler)
	http.HandleFunc("/user", listUsers)
	// monitor localhost:7777
	err := http.ListenAndServe("0.0.0.0:7778", nil)
	if err != nil {
		log.Fatal("ListenAndServe Monitoring Error", err.Error())
	}
}