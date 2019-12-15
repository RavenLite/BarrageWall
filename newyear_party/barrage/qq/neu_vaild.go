package qq

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"regexp"
	"strconv"
	"strings"
)

// 一网通办正则表达式
var ltReg = regexp.MustCompile("input type=\"hidden\" id=\"lt\" name=\"lt\" value=\"(.*?)\" />")
var executionReg = regexp.MustCompile("input type=\"hidden\" name=\"execution\" value=\"(.*?)\" />")

func ValidStudent(userId, password string) bool  {

	var err error

	// cookie Jar和client 用于维护对话
	jar,_ := cookiejar.New(nil)
	client := &http.Client{
		Jar:jar,
	}
	req, err := http.NewRequest("GET","https://pass.neu.edu.cn/tpass/login",nil)
	if err != nil {
		fmt.Printf("Create Get request https://pass.neu.edu.cn/tpass/login for %s error", userId)
		return false
	}
	resp,err := client.Do(req)
	defer _closeResp(resp)

	if err != nil {
		fmt.Printf("Request https://pass.neu.edu.cn/tpass/login for %s error", userId)
		return false
	}

	body,err :=ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Read request body for %s error", userId)
		return false
	}
	lt := ltReg.FindAllStringSubmatch(string(body), -1)[0][1]
	execution := executionReg.FindAllStringSubmatch(string(body), -1)[0][1]

	// 构造提交表单
	postStr := "rsa="+userId+password+lt+"&ul="+strconv.Itoa(len(userId))+"&pl="+strconv.Itoa(len(password))+"&lt="+lt+"&execution="+execution+"&_eventId=submit"
	postData := strings.NewReader(postStr)
	loginRequest,err := http.NewRequest("POST","https://pass.neu.edu.cn/tpass/login",postData)
	if err != nil {
		fmt.Printf("Create POST request https://pass.neu.edu.cn/tpass/login for %s error", userId)
		return false
	}
	// 这个header用于设定表单提交方式
	loginRequest.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	// User-Agent伪装头 其实没什么用
	loginRequest.Header.Add("User-Agent","Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/74.0.3729.169 Safari/537.36")
	loginResponse,err  := client.Do(loginRequest)
	defer _closeResp(loginResponse)
	if err != nil {
		fmt.Printf("Read request body for %s error", userId)
		return false
	}
	loginPage, _ := ioutil.ReadAll(loginResponse.Body)
	if strings.Contains(string(loginPage), "统一身份认证") {
		return false
	}
	return true
}

// 关闭response 防止内存溢出
func _closeResp(resp *http.Response)  {
	err := resp.Body.Close()
	if err != nil {
		fmt.Println("close Resp error")
	}
}
