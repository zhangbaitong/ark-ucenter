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
	"net/url"
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

func (oauth *OAuth) GetAuthorize(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Println("GetAuthorize:\r\n")
	resp := oauth.Server.NewResponse()
	defer resp.Close()

	queryForm, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		resp.SetError("get http param error", "")
		osin.OutputJSON(resp, w, r)
		return
	}
	if !checkCodeRequest(resp, queryForm) {
		osin.OutputJSON(resp, w, r)
		return
	}

	acname := oauth.Logged(w, r)
	if acname != "" {
		//已经登录，则返回页面，出现 授权按钮+权限列表
		oauth.View.HTML(w, http.StatusOK, "oauth", map[string]string{"AuthorizeDisplay": "block", "LoginDisplay": "none", "RequestURI": r.RequestURI})
	} else {
		//未登录，则返回页面，出现 用户名密码框+授权并登陆按钮+权限列表
		oauth.View.HTML(w, http.StatusOK, "oauth", map[string]string{"AuthorizeDisplay": "none", "LoginDisplay": "block", "RequestURI": r.RequestURI})
	}
}

//检查应用是否有权限访问其申请资源，以及资源是否已启用
func checkCodeRequest(w *osin.Response, queryForm map[string][]string) bool {
	//校验参数是否完整
	return true
	
	if queryForm["open_id"] == nil {
		w.SetError("param open_id can not be empty", "")
		return false
	}
	
	if queryForm["scope"] == nil {
		w.SetError("param scope can not be empty", "")
		return false
	}

	//检查资源是否启用
	resId := "AT"
	status := GetResStatus(resId)
	if status != 1 {
		w.SetError("resource ["+resId+"] is not enable", "")
		return false
	}

	//通过openId获取acId
	openId := queryForm["open_id"][0]
	if openId == "" {
		w.SetError("open_id can not be empty", "")
		return false
	}
	fmt.Println("openId", openId)

	acId := GetAcIdByresIdAndopenId(openId)
	if acId <= 0 {
		w.SetError("can not find acid ", "")
		return false
	}

	//通过资源ID和acId获取权限
	priviliges := GetPriviliges(resId, acId)
	if priviliges == "" {
		w.SetError("no priviliges", "")
		return false
	}

	//校验参数中的权限是否被运行访问
	scopeArr := strings.Split(queryForm["scope"][0], ",")
	if len(scopeArr) <= 0 {
		w.SetError("no priviliges specified", "")
		return false
	}
	for i := 0; i < len(scopeArr); i++ {
		if !strings.Contains(priviliges, scopeArr[i]) {
			w.SetError("the app has no access to priviliges ["+scopeArr[i]+"]", "")
			return false
		}
	}

	return true
}

func (oauth *OAuth) PostAuthorize(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Println("PostAuthorize:\r\n")

	apiChoose := r.FormValue("getUsername")
	fmt.Println("api_choose", apiChoose)

	acname := oauth.Logged(w, r)
	fmt.Println("acname=",acname)
	if acname == "" {
		//使用提交的表单登陆
		acname, _ = oauth.Login(w, r)
		fmt.Println("11acname=",acname)
		//登陆失败
		if acname == "" {
			//返回页面，出现 登陆失败提示，用户名密码框+授权并登陆按钮+权限列表
			oauth.View.HTML(w, http.StatusOK, "oauth", nil)
			return
		}
	}

	//用户登陆成功，并确认授权，则进行下一步,根据请求,发放code 或token
	resp := oauth.Server.NewResponse()
	defer resp.Close()
	ar := oauth.Server.HandleAuthorizeRequest(resp, r)
	if ar != nil {
		//发放code 或token ,附加到redirect_uri后，并跳转
		//存储acname，acid,rsid,clientid,clientSecret等必要信息
		//ar.UserData = struct{ Acname string }{Acname: acname}
		fmt.Println("Write Authorize Begin:")
		acid:=getAcId(acname)
		ar.UserData = ATUserData{Acname:acname,Acid:acid}
		ar.Authorized = true

		oauth.Server.FinishAuthorizeRequest(resp, r, ar)
		fmt.Println("Write Authorize End:")
	}
	osin.OutputJSON(resp, w, r)
}

func (oauth *OAuth) Token(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Println("Token:\r\n")
	resp := oauth.Server.NewResponse()
	defer resp.Close()

	if ar := oauth.Server.HandleAccessRequest(resp, r); ar != nil {
		switch ar.Type {
		case osin.AUTHORIZATION_CODE:
			ar.Authorized = true
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

		oauth.Server.FinishAccessRequest(resp, r, ar)
	}
	osin.OutputJSON(resp, w, r)
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
	strSQL := fmt.Sprintf("select count(ac_name) from account_tab where ac_name='%s' and ac_password='%s'", user.Acname,user.Password)
		fmt.Println(strSQL)
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
	acid := getAcId(user_name)
	if accessData.Client == nil {
		fmt.Println("Get Client Faild!!!")
	}
	client_id := accessData.Client.GetId()

	//fmt.Printf("acid=%d;client_id=%s\r\n",acid,client_id)
	jr := make(map[string]interface{})
	jr["client_id"] = client_id
	if acid != -1 && client_id != "" {
		openId := getOpenId("000001", client_id, acid)
		jr["openid"] = openId
	}

	result, err := json.Marshal(jr)
	if err != nil {
		result = []byte("")
	}
	w.Write(result)
}

func getAcId(acName string) int {
	strSQL := fmt.Sprintf("select ac_id from account_tab where ac_name='%s' limit 1", acName)
	mydb := common.GetDB()
	if mydb == nil {
		fmt.Println("get db connection error")
		return -1
	}
	defer common.FreeDB(mydb)
	rows, err := mydb.Query(strSQL)
	if err != nil {
		return -1
	} else {
		defer rows.Close()
		var acid int
		for rows.Next() {
			rows.Scan(&acid)
		}
		return acid
	}
	return -1
}

//func (oauth *OAuth)getOpenId(clientId string, acid int) string {
func getOpenId(res_id string, clientId string, acid int) string {
	strSQL := fmt.Sprintf("select openid from openid_tab where res_id='%s' and client_id='%s' and acid=%d limit 1", res_id, clientId, acid)
	//fmt.Println(strSQL)
	mydb := common.GetDB()
	if mydb == nil {
		fmt.Println("get db connection error")
		return ""
	}
	defer common.FreeDB(mydb)
	rows, err := mydb.Query(strSQL)
	if err != nil {
		return ""
	} else {
		defer rows.Close()
		var openId string
		for rows.Next() {
			rows.Scan(&openId)
		}
		return openId
	}
}

func (oauth *OAuth) CheckPrivilige(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Println("CheckPrivilige:\r\n")

	strToken := r.FormValue("token")
	strPrivilige := r.FormValue("privilige")
	if strToken == "" {
		strBody := []byte("{\"Code\":1,\"Message\":\"token can not be empty \"}")
		w.Write(strBody)
		return
	}
	if strPrivilige == "" {
		strBody := []byte("{\"Code\":1,\"Message\":\"privilige can not be empty \"}")
		w.Write(strBody)
		return
	}

	ret, err := oauth.Server.Storage.LoadAccess(strToken)
	fmt.Println(err)
	if err != nil {
		strBody := []byte("{\"Code\":1,\"Message\":\"user token not exist \"}")
		w.Write(strBody)
		return
	}

	fmt.Println("ret.Scope=",ret.Scope)
	fmt.Println("strPrivilige=",strPrivilige)
	if !strings.Contains(ret.Scope, strPrivilige) {
		strBody := []byte("{\"Code\":1,\"Message\":\"no privilige\"}")
		w.Write(strBody)
		return
	} else {
		strBody := []byte("{\"Code\":0,\"Message\":\"OK \"}")
		w.Write(strBody)
	}
}
