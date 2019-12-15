package db

import (
	"encoding/json"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"github.com/satori/go.uuid"
)

type User struct {
	Name string
	Image string
	StudentId string
}

func GetUser(userUuid string) (User, error)  {
	var user User
	// build connection
	conn,err := redis.Dial("tcp","jupyter.neuyan.com:6379", redis.DialDatabase(1), redis.DialPassword("weneudb2019"))
	if err != nil {
		fmt.Println("connect redis error :",err)
		return user, err
	}
	defer conn.Close()

	// set value, add a new user json to the end of the list
	userJson, err := redis.String(conn.Do("GET", userUuid))
	if err != nil {
		fmt.Println("redis set error:", err)
		return user, err
	}
	err = json.Unmarshal([]byte(userJson), &user)
	if err != nil{
		fmt.Println("parse JSON error ")
		return user, err
	}
	return user, nil
}


func AddUser(user User) string  {
	uid,_ := uuid.NewV4()
	uuidStr := fmt.Sprint(uid)
	userJson, err := json.Marshal(user)
	if err != nil {
		fmt.Println("Transferring Error:", err)
	}

	// build connection
	conn,err := redis.Dial("tcp","jupyter.neuyan.com:6379", redis.DialDatabase(1), redis.DialPassword("weneudb2019"))
	if err != nil {
		fmt.Println("connect redis error :",err)
		return ""
	}
	defer conn.Close()

	// set value, add a new user json to the end of the list
	_, err = conn.Do("SET", uuidStr, userJson)
	_,err = conn.Do("expire",uuidStr, "1296000")
	if err != nil {
		fmt.Println("redis set error:", err)
	}
	return uuidStr
}

