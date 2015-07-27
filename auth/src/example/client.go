package main

// Open url in browser:
// http://localhost:8080/

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/RangelReale/osin"
	"net/http"
	_ "net/url"
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

	http.HandleFunc("/client", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("callback!!!!\r\n")
		r.ParseForm()

		w.Write([]byte("<html><body>"))
		//		aurl := fmt.Sprintf("https://connect.funzhou.cn/oauth2/token?grant_type=client_credentials&client_id=1234&client_secret=aabbccdd&scope=get_user_info,getUsername,list_photo,upload_pic")

		aurl := fmt.Sprintf("https://connect.funzhou.cn/oauth2/token?grant_type=client_credentials&scope=get_user_info,getUsername,list_photo,upload_pic")
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
		defer w.Write([]byte("</body></html>"))
	})

	fmt.Println("Server is start at ", time.Now().String(), " , on port 8080")
	http.ListenAndServe(":8080", nil)
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
