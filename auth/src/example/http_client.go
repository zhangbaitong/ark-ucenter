package main

import (
	"fmt"
	"net/url"
	"common"
	"gopkg.in/mgo.v2/bson" 
	_"encoding/json"
	_"encoding/base64"
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
	var Url *url.URL
	Url, err := url.Parse("http://sms.gyjbh.nvwayun.com")
	if err != nil {
	    panic("boom")
	}

	mesaage:=make(map[string]string)
	mesaage["mobile"]="18585816540"
	mesaage["msg"]="888888(动态验证码),请在5分钟内输入该验证码."
	strData, err := json.Marshal(mesaage)
	if err != nil {
		return
	}

	Url.Path += "/application/api"
	strTnterfaceKey:="ad79bd61-4cc8-f4a4-2811-55e0117e6cc4"
	strInterfaceSign:="4bf38c7e184df4087910038afc7df8b9b899aa2f"
	strSend:=base64.StdEncoding.EncodeToString(strData)
	parameters := url.Values{}
	parameters.Add("data", strSend)
	parameters.Add("interface_key", strTnterfaceKey)
	parameters.Add("interface_sign", strInterfaceSign)
	Url.RawQuery = parameters.Encode()
	fmt.Printf("Encoded URL is %q\n", Url.String())

	strResult,err:=common.Invoker(common.HTTP_GET,Url.String(),"")
	if err!=nil {
		return
	}

	fmt.Println(strResult)

	//lue:=url.Values{"data_type":{"1"},"start_time": {"1439395200"},"end_time":{"1439740800"}}
	value:=url.Values{"data_type":{"1"},"start_time": {"1441078590"},"end_time":{"1441078590"}}
	strBody,err:=common.Invoker(common.HTTP_POST,"https://connect.funzhou.cn/manage/export_data",value)
	if err!=nil {
		fmt.Println(err)
		return
	}
	fmt.Println(strBody)	

	strVerifyCode:="912297" 
	value:=url.Values{"mobile": {"18585816540"},"verify_code":{strVerifyCode},"show_name":{"tomtest"}}
	//value:=url.Values{"mobile": {"15519028660"},"verify_code":{strVerifyCode}}
	//value:=url.Values{"mobile": {"18585816540"},"verify_code":{strVerifyCode}}
	//value:=url.Values{"mobile": {"18984550575"},"verify_code":{strVerifyCode}}
	//value:=url.Values{"mobile": {"18585816540"},"verify_code":{strVerifyCode}}
	strBody,err:=common.Invoker(common.HTTP_POST,"https://connect.funzhou.cn/user/get_verify_code",value)
	//strBody,err:=common.Invoker(common.HTTP_POST,"https://connect.funzhou.cn/user/check_verify_code",value)
	if err!=nil {
		fmt.Println(err)
		return
	}
	fmt.Println(strBody)	

	value:=url.Values{"openid": {"55e170e1e138230f4zhw"}}
	//value:=url.Values{"id": {"55cabc1de1382334ec000006"},"need_send_sms":{"0"}}
	//value:=url.Values{"id": {"55e14babe138232c79000001"}}
	//value:=url.Values{"mobile": {"18585816666"}}
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

	//password:=common.MD5("111111")
	//value:=url.Values{"acname":{"18984550001"},"password":{password},"name": {"zhanghuawen"},"company_name": {"infobird"},"job":{"yanfa"},"company_addr":{""},"email":{"18984550000@qq.com"},"mobile":{"18984550000"}}
	//value:=url.Values{"acname":{"18984550006"},"password":{password},"name": {"zhanghuawen"},"source_id": {"funzhou_0001"},"company_name": {"infobird"},"job":{"yanfa"},"company_addr":{""},"email":{"18984550006@qq.com"},"mobile":{"18984550006"}}
	//value:=url.Values{"name": {"zhanghuawen"},"source_id": {"funzhou_0001"},"company_name": {"infobird"},"job":{"yanfa"},"company_addr":{""},"email":{"18984550006@qq.com"},"mobile":{"18984550006"}}
	//value:=url.Values{"reg_type":{"1"},"show_name":{"tomtest"},"name": {"zhanghuawen"},"source_id": {"funzhou_0001"},"company_name": {"infobird"},"job":{"yanfa"},"haha11":{"yanfa"},"company_addr":{""},"email":{"18984550008@qq.com"},"mobile":{"18984550008"}}
	value:=url.Values{"reg_type":{"1"},"show_name":{"tomtest"}}
	strBody,err:=common.Invoker(common.HTTP_POST,"https://connect.funzhou.cn/user/register",value)
	if err!=nil {
		fmt.Println(err)
		return
	}
	fmt.Println(strBody)	

	// "email": "zhw@sina.com",
	value:=url.Values{"mobile": {"18984550575"}}
	strBody,err:=common.Invoker(common.HTTP_POST,"https://connect.funzhou.cn/user/exist",value)
	if err!=nil {
		fmt.Println(err)
		return
	}
	fmt.Println(strBody)	
	password:=common.MD5("123456")
	new_password:=common.MD5("123456")
	value:=url.Values{"acname": {"gy_01@qq.com"},"password":{password},"new_password":{new_password}}
	//strBody,err:=common.Invoker(common.HTTP_POST,"https://connect.funzhou.cn/user/register",value)
	//strBody,err:=common.Invoker(common.HTTP_POST,"https://connect.funzhou.cn/user/change_password",value)
	strBody,err:=common.Invoker(common.HTTP_POST,"https://connect.funzhou.cn/manage/reset_password",value)
	if err!=nil {
		fmt.Println(err)
		return
	}
	fmt.Println(strBody)	

	value:=url.Values{"fieldes": {"qq,email,mobile,weibo,openid"}}
	strBody,err:=common.Invoker(common.HTTP_GET,"https://connect.funzhou.cn/manage/get_only_check_list",value)
	//strBody,err:=common.Invoker(common.HTTP_POST,"https://connect.funzhou.cn/manage/update_only_check_list",value)
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
	fmt.Println(strBody)	

	//value:=url.Values{"id": {"561b1f35e138232394000001"},"show_name":{"tomtest"},"name": {"zhanghuawen"},"source_id": {"funzhou_0001"},"company_name": {"infobird"},"job":{"yanfa"},"haha":{"yanfa66"},"company_addr":{""},"email":{"18984550008@qq.com"},"mobile":{"18984550008"}}
	value:=url.Values{"id": {"561b1f35e138232394000001"},"show_name":{"tomtest"},"name": {"zhanghuawen111"},"source_id": {"funzhou_0001"},"company_name": {"infobird"},"job":{"yanfa1111"},"email":{"18984550008@qq.com"},"mobile":{"18984550008"}}
	strBody,err:=common.Invoker(common.HTTP_POST,"https://connect.funzhou.cn/user/set_user_info",value)
	if err!=nil {
		fmt.Println(err)
		return
	}
	fmt.Println(strBody)	
	
	//value:=url.Values{"id": {"561b1f35e138232394000001"},"show_name":{"tomtest"}}
	value:=url.Values{"haha": {"yanfa66"},"show_name":{"tomtest"}}
	strBody,err:=common.Invoker(common.HTTP_POST,"https://connect.funzhou.cn/user/get_user_info",value)
	if err!=nil {
		fmt.Println(err)
		return
	}
	fmt.Println(strBody)	

	//value:=url.Values{"start_time": {"1441728000"},"end_time": {"1441814400"}}
	//value:=url.Values{"start_time": {"1441814400"},"end_time": {"1441900800"}}
	//value:=url.Values{"start_time": {"1441900800"},"end_time": {"1441987200"}}
	//value:=url.Values{"start_time": {"1441987200"},"end_time": {"1442073600"}} 
	value:=url.Values{"data_type": {"1"},"start_time": {"1441728000"},"end_time": {"1442073600"}} 
	//strBody,err:=common.Invoker(common.HTTP_POST,"https://connect.funzhou.cn/manage/user_stat",value)
	strBody,err:=common.Invoker(common.HTTP_POST,"https://connect.funzhou.cn/manage/export_data",value)
	if err!=nil {
		fmt.Println(err)
		return
	}
	fmt.Println(strBody)	
*/
	value:=url.Values{"user_list": {"561b1f35e138232394000001,56188ea3e1382317b3000001"}}
	strBody,err:=common.Invoker(common.HTTP_POST,"https://connect.funzhou.cn/user/get_user_list",value)
	if err!=nil {
		fmt.Println(err)
		return
	}
	fmt.Println(strBody)	


}
