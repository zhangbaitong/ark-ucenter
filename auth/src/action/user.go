package action

import (
	"common"
	"errors"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/dchest/authcookie"
	"gopkg.in/mgo.v2/bson" 
	"time"
	"net/http"
	"encoding/json"
	"strings"
)
const (
	//cookie加密、解密使用
	KEY        string = "QAZWERT4556"
	COOKIENAME string = "MNBVCXZ"
)

type Response struct {
	Code int
	Message string
}

//生成cookie，放到reponse对象
func GenerateCookie(w http.ResponseWriter, r *http.Request, userNmae string, number int) {
	timeLength := 24 * time.Hour
	cookieValue := authcookie.NewSinceNow(userNmae, timeLength, []byte(KEY))
	expire := time.Now().Add(timeLength)
	cookie := http.Cookie{Name: COOKIENAME, Value: cookieValue, Path: "/", Expires: expire, MaxAge: 86400}
	http.SetCookie(w, &cookie)
}

func GetCookieName(req *http.Request) string {
	cookie, err := req.Cookie(COOKIENAME)
	if err == nil {
		return authcookie.Login(cookie.Value, []byte(KEY))
	}
	return ""
}

func Register(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	req.ParseForm()
	reg_type := req.FormValue("reg_type")
	if reg_type=="1" {
		strBody := []byte("{\"Code\":0,\"Message\":\"ok\"}")	
		Info:=make(map[string]string)
		for k, v := range req.Form {
			if k== "reg_type" {
				continue
			}
			Info[k]=v[0]
		}

		UserData,code:=MultiRegister(&Info)
		if code==0 {
			strUser, err := json.Marshal(UserData)
			if err != nil {
				strBody = []byte(fmt.Sprintf("{\"Code\":%d,\"Message\":\"%s\"}",UNKNOWN_ERROR,GetError(UNKNOWN_ERROR)))
			} else {
				strTemp:="{\"Code\":0,\"Message\":\""+string(strUser)+"\"}"
				strBody = []byte(strTemp)	
			}
		} else{
			strBody = []byte(fmt.Sprintf("{\"Code\":%d,\"Message\":\"%s\"}",code,GetError(code)))
		}
		w.Write(strBody)
		return 
	}

	fmt.Println(req.Form)

	acname := req.FormValue("acname")
	password := req.FormValue("password")
	if acname == ""  || password==""{
		strBody := []byte(fmt.Sprintf("{\"Code\":%d,\"Message\":\"%s\"}",PARAM_ERROR,GetError(PARAM_ERROR)))
		w.Write(strBody)
		return 
	}
	strBody := []byte("{\"Code\":0,\"Message\":\"ok\"}")	
	_,ok:=GetUser(acname)
	if ok {
		strBody = []byte(fmt.Sprintf("{\"Code\":%d,\"Message\":\"%s\"}",USER_EX,GetError(USER_EX)))
		w.Write(strBody)
		return 
	}

	account := Account{Ac_name: acname, Ac_password: password}
	account.Id=bson.NewObjectId()
	ok = RegisterInsert(&account)
	if !ok {
		strBody = []byte(fmt.Sprintf("{\"Code\":%d,\"Message\":\"%s\"}",INSERT_DB_ERROR,GetError(INSERT_DB_ERROR)))
		w.Write(strBody)
		return
	}

	//save others info
	Info:=make(map[string]string)
	for k, v := range req.Form {
		if k== "acname" || k== "password"  || k== "reg_type" {
			continue
		}

		Info[k]=v[0]
	}

	UserData,_:=GetUser(acname)
	UserInfo:=ATUserInfo{}
	UserInfo.Id=account.Id
	UserInfo.Ac_id=UserData.Ac_id
	UserInfo.Info=Info
	InsertUserInfo(&UserInfo)

	InfoAll:=UserInfoAll{}
	InfoAll.Id              =account.Id.Hex()
	InfoAll.Ac_name 		=UserData.Ac_name
	InfoAll.Status		=UserData.Status
	InfoAll.Source   		=UserData.Source
	InfoAll.Create_time  =UserData.Create_time
	InfoAll.Info           =Info
	strUser, err := json.Marshal(InfoAll)
	if err != nil {
		strBody = []byte(fmt.Sprintf("{\"Code\":%d,\"Message\":\"%s\"}",UNKNOWN_ERROR,GetError(UNKNOWN_ERROR)))
	} else {
		 var response Response		
		 response.Code=0
		 response.Message=string(strUser)
		 strBody, _ = json.Marshal(response)
	}

	w.Write(strBody)
}

//登录
func UserLogin(w http.ResponseWriter, r *http.Request) string {
	fmt.Println("Login\r\n")
	acname := GetCookieName(r)
	if acname == "" {
		//使用提交的表单登陆
		acname, _ = Login(w, r)
		//登陆失败
		if acname == "" {
			//返回页面，出现 登陆失败提示，用户名密码框+授权并登陆按钮+权限列表
			common.ForwardPage(w, "./static/public/oauth2/login.html", map[string]string{"RequestURI": r.RequestURI})
		}
	}
	return acname
}

func LoginCenter(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	acname := req.FormValue("acname")
	password := req.FormValue("password")
	strBody:=[]byte("")
	if acname == "" || password == "" {
		strBody = []byte("{\"Code\":1,\"Message\":\"user name or password is not empty\"}")
		w.Write(strBody)
	}
	user := User{Acname: acname, Password: password}
	ok := LoginQuery(&user)
	if ok {
		UserData,_:=GetUser(acname)
		UserInfo,_:=GetUserInfoM(UserData.Mid)
		//GenerateCookie(w, req, user.Acname, 1)
		InfoAll:=UserInfoAll{}
		InfoAll.Id              =UserData.Mid
		InfoAll.Ac_name 		=UserData.Ac_name
		InfoAll.Status		=UserData.Status
		InfoAll.Source   		=UserData.Source
		InfoAll.Create_time  =UserData.Create_time
		InfoAll.Info           =UserInfo.Info
		strUser, err := json.Marshal(InfoAll)
		if err != nil {
			strBody = []byte("{\"Code\":1,\"Message\":\"json encode faild\"}")
		} else {
			 var response Response		
			 response.Code=0
			 response.Message=string(strUser)
			 strBody, _ = json.Marshal(response)
		}
	} else {
		strBody = []byte(fmt.Sprintf("{\"Code\":%d,\"Message\":\"%s\"}",USER_NOT_EX,GetError(USER_NOT_EX)))
	}
	w.Write(strBody)
}

func Login(w http.ResponseWriter, req *http.Request) (string, error) {
	acname := req.FormValue("acname")
	password := req.FormValue("password")
	if acname == "" || password == "" {
		return "", errors.New("未输入用户名和密码！")
	}
	user := User{Acname: acname, Password: password}
	ok := LoginQuery(&user)
	if ok {
		GenerateCookie(w, req, user.Acname, 1)
		return acname, nil
	} else {
		return "", errors.New("用户名或密码错误！")
	}
}

func Logout(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	cookie := http.Cookie{Name: COOKIENAME, Path: "/", MaxAge: -1}
	http.SetCookie(w, &cookie)
}

//通过openId获取用户资源权限列表
func SetUserInfo(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	
	strBody := []byte("{\"Code\":0,\"Message\":\"ok\"}")
	id := req.FormValue("id")
	//acname := GetCookieName(req)
	if id == "" {
		strBody := []byte(fmt.Sprintf("{\"Code\":%d,\"Message\":\"%s\"}",PARAM_ERROR,GetError(PARAM_ERROR)))
		w.Write(strBody)
		return 
	}

	//acname:="18585816540"
	UserData,ok := GetUserById(id)
	if !ok {
		strBody := []byte(fmt.Sprintf("{\"Code\":%d,\"Message\":\"%s\"}",USER_NOT_EX,GetError(USER_NOT_EX)))
		w.Write(strBody)
		return
	}

	//var  info map[string]string
	info := make(map[string]string)
	for k, v := range req.Form {
		if k== "id" {
			continue
		}
		info[k] = v[0]
	}

	var UserInfo ATUserInfo
	UserInfo.Ac_id = UserData.Ac_id
	UserInfo.Info = info
	ok = UpdateUserInfo(&UserInfo)
	if !ok {
		strBody = []byte(fmt.Sprintf("{\"Code\":%d,\"Message\":\"%s\"}",UPDATE_DB_ERROR,GetError(UPDATE_DB_ERROR)))
	}
	w.Write(strBody)	
}

func GetUserInfo(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	strBody := []byte("{\"Code\":0,\"Message\":\"ok\"}")
	id := req.FormValue("id")
	if id == "" {
		strBody := []byte(fmt.Sprintf("{\"Code\":%d,\"Message\":\"%s\"}",PARAM_ERROR,GetError(PARAM_ERROR)))
		w.Write(strBody)
		return 
	}

	UserData,ok := GetUserById(id)
	if !ok {
		strBody := []byte(fmt.Sprintf("{\"Code\":%d,\"Message\":\"%s\"}",USER_NOT_EX,GetError(USER_NOT_EX)))
		w.Write(strBody)
		return 
	}

	UserInfo,ok:=GetUserInfoM(UserData.Mid)
	InfoAll:=UserInfoAll{}
	InfoAll.Id              =UserData.Mid
	InfoAll.Ac_name 		=UserData.Ac_name
	InfoAll.Status		=UserData.Status
	InfoAll.Source   		=UserData.Source
	InfoAll.Create_time  =UserData.Create_time
	InfoAll.Info           =UserInfo.Info
	strUser, err := json.Marshal(InfoAll)
	if err != nil {
		strBody = []byte("{\"Code\":1,\"Message\":\"json encode faild\"}")
	} else {
		 var response Response		
		 response.Code=0
		 response.Message=string(strUser)
		 strBody, _ = json.Marshal(response)
	}

	w.Write(strBody)
}

func GetUserList(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	strBody := []byte("{\"Code\":0,\"Message\":\"ok\"}")
	user_list := req.FormValue("user_list")
	fmt.Println(user_list)
	if user_list == "" {
		strBody = []byte(fmt.Sprintf("{\"Code\":%d,\"Message\":\"%s\"}",LIST_EMPTY,GetError(LIST_EMPTY)))
		w.Write(strBody)
		return
	}

	UserList:=strings.Split(user_list, ",")
	List := make(map[string]UserInfoAll)
	for i := 0; i < len(UserList); i++ {
		UserData,ok := GetUserById(UserList[i])
		if !ok {
			continue
		}

		UserInfo,ok:=GetUserInfoM(UserData.Mid)
		InfoAll:=UserInfoAll{}
		InfoAll.Id              =UserData.Mid
		InfoAll.Ac_name 		=UserData.Ac_name
		InfoAll.Status		=UserData.Status
		InfoAll.Source   		=UserData.Source
		InfoAll.Create_time  =UserData.Create_time
		InfoAll.Info           =UserInfo.Info
		List[UserList[i]]=InfoAll
	}

	strUser, err := json.Marshal(&List)
	if err != nil {
		strBody = []byte("{\"Code\":1,\"Message\":\"json encode faild\"}")
	} else {
		var response Response		
		response.Code=0
		response.Message=string(strUser)
		strBody, _ = json.Marshal(response)
	}
	w.Write(strBody)
}
/*
func MultiLogin(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	strBody := []byte("{\"Code\":0,\"Message\":\"ok\"}")
	strName := req.FormValue("acname")
	strPassword := req.FormValue("password")
	if strName == "" || strPassword == "" {
		strBody = []byte("{\"Code\":1,\"Message\":\"password is empty!!\"}")
		w.Write(strBody)
		return
	}

	strACName, _ := LoginMulti(strName, strPassword)
	if len(strACName) > 0 {
		GenerateCookie(w, req, strACName, 1)
		w.Write(strBody)
		return
	} else {
		strBody = []byte("{\"Code\":1,\"Message\":\"password is empty!!\"}")
		return
	}
}
*/