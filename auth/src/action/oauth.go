package action

import (
	"common"
	"encoding/json"
	"fmt"
	"github.com/RangelReale/osin"
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

	Res struct {
		Resname  string
		Rescname string
	}
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
	acname := UserLogin(w, r)
	if acname == "" {
		return
	}
	if !checkAuthorize(oauth, w, r, acname) {
		return
	}
	doAuthorizeRequest(oauth, w, r)
}

//检查是否登录，未登录，则返回登录页
func checkLogin(oauth *OAuth, w http.ResponseWriter, r *http.Request) bool {
	fmt.Println("checkLogin\r\n")
	acname := GetCookieName(r)
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
		acname = GetCookieName(r)
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
		acname := GetCookieName(r)
		if acname == "" {
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
		ok := LoginQuery(&user)
		if ok {
			GenerateCookie(w, r, user.Acname, 1)
			ar.Authorized = true
		} else {
			//通过redirect_uri 返回错误约定 并跳转到改redirect_uri
		}
	case osin.CLIENT_CREDENTIALS:
		//校验appId和appKey是否正确
		if ar.Client.GetSecret() != GetAppKey(ar.Client.GetId()) {
			ar.Authorized = false
			return ar
		}

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
			}
		}
		//重新给token绑定资源
		ar.Scope = resources
	case osin.ASSERTION:
		if ar.AssertionType == "urn:osin.example.complete" && ar.Assertion == "osin.data" {
			ar.Authorized = true
		}
	}

	return ar
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
	resp := oauth.Server.NewResponse()

	flag := true
	queryUrl := common.GetUrlParam(r)
	clientId := queryUrl["client_id"][0]

	token := queryUrl["token"][0]
	storage, err := oauth.Server.Storage.LoadAccess(token)

	if err != nil {
		fmt.Println("get token storage failure")
		flag = false
	} else {
		if storage.CreatedAt.Add(time.Duration(3600) * time.Second).Before(oauth.Server.Now()) {
			flag = false
			resp.SetError("invalid_grant test", "")
		} else {
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
		}
	}

	if flag {
		resp.Output["code"] = 0
		resp.Output["data"] = storage.Scope
	} else {
		resp.Output["code"] = 1
	}
	//	common.Write(w, ret)
	osin.OutputJSON(resp, w, r)

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
