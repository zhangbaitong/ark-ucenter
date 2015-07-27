package action

import (
	"common"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/RangelReale/osin"
	"github.com/dchest/authcookie"
	"github.com/julienschmidt/httprouter"
	"github.com/unrolled/render"
	"net/http"
	"strings"
	"time"
)

type (
	OAuth struct {
		Server *osin.Server
		View   *render.Render
	}

	User struct {
		Acname   string
		Password string
	}

	Res struct {
		Resname  string
		Rescname string
	}
)

const (
	//cookie加密、解密使用
	KEY        string = "QAZWERT4556"
	COOKIENAME string = "MNBVCXZ"
)

func NewOAuth() *OAuth {

	sconfig := osin.NewServerConfig()
	sconfig.AllowedAuthorizeTypes = osin.AllowedAuthorizeType{osin.CODE, osin.TOKEN}
	sconfig.AllowedAccessTypes = osin.AllowedAccessType{osin.AUTHORIZATION_CODE,
		osin.REFRESH_TOKEN, osin.PASSWORD, osin.CLIENT_CREDENTIALS, osin.ASSERTION}
	sconfig.AllowGetAccessRequest = true
	sconfig.AllowClientSecretInParams = true

	oauth := OAuth{
		Server: osin.NewServer(sconfig, NewATStorage()),
		View:   render.New(),
	}
	return &oauth
}

//申请获取授权码
func (oauth *OAuth) GetAuthorize(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Println("GetAuthorize:\r\n")
	if !checkLogin(oauth, w, r) {
		return
	}
	if !checkAuthorize(oauth, w, r, "") {
		return
	}
	doAuthorizeRequest(oauth, w, r)
}

func (oauth *OAuth) PostLogin(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	acname := login(oauth, w, r)
	if acname == "" {
		return
	}
	if !checkAuthorize(oauth, w, r, acname) {
		return
	}
	doAuthorizeRequest(oauth, w, r)
}

//登录
func login(oauth *OAuth, w http.ResponseWriter, r *http.Request) string {
	fmt.Println("Login\r\n")
	acname := oauth.Logged(w, r)
	if acname == "" {
		//使用提交的表单登陆
		acname, _ = oauth.Login(w, r)
		//登陆失败
		if acname == "" {
			//返回页面，出现 登陆失败提示，用户名密码框+授权并登陆按钮+权限列表
			common.ForwardPage(w, "./static/public/oauth2/login.html", map[string]string{"RequestURI": r.RequestURI})
		}
	}
	return acname
}

//检查是否登录，未登录，则返回登录页
func checkLogin(oauth *OAuth, w http.ResponseWriter, r *http.Request) bool {
	fmt.Println("checkLogin\r\n")
	acname := oauth.Logged(w, r)
	fmt.Println("checkLogin acname", acname)
	if acname == "" {
		common.ForwardPage(w, "./static/public/oauth2/login.html", map[string]string{"RequestURI": "/oauth2/login?" + r.URL.RawQuery})
		return false
	}
	return true
}

//检查申请资源是否被授权，如有有未授权的资源，则返回授权页
func checkAuthorize(oauth *OAuth, w http.ResponseWriter, r *http.Request, acname string) bool {
	sliceRes := []Res{}
	strRes := ""
	queryForm := common.GetUrlParam(r)
	arrScope := strings.Split(queryForm["scope"][0], ",")
	clientId := queryForm["client_id"][0]
	if acname == "" {
		acname = oauth.Logged(w, r)
	}
	openId := GetOpenIdByacName(acname, clientId)

	for i := 0; i < len(arrScope); i++ {
		resId := GetResId(arrScope[i])
		if resId > 0 {
			if !IsPersonConfered(clientId, openId, resId) {
				resCname := GetResCname(arrScope[i])
				res := Res{Resname: arrScope[i], Rescname: resCname}
				sliceRes = append(sliceRes, res)
			} else {
				if strRes == "" {
					strRes += arrScope[i]
				} else {
					strRes += "," + arrScope[i]
				}
			}
		}
	}

	if len(sliceRes) > 0 {
		requestURI := "/oauth2/authorize?response_type=" + queryForm["response_type"][0] + "&client_id=" + queryForm["client_id"][0] + "&redirect_uri=" + queryForm["redirect_uri"][0] + "&state=" + queryForm["state"][0]
		common.ForwardPage(w, "./static/public/oauth2/oauth.html", map[string]interface{}{"RequestURI": requestURI, "sliceRes": sliceRes, "strRes": strRes})
		return false
	}
	return true
}

//绑定授权码且返回授权码
func doAuthorizeRequest(oauth *OAuth, w http.ResponseWriter, r *http.Request) {
	//用户登陆成功，并确认授权，则进行下一步,根据请求,发放code 或token
	resp := oauth.Server.NewResponse()
	defer resp.Close()
	ar := oauth.Server.HandleAuthorizeRequest(resp, r)
	if ar != nil {
		//发放code 或token ,附加到redirect_uri后，并跳转
		//存储acname，acid,rsid,clientid,clientSecret等必要信息
		acname := oauth.Logged(w, r)
		if acname==""{
			acname = r.FormValue("acname")
		}
		acid := GetAcId(acname)
		ar.UserData = ATUserData{Ac_name: acname, Ac_id: acid}
		ar.Authorized = true

		oauth.Server.FinishAuthorizeRequest(resp, r, ar)
	}
	osin.OutputJSON(resp, w, r)
}

func (oauth *OAuth) PostAuthorize(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	doAuthorizeRequest(oauth, w, r)
}

func (oauth *OAuth) Token(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Println("Token:\r\n")
	resp := oauth.Server.NewResponse()
	defer resp.Close()

	if ar := oauth.Server.HandleAccessRequest(resp, r); ar != nil {
		checkAccessRequest(oauth, w, r, ar)
		oauth.Server.FinishAccessRequest(resp, r, ar)
	}
	osin.OutputJSON(resp, w, r)
}

//检查应用是否有权限访问其申请资源，以及资源是否已启用
func checkAccessRequest(oauth *OAuth, w http.ResponseWriter, r *http.Request, ar *osin.AccessRequest) *osin.AccessRequest {
	switch ar.Type {
	case osin.AUTHORIZATION_CODE:
		ar.Authorized = true

		//校验申请的资源是否已经给第三方应用授权
		resources := ""
		arrScope := strings.Split(ar.Scope, ",")
		for i := 0; i < len(arrScope); i++ {
			resId := GetResId(arrScope[i])
			if IsAppConfered(ar.Client.GetId(), resId) {
				if resources == "" {
					resources += arrScope[i]
				} else {
					resources += "," + arrScope[i]
				}

				//写入用户授权表
				userData := ar.UserData.(map[string]interface{})
				acId := int(userData["Ac_id"].(float64))
				openId := GetOpenId(acId, ar.Client.GetId())
				if !IsPersonConfered(ar.Client.GetId(), openId, resId) {
					InsertPersonConfered(ar.Client.GetId(), openId, resId)
				}
			}
		}

		//重新给token绑定资源
		ar.Scope = resources
	case osin.REFRESH_TOKEN:
		ar.Authorized = true
	case osin.PASSWORD:
		user := User{Acname: ar.Username, Password: ar.Password}
		ok := oauth.LoginQuery(&user)
		if ok {
			oauth.GenerateCookie(w, r, user.Acname, 1)
			ar.Authorized = true
		} else {
			//通过redirect_uri 返回错误约定 并跳转到改redirect_uri
		}
	case osin.CLIENT_CREDENTIALS:
		ar.Authorized = true
	case osin.ASSERTION:
		if ar.AssertionType == "urn:osin.example.complete" && ar.Assertion == "osin.data" {
			ar.Authorized = true
		}
	}

	return ar
}

func (oauth *OAuth) Logged(w http.ResponseWriter, req *http.Request) string {
	cookie, err := req.Cookie(COOKIENAME)
	if err == nil {
		return authcookie.Login(cookie.Value, []byte(KEY))
	}
	return ""
}

func (oauth *OAuth) Login(w http.ResponseWriter, req *http.Request) (string, error) {
	acname := req.FormValue("acname")
	password := req.FormValue("password")
	if acname == "" || password == "" {
		return "", errors.New("未输入用户名和密码！")
	}
	user := User{Acname: acname, Password: password}
	ok := oauth.LoginQuery(&user)
	if ok {
		oauth.GenerateCookie(w, req, user.Acname, 1)
		return acname, nil
	} else {
		return "", errors.New("用户名或密码错误！")
	}
}

//登录插入
func (oauth *OAuth) LoginQuery(user *User) bool {
	strSQL := fmt.Sprintf("select count(ac_name) from account_tab where ac_name='%s' and ac_password='%s'", user.Acname, user.Password)
	rows, err := common.GetDB().Query(strSQL)
	if err != nil {
		return false
	} else {
		defer rows.Close()
		var nCount int
		for rows.Next() {
			rows.Scan(&nCount)
		}
		if nCount == 0 {
			return false
		}
	}
	return true
}

type Logout struct {
}

func NewLogout() *Logout {
	return new(Logout)
}

func (l *Logout) Get(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	cookie := http.Cookie{Name: COOKIENAME, Path: "/", MaxAge: -1}
	http.SetCookie(w, &cookie)
}

//生成cookie，放到reponse对象
func (oauth *OAuth) GenerateCookie(w http.ResponseWriter, r *http.Request, userNmae string, number int) {
	timeLength := 24 * time.Hour
	cookieValue := authcookie.NewSinceNow(userNmae, timeLength, []byte(KEY))
	expire := time.Now().Add(timeLength)
	cookie := http.Cookie{Name: COOKIENAME, Value: cookieValue, Path: "/", Expires: expire, MaxAge: 86400}
	http.SetCookie(w, &cookie)
}

func (oauth *OAuth) Get(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	//用户登陆成功，并确认授权，则进行下一步,根据请求,发放code 或token
	strToken := req.Header.Get("Access_token")
	accessData, _ := oauth.Server.Storage.LoadAccess(strToken)
	fmt.Println(accessData)
	UserData := accessData.UserData.(map[string]interface{})
	user_name := UserData["Acname"].(string)
	acid := GetAcId(user_name)
	if accessData.Client == nil {
		fmt.Println("Get Client Faild!!!")
	}
	client_id := accessData.Client.GetId()

	//fmt.Printf("acid=%d;client_id=%s\r\n",acid,client_id)
	jr := make(map[string]interface{})
	jr["client_id"] = client_id
	if acid != -1 && client_id != "" {
		openId := GetOpenId(acid, client_id)
		jr["openid"] = openId
	}

	result, err := json.Marshal(jr)
	if err != nil {
		result = []byte("")
	}
	w.Write(result)
}

//检查token拥有的资源
func (oauth *OAuth) CheckPrivilige(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	flag := true
	queryUrl := common.GetUrlParam(r)
	clientId := queryUrl["client_id"][0]

	token := queryUrl["token"][0]
	storage, err := oauth.Server.Storage.LoadAccess(token)
	if err != nil {
		fmt.Println("get token storage failure")
	}

	if clientId != storage.Client.GetId() {
		flag = false
	} else {
		openId := ""
		if queryUrl["open_id"] != nil {
			openId = queryUrl["open_id"][0]
		}

		if openId != "" {
			userData := storage.UserData.(map[string]interface{})
			acId := int(userData["Ac_id"].(float64))
			storageOpenId := GetOpenId(acId, clientId)
			if openId != storageOpenId {
				flag = false
			}
		}
	}

	ret := make(map[string]interface{})
	if flag {
		ret["code"] = 0
		ret["data"] = storage.Scope
	} else {
		ret["code"] = 1
	}
	common.Write(w, ret)
}

//通过openId获取用户资源权限列表
func (oauth *OAuth) QueryPersonResList(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	queryUrl := common.GetUrlParam(r)
	openId := queryUrl["open_id"][0]
	personRes := GetPersonResList(openId)
	ret := make(map[string]interface{})
	ret["personRes"] = personRes
	common.Write(w, ret)
}

func (oauth *OAuth) SetUserInfo(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	strBody := []byte("{\"Code\":0,\"Message\":\"ok\"}")	
	acname := oauth.Logged(w, req)
	if acname=="" {
		strBody = []byte("{\"Code\":1,\"Message\":\"user need login\"}")
		w.Write(strBody)
		return
	}

	//acname:="18585816540"
	acid := GetAcId(acname)
	if acid==-1 {
		strBody = []byte("{\"Code\":1,\"Message\":\"user not exist\"}")
		w.Write(strBody)
		return
	}
	
	//var  info map[string]string
	req.ParseForm()
	info:=make(map[string] string)
	for k, v := range req.Form {
		info[k]= v[0]
	}

	var UserInfo ATUserInfo
	UserInfo.Ac_id=acid
	UserInfo.Info=info
	ok:=UpdateUserInfo(&UserInfo)
	if !ok {
		strBody = []byte("{\"Code\":1,\"Message\":\"save data faild\"}")
	}
	w.Write(strBody)
}

func (oauth *OAuth) MultiLogin(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	strBody := []byte("{\"Code\":0,\"Message\":\"ok\"}")	
	strFields := req.FormValue("fields")
	strValue := req.FormValue("Value")
	strPassword := req.FormValue("password")
	if strValue == "" || strPassword == "" {
		strBody = []byte("{\"Code\":1,\"Message\":\"password is empty!!\"}")
		w.Write(strBody)
		return 
	}

	strName,_ := LoginMulti(strFields,strValue,strPassword)
	if len(strName)>0 {
		oauth.GenerateCookie(w, req, strName, 1)
		w.Write(strBody)
		return 
	} else {
		return 
	}

	if !checkAuthorize(oauth, w, req, strName) {
		return
	}
	doAuthorizeRequest(oauth, w, req)
}
