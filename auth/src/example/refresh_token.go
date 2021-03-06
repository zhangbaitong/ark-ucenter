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
	// Application destination - REFRESH
	http.HandleFunc("/appauth/refresh", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()

		w.Write([]byte("<html><body>"))
		w.Write([]byte("APP AUTH - REFRESH<br/>"))
		defer w.Write([]byte("</body></html>"))

		jr := make(map[string]interface{})

		err := DownloadAccessToken(fmt.Sprintf("https://connect.funzhou.cn/oauth2/token?grant_type=refresh_token&refresh_token=4MDk1xa-SuGn-1JnAjDSkA"),
			&osin.BasicAuth{Username: "1234", Password: "aabbccdd"}, jr)
		if err != nil {
			w.Write([]byte(err.Error()))
			w.Write([]byte("<br/>"))
		}

		// show json error
		if erd, ok := jr["error"]; ok {
			w.Write([]byte(fmt.Sprintf("ERROR: %s<br/>\n", erd)))
		}

		// show json access token
		if at, ok := jr["access_token"]; ok {
			w.Write([]byte(fmt.Sprintf("ACCESS TOKEN: %s<br/>\n", at)))
		}

		w.Write([]byte(fmt.Sprintf("FULL RESULT: %+v<br/>\n", jr)))

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
		fmt.Println(" presp.StatusCode", presp.StatusCode)
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
