package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"encoding/json"
)

type RequestData struct {
	Version  string
	Method   string
	Params   string
}

func httpGet() {
	resp, err := http.Get("http://127.0.0.1:8080/v1/version")
	if err != nil {
		// handle error
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// handle error
	}

	fmt.Println(string(body))
}
func httpPost() {
	//var data map[string]string
	//data["imageName"]="centos:latest"
	//data["containerName"]="tomzhao"
	//post_data := RequestData{Version: "1.0", Method: "/auth/register", Params: "{\"Ac_name\":\"tomzhao1111\",\"Ac_password\":\"111111\",\"Email\":\"gyzly@sina.com1111\",\"Mobile\":\"18585816541\"}"}
	//post_data := RequestData{Version: "1.0", Method: "/auth/login", Params: "{\"user_name\":\"tomzhao\",\"password\":\"111111\"}"}
	//post_data := RequestData{Version: "1.0", Method: "/auth/logout", Params: "{\"ac_name\":\"tomzhao\",\"ac_password\":\"111111\"}"}
	//post_data := RequestData{Version: "1.0", Method: "/auth/getacid", Params: "{\"openid\":\"tomzhao\"}"}
	post_data := RequestData{Version: "1.0", Method: "/auth/changepw", Params: "{\"ac_name\":\"tomzhao\",\"old_password\":\"333333\",\"new_password\":\"444444\"}"}
	strPostData, _ := json.Marshal(post_data)
	strTemp := "request=" + string(strPostData)
	resp, err := http.Post("http://127.0.0.1:8080/auth/changepw",
		"application/x-www-form-urlencoded", strings.NewReader(strTemp))
	//"application/json",strings.NewReader(strTemp))
	if err != nil {
		fmt.Println(err)
	}

	defer resp.Body.Close()
	// body, err := ioutil.ReadAll(resp.Body)
	fmt.Println("-----------", resp.Header.Get("Set-Cookie"))
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// handle error
	}

	fmt.Println(string(body))
}

func httpPostForm() {
	resp, err := http.PostForm("http://127.0.0.1:8080/v1/version",
		url.Values{"key": {"Value"}, "id": {"123"}})

	if err != nil {
		// handle error
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// handle error
	}

	fmt.Println(string(body))

}

func httpDo() {
	client := &http.Client{}

	req, err := http.NewRequest("POST", "http://127.0.0.1:8080/v1/version", strings.NewReader("name=cjb"))
	if err != nil {
		// handle error
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Cookie", "name=anny")

	resp, err := client.Do(req)

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// handle error
	}

	fmt.Println(string(body))
}

//
type Client struct {
	Id          string
	Secret      string
	RedirectUri string
	UserData    interface{}
}

func fromJSON(jsonBytes string, obj interface{}) bool {
	err := json.Unmarshal([]byte(jsonBytes), &obj)
	if err != nil {
		fmt.Println("Methdo - fromJSON : ", err)
		return false
	}
	return true
}

func toJSON(obj interface{}) string {
	ret, err := json.Marshal(obj)
	if err != nil {
		fmt.Println("Methdo - toJSON : ", err)
		return ""
	}
	return string(ret)
}

func main() {

/*
*/
	//strClient:="{\"Id\":\"1234\",\"Secret\":\"aabbccdd\",\"RedirectUri\":\"http://localhost:8080\",\"UserData\":\"\"}"
	var client Client
	client.Id="1234"
	client.Secret="aabbccdd"
	client.RedirectUri="http://localhost:8080"
	client.UserData=""
	strTem:=toJSON(client)
	fmt.Println("strTem=", strTem)

	var client_test Client
	strClient:="{\"Id\":\"1234\",\"Secret\":\"aabbccdd\",\"RedirectUri\":\"http://localhost:8080\",\"UserData\":\"\"}"
	fmt.Println("strClient=", strClient)
	fromJSON(strClient,&client_test)
	fmt.Println(client_test)
	//strImage,strTag:=GetImage("10.122.75.228:5000/centostest:latest")
	//fmt.Println("strImage=", strImage)
	//fmt.Println("strTag=", strTag)
	//return
	//httpGet()
	//httpPost()
	//httpPostForm()
	//httpDo()
}
