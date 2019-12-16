# BarrageWall
Barrage Wall
## Login 
### URL 
/login  
### Request Body 
```json 
{
	"code":"" , 
	"rawData":"",
	"signature":"",
	"neuId":"",
	"neuPassword":""
}
```  
key | description 
-|- 
code | get from qq.login() 
rawData | get from qq.getUserInfo() 
signature | get from qq.getUserInfo() 
neuId | student id 
neuPassword | password 

### Response Data 
```json
{"code":0, "msg":"OK", "session":""}
``` 
code | description 
-|-
0 | OK 
-3 | NEU authentication fail 
-4 | QQ valid signature fail 

## Send
### URL 
/test-ws 
### Request Method 
GET
### Request Params
s=session
