package main

import (
	"BarrageWall/newyear_party/barrage/db"
	"BarrageWall/newyear_party/barrage/qq"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
)

type QQLoginRequest struct {
	code string
	rawData string
	signature string
	neuId string
	neuPassword string
}

type QQUserInfo struct {
	NickName string `json:"nickName"`
	AvatarUrl string `json:"avatarUrl"`
}

func QQLogin(w http.ResponseWriter, r *http.Request)  {
	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(403)
		return
	}
	var requestData QQLoginRequest
	_ = json.Unmarshal(requestBody, &requestData)
	qqSession, _ := qq.Code2Session(requestData.code)
	sessionKey := qqSession.SessionKey
	if qq.UserInfoValid(sessionKey, requestData.rawData, requestData.signature) {
		if qq.ValidStudent(requestData.neuId, requestData.neuPassword) {
			var userInfo QQUserInfo
			_ = json.Unmarshal([]byte(requestData.rawData), &userInfo)
			user := db.User{Name:userInfo.NickName, Image:userInfo.AvatarUrl, StudentId: requestData.neuId}
			uuidStr := db.AddUser(user)
			_, _ = io.WriteString(w, `{"code":0, "msg":"OK", "session":"`+uuidStr+`"}`)
		} else {
			_, _ = io.WriteString(w, `{"code":-3, "msg":"error neu id or password", "session":""}`)
			return
		}
	} else {
		_, _ = io.WriteString(w, `{"code":-4, "msg":"error qq data", "session":""}`)
		return
	}
}
