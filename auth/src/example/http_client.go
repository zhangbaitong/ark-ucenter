package main

import (
	"fmt"
	"net/url"
	"common"
	"gopkg.in/mgo.v2/bson" 
	_"encoding/json"
)

type Response struct {
	Code int
	Message string
}

type ATUserInfo struct {
      Id bson.ObjectId "_id"
	Ac_id   int
	Info map[string] string
}

func main() {
/*
	//value:=url.Values{"acname": {"yanglinkang@qq.com "}}
	//value:=url.Values{"id": {"55cabc1de1382334ec000006"},"need_send_sms":{"0"}}
	value:=url.Values{"id": {"55cabb3ae1382334ec000004"}}
	//value:=url.Values{"email": {"zhw@sina.com"}}
	//value:=url.Values{"company": {"infobird1"}}
	strBody,err:=common.Invoker(common.HTTP_POST,"https://connect.funzhou.cn/user/get_user_info",value)
	if err!=nil {
		fmt.Println(err)
		return
	}
	var result Response
	json.Unmarshal([]byte(strBody),&result)
	fmt.Println(result)	
	fmt.Println(strBody)	

	value:=url.Values{"reg_type":{"1"},"company_name": {"infobird"},"job":{"yanfa"},"company_addr":{""},"email":{"tomtom33@infobird.com"},"mobile":{"1382881522555"}}
	strBody,err:=common.Invoker(common.HTTP_POST,"https://connect.funzhou.cn/user/register",value)
	if err!=nil {
		fmt.Println(err)
		return
	}
	fmt.Println(strBody)	

	// "email": "zhw@sina.com",
	value:=url.Values{"email": {"zhw11@sina.com"}}
	strBody,err:=common.Invoker(common.HTTP_POST,"https://connect.funzhou.cn/user/exist",value)
	if err!=nil {
		fmt.Println(err)
		return
	}
	fmt.Println(strBody)	
*/
	password:=common.MD5("333333")
	new_password:=common.MD5("333333")
	value:=url.Values{"acname": {"tomzhao44@qq.com "},"password":{password},"new_password":{new_password}}
	//strBody,err:=common.Invoker(common.HTTP_POST,"https://connect.funzhou.cn/user/register",value)
	strBody,err:=common.Invoker(common.HTTP_POST,"https://connect.funzhou.cn/user/change_password",value)
	//strBody,err:=common.Invoker(common.HTTP_POST,"https://connect.funzhou.cn/manage/reset_password",value)
	if err!=nil {
		fmt.Println(err)
		return
	}
	fmt.Println(strBody)	
/*
	value:=url.Values{"acname": {"zhw"},"password":{"111111"},"email":{"zhw@sina.com"}}
	strBody,err:=common.Invoker(common.HTTP_POST,"https://connect.funzhou.cn/user/register",value)
	if err!=nil {
		fmt.Println(err)
		return
	}
	fmt.Println(strBody)	

	password:=common.MD5("33335")
	value:=url.Values{"acname": {"muling@qq.com"},"password":{password}}
	strBody,err:=common.Invoker(common.HTTP_POST,"https://connect.funzhou.cn/user/login",value)
	if err!=nil {
		fmt.Println(err)
		return
	}
	var result Response
	json.Unmarshal([]byte(strBody),&result)
	fmt.Println(result)	
	fmt.Println(strBody)	

	value:=url.Values{"id": {"55c9e63ae138231dac000011"},"name":{"张华文"}}
	//strBody,err:=common.Invoker(common.HTTP_POST,"https://connect.funzhou.cn/user/set_user_info",value)
	strBody,err:=common.Invoker(common.HTTP_POST,"http://127.0.0.1:8080/user/set_user_info",value)
	if err!=nil {
		fmt.Println(err)
		return
	}
	var result Response
	json.Unmarshal([]byte(strBody),&result)
	fmt.Println(result)	
	fmt.Println(strBody)	
	
	value:=url.Values{"id": {"55c0891de1382334bb000002"}}
	strBody,err:=common.Invoker(common.HTTP_POST,"https://connect.funzhou.cn/user/get_user_info",value)
	if err!=nil {
		fmt.Println(err)
		return
	}
	var result Response
	json.Unmarshal([]byte(strBody),&result)
	fmt.Println(result)	
	fmt.Println(strBody)	

	value:=url.Values{"user_list": {"55c0891de1382334bb000002,55c1804be13823298d000331"}}
	strBody,err:=common.Invoker(common.HTTP_POST,"https://connect.funzhou.cn/user/get_user_list",value)
	if err!=nil {
		fmt.Println(err)
		return
	}
	var result Response
	json.Unmarshal([]byte(strBody),&result)
	fmt.Println(result)	
	//fmt.Println(strBody)	
*/
}
