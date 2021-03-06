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
	"encoding/hex"
	"encoding/base64"
	"math/rand"
	"strconv" 
	"net/url"
)
const (
	//cookie加密、解密使用
	KEY        string = "QAZWERT4556"
	COOKIENAME string = "MNBVCXZ"
)

const (
//SMS string="http://sms.gyjbh.nvwayun.com/application/api?data=%s&interface_key=%s&interface_sign=%s"
SMS string="http://sms.gyjbh.nvwayun.com/application/api?data=%s&interface_key=%s&interface_sign=%s"
)
var (
	SMSHost string
	TnterfaceKey string
	InterfaceSign string
	CheckText string
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

func Exist(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	req.ParseForm()
	Info:=make(map[string]bool)
	fmt.Println(req.Form)
	var ok bool
	for k, v := range req.Form {
		if k== "acname" {
			_,ok=GetUser(v[0])
		} else {
			_,ok=isUserExist("info."+strings.ToLower(k),v[0])
		}
		Info[k]=ok
	}

	strBody := []byte("{\"Code\":0,\"Message\":\"ok\"}")	
	strUser, err := json.Marshal(Info)
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

func Register(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	req.ParseForm()
	fmt.Println(req.Form)
	reg_type := req.FormValue("reg_type")
	show_id := req.FormValue("show_id")
	if reg_type=="1" {
		strBody := []byte("{\"Code\":0,\"Message\":\"ok\"}")	
		Info:=make(map[string]string)
		for k, v := range req.Form {
			if k== "reg_type" || k== "show_id" {
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
				 var response Response		
				 response.Code=0
				 response.Message=string(strUser)
				 strBody, _ = json.Marshal(response)
			}
		} else{
			if code==USER_EX {
				strBody= []byte(fmt.Sprintf("{\"Code\":%d,\"Message\":\"{\\\"Id\\\":\\\"%s\\\"}\"}",code,UserData.Id))
			} else {
				strBody = []byte(fmt.Sprintf("{\"Code\":%d,\"Message\":\"%s\"}",code,GetError(code)))
			}
		}
		w.Write(strBody)

		SetUserLocus(UserData.Id,UserData.Id,show_id,&Info)

		return 
	}

	acname := req.FormValue("acname")
	password := req.FormValue("password")

	source_id := req.FormValue("source_id")
	nSource:=0
	if source_id == "" {
		source_id="funzhou"
	}

	if acname == ""  || password==""{
		strBody := []byte(fmt.Sprintf("{\"Code\":%d,\"Message\":\"%s\"}",PARAM_ERROR,GetError(PARAM_ERROR)))
		w.Write(strBody)
		return 
	}

	//save others info
	Info:=make(map[string]string)
	for k, v := range req.Form {
		if k== "acname" || k== "password"  || k== "source" || k== "source_id" || k== "reg_type" || k== "show_id" {
			continue
		}

		Info[k]=v[0]
	}

	strBody := []byte("{\"Code\":0,\"Message\":\"ok\"}")	
	UserData,ok:=GetUser(acname)
	if ok {
		strBody= []byte(fmt.Sprintf("{\"Code\":%d,\"Message\":\"{\\\"Id\\\":\\\"%s\\\"}\"}",USER_EX,UserData.Mid))
		SetUserLocus(UserData.Mid,UserData.Mid,show_id,&Info)
		w.Write(strBody)
		return 
	}

	//password=common.MD5(password)
	Id:=bson.NewObjectId()
	ok = RegisterInsert(acname,password,Id.Hex(),nSource,source_id)
	if !ok {
		strBody = []byte(fmt.Sprintf("{\"Code\":%d,\"Message\":\"%s\"}",INSERT_DB_ERROR,GetError(INSERT_DB_ERROR)))
		w.Write(strBody)
		return
	}

	UserData,_=GetUser(acname)
	UserInfo:=ATUserInfo{}
	UserInfo.Id=Id
	UserInfo.Ac_id=UserData.Ac_id
	UserInfo.Create_time=UserData.Create_time
	UserInfo.Info=Info
	InsertUserInfo(&UserInfo)

	InfoAll:=UserInfoAll{}
	InfoAll.Id              =Id.Hex()
	InfoAll.Ac_name 		=UserData.Ac_name
	InfoAll.Status		=UserData.Status
	InfoAll.Source   		=UserData.Source
	InfoAll.Source_id   	=UserData.Source_id
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

		SetUserLocus(UserData.Mid,UserData.Mid,show_id,&Info)
		
	}

	w.Write(strBody)
}

func GetSMSUrl(strMobile,strMessage string) string {
	var Url *url.URL 
	//Url, err := url.Parse("http://sms.gyjbh.nvwayun.com")
	Url, err := url.Parse(SMSHost)
	if err != nil {
		return ""
	}

	mesaage:=make(map[string]string)
	mesaage["mobile"]=strMobile
	mesaage["msg"]=strMessage
	strData, err := json.Marshal(mesaage)
	if err != nil {
		return ""
	}

	Url.Path += "/application/api"
	strSend:=base64.StdEncoding.EncodeToString(strData)
	parameters := url.Values{}
	parameters.Add("data", strSend)
	parameters.Add("interface_key", TnterfaceKey)
	parameters.Add("interface_sign", InterfaceSign)
	Url.RawQuery = parameters.Encode()

	return Url.String()
}

func GetVerifyCode(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	req.ParseForm()
	fmt.Println(req.Form)
	mobile := req.FormValue("mobile")

	source_id := req.FormValue("source_id")
	show_id := req.FormValue("show_id")
	nSource:=1
	if source_id == "" {
		source_id="dzq"
	}

	if mobile == "" {
		strBody := []byte(fmt.Sprintf("{\"Code\":%d,\"Message\":\"%s\"}",PARAM_ERROR,GetError(PARAM_ERROR)))
		w.Write(strBody)
		return 
	}

	strBody := []byte("{\"Code\":0,\"Message\":\"ok\"}")	
	_,ok:=GetUser(mobile)
	if !ok {
		var Id bson.ObjectId 
		UserInfo,ok:=isUserExist("info.mobile",mobile)
		if ok {
			Id=UserInfo.Id
		} else {
			Id=bson.NewObjectId()
		}
		ok_reg := RegisterInsert(mobile,"",Id.Hex(),nSource,source_id)
		if !ok_reg {
			strBody = []byte(fmt.Sprintf("{\"Code\":%d,\"Message\":\"%s\"}",INSERT_DB_ERROR,GetError(INSERT_DB_ERROR)))
			w.Write(strBody)
			return
		}

		Info:=make(map[string]string)
		UserData,_:=GetUser(mobile)
		UserInfo_I:=ATUserInfo{}
		UserInfo_I.Id=Id
		UserInfo_I.Ac_id=UserData.Ac_id
		UserInfo_I.Create_time=UserData.Create_time
		Info["mobile"]=mobile
		UserInfo_I.Info=Info

		if !ok {
			InsertUserInfo(&UserInfo_I)

			SetUserLocus(UserData.Mid,UserData.Mid,show_id,&Info)

		} else {
			ok = UpdateUserInfo(&UserInfo_I,show_id)
			if !ok {
				strBody = []byte(fmt.Sprintf("{\"Code\":%d,\"Message\":\"%s\"}",UPDATE_DB_ERROR,GetError(UPDATE_DB_ERROR)))
			}
		}
	} else {
		UserInfo,ok:=isUserExist("info.mobile",mobile)
		if ok {
			Info:=make(map[string]string)
			Info["mobile"]=mobile
			SetUserLocus(UserInfo.Id.Hex(),UserInfo.Id.Hex(),show_id,&Info)
		}
	}

	code,ok:=Get_Verify_Code(mobile)
	if !ok {
		//Build Verify Code
	    rand.Seed( time.Now().UTC().UnixNano())
	    code=rand.Int()%1000000
	    ok=SaveVerifyCode(mobile,code)
	    if !ok {
			strBody = []byte(fmt.Sprintf("{\"Code\":%d,\"Message\":\"%s\"}",INSERT_DB_ERROR,GetError(INSERT_DB_ERROR)))
			w.Write(strBody)
			return    	
	    }
	}
	strMessage:=fmt.Sprintf(CheckText,code)

	//send SMS
	strSendURL:=GetSMSUrl(mobile,strMessage)
	if strSendURL=="" {
		strBody = []byte(fmt.Sprintf("{\"Code\":%d,\"Message\":\"%s\"}",SMS_ERROR,GetError(SMS_ERROR)))
		w.Write(strBody)
		return
	}

	strResult,err:=common.Invoker(common.HTTP_GET,strSendURL,"")
	if err!=nil {
		strBody = []byte(fmt.Sprintf("{\"Code\":%d,\"Message\":\"%s\"}",SMS_ERROR,GetError(SMS_ERROR)))
		w.Write(strBody)
		return
	}

	fmt.Println(strResult)
	result:=make(map[string]string)
	json.Unmarshal([]byte(strResult),&result)	

	if result["result"]=="0" {
		strBody = []byte(fmt.Sprintf("{\"Code\":0,\"Message\":\"{\"verify_code\":\"%06d\"}\"}",code))
	} else {
		strBody = []byte(fmt.Sprintf("{\"Code\":%d,\"Message\":\"%s  error code %s\"}",SMS_ERROR,GetError(SMS_ERROR),result["result"]))
	}

	w.Write(strBody)
}

func CheckVerifyCode(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	req.ParseForm()
	fmt.Println(req.Form)
	mobile := req.FormValue("mobile")
	verify_code := req.FormValue("verify_code")

	strBody := []byte("{\"Code\":0,\"Message\":\"ok\"}")	
	if mobile == "" || verify_code=="" {
		strBody = []byte(fmt.Sprintf("{\"Code\":%d,\"Message\":\"%s\"}",PARAM_ERROR,GetError(PARAM_ERROR)))
		w.Write(strBody)
		return 
	}

	code, err := strconv.Atoi(verify_code)
	if err != nil {
		strBody = []byte(fmt.Sprintf("{\"Code\":%d,\"Message\":\"%s\"}",PARAM_ERROR,GetError(PARAM_ERROR)))
		w.Write(strBody)
		return 
	} 

	ok:=VerifyCodeCheck(mobile,code)
	if !ok {
		strBody = []byte(fmt.Sprintf("{\"Code\":%d,\"Message\":\"%s\"}",CHECK_VERIFY_ERROR,GetError(CHECK_VERIFY_ERROR)))
		w.Write(strBody)
		return 
	}

	UserData,_:=GetUser(mobile)
	UserInfo,_:=GetUserInfoM(UserData.Mid)
	InfoAll:=UserInfoAll{}
	InfoAll.Id              =UserData.Mid
	InfoAll.Ac_name 		=UserData.Ac_name
	InfoAll.Status		=UserData.Status
	InfoAll.Source   		=UserData.Source
	InfoAll.Source_id   	=UserData.Source_id
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

//登录
func UserLogin(w http.ResponseWriter, r *http.Request) string {
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
	req.ParseForm()
	fmt.Println(req.Form)

	acname := req.FormValue("acname")
	password := req.FormValue("password")
	strBody:=[]byte("")
	if acname == "" || password == "" {
		strBody = []byte("{\"Code\":1,\"Message\":\"user name or password is not empty\"}")
		w.Write(strBody)
	}
	//password=common.MD5(password)
	ok := LoginQuery(acname,password)
	if ok {
		UserData,_:=GetUser(acname)
		UserInfo,_:=GetUserInfoM(UserData.Mid)
		//GenerateCookie(w, req, user.Acname, 1)
		InfoAll:=UserInfoAll{}
		InfoAll.Id              =UserData.Mid
		InfoAll.Ac_name 		=UserData.Ac_name
		InfoAll.Status		=UserData.Status
		InfoAll.Source   		=UserData.Source
		InfoAll.Source_id	=UserData.Source_id
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
	password=common.MD5(password)
	ok := LoginQuery(acname,password)
	if ok {
		GenerateCookie(w, req, acname, 1)
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
	req.ParseForm()
	fmt.Println(req.Form)

	strBody := []byte("{\"Code\":0,\"Message\":\"ok\"}")
	id := req.FormValue("id")
	show_id := req.FormValue("show_id")
	//acname := GetCookieName(req)
	if id == "" {
		strBody := []byte(fmt.Sprintf("{\"Code\":%d,\"Message\":\"%s\"}",PARAM_ERROR,GetError(PARAM_ERROR)))
		w.Write(strBody)
		return 
	}

	d, err := hex.DecodeString(id)
	if err != nil || len(d) != 12 {
		//Invalid input to ObjectIdHex
		strBody = []byte(fmt.Sprintf("{\"Code\":%d,\"Message\":\"%s\"}",PARAM_ERROR,GetError(PARAM_ERROR)))
		w.Write(strBody)
		return 
	}
	
	//acname:="18585816540"
	UserData,ok := GetUserById(id)
	var UserInfo ATUserInfo
	if ok {
		UserInfo.Ac_id = UserData.Ac_id
	} else {
		UserInfo.Ac_id = -1
	}

	//var  info map[string]string
	info := make(map[string]string)
	for k, v := range req.Form {
		if k== "id" ||  k== "show_id"{
			continue
		}
		info[k] = v[0]
	}

	UserInfo.Info = info
	UserInfo.Id=bson.ObjectIdHex(id)
	ok = UpdateUserInfo(&UserInfo,show_id)
	if !ok {
		strBody = []byte(fmt.Sprintf("{\"Code\":%d,\"Message\":\"%s\"}",UPDATE_DB_ERROR,GetError(UPDATE_DB_ERROR)))
	}
	w.Write(strBody)	
}

func GetUserInfo(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	req.ParseForm()
	var ok bool
	var UserData* ATUserData
	strBody := []byte("{\"Code\":0,\"Message\":\"ok\"}")
	start:=time.Now().UTC().UnixNano()
	show_id := req.FormValue("show_id")
	if len(show_id) == 0 {
		for k, v := range req.Form {
			if len(v[0])==0 {
				strBody = []byte(fmt.Sprintf("{\"Code\":%d,\"Message\":\"%s\"}",PARAM_ERROR,GetError(PARAM_ERROR)))
				w.Write(strBody)
				return 
			}
			if k== "acname"||  k=="id" {
				var strID string
				strID=v[0]
				if k== "acname" {
					UserData,ok=GetUser(v[0])
					if !ok {
						strBody = []byte(fmt.Sprintf("{\"Code\":%d,\"Message\":\"%s\"}",USER_NOT_EX,GetError(USER_NOT_EX)))
						w.Write(strBody)
						return 
					}
					strID=UserData.Mid
				} 

				if k=="id" {
					d, err := hex.DecodeString(v[0])
					if err != nil || len(d) != 12 {
						//Invalid input to ObjectIdHex
						strBody = []byte(fmt.Sprintf("{\"Code\":%d,\"Message\":\"%s\"}",PARAM_ERROR,GetError(PARAM_ERROR)))
						w.Write(strBody)
						return 
					}

					UserData,ok = GetUserById(v[0])				
				}

				UserInfo,ok2:=GetUserInfoM(strID)
				if !ok2 {
					strBody = []byte(fmt.Sprintf("{\"Code\":%d,\"Message\":\"%s\"}",USER_NOT_EX,GetError(USER_NOT_EX)))
					w.Write(strBody)
					return 
				}

				InfoAll:=UserInfoAll{}
				if ok {
					InfoAll.Ac_name 		=UserData.Ac_name
					InfoAll.Status		=UserData.Status
					InfoAll.Source   		=UserData.Source
					InfoAll.Source_id	=UserData.Source_id
					InfoAll.Create_time  =UserData.Create_time
				}
				InfoAll.Id              =strID
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
				strNode:="info."+strings.ToLower(k)
				fmt.Println(strNode,":",v[0])
				UserInfo,ok:=isUserExist(strNode,v[0])
				if !ok {
					strBody = []byte(fmt.Sprintf("{\"Code\":%d,\"Message\":\"%s\"}",USER_NOT_EX,GetError(USER_NOT_EX)))
					w.Write(strBody)
					return 
				}			

				InfoAll:=UserInfoAll{}
				InfoAll.Id=UserInfo.Id.Hex()
				if UserInfo.Ac_id>-1 {
					UserData,ok=GetUserByAcId(UserInfo.Ac_id)
					if ok {
						InfoAll.Id              =UserData.Mid
						InfoAll.Ac_name 		=UserData.Ac_name
						InfoAll.Status		=UserData.Status
						InfoAll.Source   		=UserData.Source
						InfoAll.Source_id	=UserData.Source_id
						InfoAll.Create_time  =UserData.Create_time
					} else {
					}		

				} 

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
			}
			break
		}
	} else {
		for k, v := range req.Form {
			if k== "show_id" {
				continue
			}
			if len(v[0])==0 {
				strBody = []byte(fmt.Sprintf("{\"Code\":%d,\"Message\":\"%s\"}",PARAM_ERROR,GetError(PARAM_ERROR)))
				w.Write(strBody)
				return 
			}
			if k== "acname"||  k=="id" {
				var strID string
				strID=v[0]
				if k== "acname" {
					UserData,ok=GetUser(v[0])
					if !ok {
						strBody = []byte(fmt.Sprintf("{\"Code\":%d,\"Message\":\"%s\"}",USER_NOT_EX,GetError(USER_NOT_EX)))
						w.Write(strBody)
						return 
					}
					strID=UserData.Mid
				} 

				if k=="id" {
					d, err := hex.DecodeString(v[0])
					if err != nil || len(d) != 12 {
						//Invalid input to ObjectIdHex
						strBody = []byte(fmt.Sprintf("{\"Code\":%d,\"Message\":\"%s\"}",PARAM_ERROR,GetError(PARAM_ERROR)))
						w.Write(strBody)
						return 
					}

					UserData,ok = GetUserById(v[0])
				}

				Locus,ok2:=isUserExistL("mid_c",strID,show_id)
				if !ok2 {
					strBody = []byte(fmt.Sprintf("{\"Code\":%d,\"Message\":\"%s\"}",USER_NOT_EX,GetError(USER_NOT_EX)))
					w.Write(strBody)
					return 
				}

				InfoAll:=UserInfoAll{}
				if ok {
					InfoAll.Ac_name 		=UserData.Ac_name
					InfoAll.Status		=UserData.Status
					InfoAll.Source   		=UserData.Source
					InfoAll.Source_id	=UserData.Source_id
					InfoAll.Create_time  =UserData.Create_time
				} else {
					InfoAll.Create_time  =Locus.Create_time
				}
				InfoAll.Id              =strID
				InfoAll.Info           =Locus.Info
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
				strNode:="info."+strings.ToLower(k)
				fmt.Println(strNode,":",v[0])
				Locus,ok:=isUserExistL(strNode,v[0],show_id)
				if !ok {
					strBody = []byte(fmt.Sprintf("{\"Code\":%d,\"Message\":\"%s\"}",USER_NOT_EX,GetError(USER_NOT_EX)))
					w.Write(strBody)
					return 
				}			

				InfoAll:=UserInfoAll{}

				InfoAll.Id=Locus.Mid_c
				InfoAll.Create_time  =Locus.Create_time
				InfoAll.Info           =Locus.Info

				strUser, err := json.Marshal(InfoAll)
				if err != nil {
					strBody = []byte("{\"Code\":1,\"Message\":\"json encode faild\"}")
				} else {
					 var response Response		
					 response.Code=0
					 response.Message=string(strUser)
					 strBody, _ = json.Marshal(response)
				}
			}
			break
		}		
	}
	w.Write(strBody)
	end:=time.Now().UTC().UnixNano()
	fmt.Println("used time ",(end-start)/1000)
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

	show_id := req.FormValue("show_id")
	UserList:=strings.Split(user_list, ",")
	List := make(map[string]UserInfoAll)
	if len(show_id) == 0 {		
		for i := 0; i < len(UserList); i++ {
			InfoAll:=UserInfoAll{}

			d, err := hex.DecodeString(UserList[i])
			if err != nil || len(d) != 12 {
				List[UserList[i]]=InfoAll
				continue
			}

			UserInfo,ok:=GetUserInfoM(UserList[i])
			if !ok {
				List[UserList[i]]=InfoAll
				continue
			}

			InfoAll.Id              =UserList[i]
			InfoAll.Info           =UserInfo.Info
			UserData,ok := GetUserById(UserList[i])
			if ok {
				InfoAll.Ac_name 		=UserData.Ac_name
				InfoAll.Status		=UserData.Status
				InfoAll.Source   		=UserData.Source
				InfoAll.Source_id	=UserData.Source_id
				InfoAll.Create_time  =UserData.Create_time
			} else {
				InfoAll.Create_time  =UserInfo.Create_time			
			}
			
			List[UserList[i]]=InfoAll
		}
	} else {
		for i := 0; i < len(UserList); i++ {
			InfoAll:=UserInfoAll{}

			d, err := hex.DecodeString(UserList[i])
			if err != nil || len(d) != 12 {
				List[UserList[i]]=InfoAll
				continue
			}

			Locus,ok:=isUserExistL("mid_c",UserList[i],show_id)
			if !ok {
				List[UserList[i]]=InfoAll
				continue
			}

			InfoAll.Id              =UserList[i]
			InfoAll.Info           =Locus.Info
			InfoAll.Create_time  =Locus.Create_time								
			List[UserList[i]]=InfoAll
		}		
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

func ChangePassword(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	req.ParseForm()
	fmt.Println(req.Form)

	acname := req.FormValue("acname")
	password := req.FormValue("password")
	new_password := req.FormValue("new_password")
	strBody:=[]byte("")
	if acname == "" || password == "" {
		strBody = []byte(fmt.Sprintf("{\"Code\":%d,\"Message\":\"%s\"}",PARAM_ERROR,GetError(PARAM_ERROR)))
		w.Write(strBody)
	}
	//password=common.MD5(password)
	ok := LoginQuery(acname,password)
	if ok {
		ok=UpdatePassword(acname,password,new_password)
		if ok {
			strBody = []byte("{\"Code\":0,\"Message\":\"ok\"}")
		} else {
			strBody = []byte(fmt.Sprintf("{\"Code\":%d,\"Message\":\"%s\"}",UPDATE_DB_ERROR,GetError(UPDATE_DB_ERROR)))
		}
	} else {
		strBody = []byte(fmt.Sprintf("{\"Code\":%d,\"Message\":\"%s\"}",USER_NOT_EX,GetError(USER_NOT_EX)))
	}
	w.Write(strBody)
}
