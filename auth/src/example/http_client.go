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
	// "email": "zhw@sina.com",
	value:=url.Values{"email": {"zhw11@sina.com"}}
	strBody,err:=common.Invoker(common.HTTP_POST,"https://connect.funzhou.cn/user/exist",value)
	if err!=nil {
		fmt.Println(err)
		return
	}
	fmt.Println(strBody)	
/*
	value:=url.Values{"acname": {"gyzly@tom.com"},"password":{"111111"},"email":{"gyzly@tom.com"}}
	strBody,err:=common.Invoker(common.HTTP_POST,"https://connect.funzhou.cn/user/register",value)
	if err!=nil {
		fmt.Println(err)
		return
	}
	fmt.Println(strBody)	

	value:=url.Values{"acname": {"zhw"},"password":{"111111"},"email":{"zhw@sina.com"}}
	strBody,err:=common.Invoker(common.HTTP_POST,"https://connect.funzhou.cn/user/register",value)
	if err!=nil {
		fmt.Println(err)
		return
	}
	fmt.Println(strBody)	

	value:=url.Values{"acname": {"zhw"},"password":{"111111"}}
	strBody,err:=common.Invoker(common.HTTP_POST,"https://connect.funzhou.cn/user/login",value)
	if err!=nil {
		fmt.Println(err)
		return
	}
	var result Response
	json.Unmarshal([]byte(strBody),&result)
	fmt.Println(result)	
	fmt.Println(strBody)	

	value:=url.Values{"id": {"55c1804be13823298d000001"},"name":{"张华文"},"mobile":{"18585816500"}}
	strBody,err:=common.Invoker(common.HTTP_POST,"https://connect.funzhou.cn/user/set_user_info",value)
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

	value:=url.Values{"user_list": {"55c0891de1382334bb000002,55c1804be13823298d000001"}}
	strBody,err:=common.Invoker(common.HTTP_POST,"https://connect.funzhou.cn/user/get_user_list",value)
	if err!=nil {
		fmt.Println(err)
		return
	}
	var result Response
	json.Unmarshal([]byte(strBody),&result)
	fmt.Println(result)	
	fmt.Println(strBody)	
*/
}
