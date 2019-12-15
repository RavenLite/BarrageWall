package qq

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type QQSessionData struct {
	OpenID string `json:"openid"`
	SessionKey string `json:"session_key"`
	Code int  `json:"errcode"`
	Msg string `json:"errmsg"`
}

func Code2Session(code string) (QQSessionData, error)  {
	var res QQSessionData
	var requestURL string  =
		"https://api.q.qq.com/sns/jscode2session?appid=1109940479&secret=0m3OpYYy4NgTnbvm&js_code="+code+"&grant_type=authorization_code"
	response, err := http.Get(requestURL)
	if err != nil {
		return res, err
	}
	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return res, err
	}
	err = json.Unmarshal(responseData, &res)
	if err != nil {
		fmt.Println("Parse data from QQ error")
		return res,err
	}
	return res, nil
}
