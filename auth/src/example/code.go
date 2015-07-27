package main

// Open url in browser:
// http://localhost:8080/

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/RangelReale/osin"
	"net/http"
	"net/url"
	"time"
)

// Client information
type Client struct {
	//type Client interface {
	// Client id
	Id          string
	Secret      string
	RedirectUri string
	UserData    interface{}
}

func (d *Client) GetId() string {
	return d.Id
}

func (d *Client) GetSecret() string {
	return d.Secret
}

func (d *Client) GetRedirectUri() string {
	return d.RedirectUri
}

func (d *Client) GetUserData() interface{} {
	return d.UserData
}

// DefaultClient stores all data in struct variables
type DefaultClient struct {
	Id          string
	Secret      string
	RedirectUri string
	UserData    interface{}
}

func (d *DefaultClient) GetId() string {
	return d.Id
}

func (d *DefaultClient) GetSecret() string {
	return d.Secret
}

func (d *DefaultClient) GetRedirectUri() string {
	return d.RedirectUri
}

func (d *DefaultClient) GetUserData() interface{} {
	return d.UserData
}

func main() {
	/*
		strClient:="{\"Id\":\"1234\",\"Secret\":\"aabbccdd\",\"RedirectUri\":\"http://localhost:8080\",\"UserData\":\"\"}"
		//strClient:="{\"Id\":\"1234\",\"Secret\":\"aabbccdd\",\"RedirectUri\":\"http://localhost:8080\",\"UserData\":\"\"}"
		var client Client
		err := json.Unmarshal([]byte(strClient), &client)
		fmt.Println(client)
		if err != nil {
			fmt.Println("Methdo - fromJSON : ", err)
		}
		return
	*/

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("<html><head><meta charset=\"utf-8\" /><title>第三方应用首页</title></head><body>"))
		w.Write([]byte("<h2>Step1：获取Authorization Code</h2>"))
		w.Write([]byte("<h4>1：检测用户是否在平台登录，未登录，跳转到登录页；</h4>"))
		w.Write([]byte("<h4>2：已登录，跳转到授权页；</h4>"))
		w.Write([]byte("<h4>3：授权项由参数scope决定；</h4>"))
		w.Write([]byte("<h4>4：此次申请权限：get_user_info,getUsername，list_photo，upload_pic</h4><br/>"))
		w.Write([]byte(fmt.Sprintf("<a href=\"https://connect.funzhou.cn/oauth2/authorize?response_type=code&client_id=1234&redirect_uri=%s&state=first&scope=get_user_info,getUsername,list_photo,upload_pic\">获取Authorization Code）</a><br/>", url.QueryEscape("http://localhost:8080/callback"))))
		//w.Write([]byte(fmt.Sprintf("<a href=\"https://connect.funzhou.cn/oauth2/authorize?response_type=code&client_id=1234&redirect_uri=%s&state=first&scope=get_user_info,getUsername\">获取Authorization Code）</a><br/>", url.QueryEscape("http://localhost:8080/callback"))))
		w.Write([]byte("</body></html>"))
	})

	// Application destination - CODE
	http.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("callback!!!!\r\n")
		r.ParseForm()

		w.Write([]byte("<html><body>"))
		w.Write([]byte("<h2>第三方登陆第二步</h2>"))
		w.Write([]byte("<h4>平台在redirect_uri后附加授权code，并返回HTTP302</h4>"))
		code := r.Form.Get("code")
		w.Write([]byte(fmt.Sprintf("<h4>由浏览器跳转到当前页并把code=%s随请求一起传递给此页的http.Request对象</h4>", code)))
		w.Write([]byte("<h4>此页的通过http.Request对象，取得授权code后，在后台访问平台Token接口(需使用HTTPBasicAuth协议传入client_id和client_secret)</h4>"))
		// build access code url
		aurl := fmt.Sprintf("https://connect.funzhou.cn/oauth2/token?grant_type=authorization_code&client_id=1234&state=xyz1&redirect_uri=%s&code=%s",
			url.QueryEscape("http://localhost:8080/callback"), url.QueryEscape(code))
		w.Write([]byte("<h4>" + aurl + "</h4>"))
		w.Write([]byte("<h4>取得如下信息</h4>"))
		jr := make(map[string]interface{})
		err := DownloadAccessToken(aurl, &osin.BasicAuth{"1234", "aabbccdd"}, jr)
		if err != nil {
			w.Write([]byte(err.Error()))
			w.Write([]byte("<br/>"))
		}
		// show json error
		if erd, ok := jr["error"]; ok {
			w.Write([]byte(fmt.Sprintf("ERROR: %s<br/>\n", erd)))
		}
		// show json access token
		token, ok := jr["access_token"]
		if ok {
			w.Write([]byte(fmt.Sprintf("授权令牌access_token: %s<br/>\n", token)))
		}
		if at, ok := jr["expires_in"]; ok {
			w.Write([]byte(fmt.Sprintf("有效期expires_in: %f秒<br/>\n", at)))
		}
		if at, ok := jr["refresh_token"]; ok {
			w.Write([]byte(fmt.Sprintf("续期刷新令牌refresh_token: %s<br/>\n", at)))
		}
		if at, ok := jr["scope"]; ok {
			w.Write([]byte(fmt.Sprintf("正式授予的权限scope: %s<br/>\n", at)))
		}
		w.Write([]byte(fmt.Sprintf("<br/><a href=./bindingAndlogin?access_token=%s>进入第三步（获取OpenID登陆)可和第二步合并</a><br/>", token)))
		defer w.Write([]byte("</body></html>"))
	})
	http.HandleFunc("/bindingAndlogin", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("bindingAndlogin!!!!\r\n")
		w.Write([]byte("<html><body>"))
		w.Write([]byte("<h2>第三方登陆第三步</h2>"))
		accessToken := r.FormValue("access_token")
		w.Write([]byte(fmt.Sprintf("<h4>第三方使用access_token去平台请求https://connect.funzhou.cn/oauth2/me?access_token=%s</h4>", accessToken)))

		w.Write([]byte("<h4>获取当前登陆用户的openid（每个应用不一致）,创建（存在就不创建）并绑定第三方自己账号并登陆</h4>"))

		jr := make(map[string]interface{})
		err := CallInterface("https://connect.funzhou.cn/oauth2/me", accessToken, jr)
		if err != nil {
			w.Write([]byte(err.Error()))
			w.Write([]byte("<br/>"))
		}
		// show json access token
		clientId, _ := jr["client_id"]
		openId, _ := jr["openid"]
		w.Write([]byte(fmt.Sprintf("<h4>相当于登陆openId:%s的用户名是openid%s", clientId, openId)))
		w.Write([]byte(fmt.Sprintf("<h4>相当于登陆openId:%s的密码名是access_token：%s", clientId, accessToken)))

		defer w.Write([]byte("</body></html>"))
	})

	fmt.Println("Server is start at ", time.Now().String(), " , on port 8080")
	http.ListenAndServe(":8080", nil)
	fmt.Println("Test end \r\n")
}

func DownloadAccessToken(url string, auth *osin.BasicAuth, output map[string]interface{}) error {
	// download access token
	preq, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	if auth != nil {
		preq.SetBasicAuth(auth.Username, auth.Password)
	}

	pclient := &http.Client{}
	presp, err := pclient.Do(preq)
	if err != nil {
		return err
	}

	if presp.StatusCode != 200 {
		return errors.New("Invalid status code")
	}
	jdec := json.NewDecoder(presp.Body)
	err = jdec.Decode(&output)
	return err
}

func CallInterface(url string, accessToken string, output map[string]interface{}) error {
	// download access token
	preq, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	preq.Header.Add("access_token", accessToken)

	pclient := &http.Client{}
	presp, err := pclient.Do(preq)
	if err != nil {
		return err
	}

	if presp.StatusCode != 200 {
		return errors.New("Invalid status code")
	}
	jdec := json.NewDecoder(presp.Body)
	err = jdec.Decode(&output)
	return err
}
