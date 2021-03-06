package common

import (
	"fmt"
	"github.com/robfig/config"
	"io/ioutil"
	"os"
	"net/http"
	"net/url"
	"strings"
)

func DisplayJson(obj_json map[string]interface{}) {
	Log().Println(obj_json)
	Log().Println("----------------------parse start------------------------")
	for k, v := range obj_json {
		switch v2 := v.(type) {
		case string:
			fmt.Println(k, "is string", v2)
		case int:
			fmt.Println(k, "is int ", v2)
		case bool:
			fmt.Println(k, "is bool", v2)
		case []interface{}:
			fmt.Println(k, "is array", v2)
			for i, iv := range v2 {
				fmt.Println(i, iv)
			}
		case map[string]interface{}:
			fmt.Println(k, "is map")
			DisplayJson(v2)
		default:
			fmt.Println(k, "is another type not handle yet")
		}
	}
	Log().Println("----------------------parse end------------------------")
}

const (
	SUCCESS     int    = 0
	FAILT       int    = 1
	SUCCESS_MSG string = "ok"
)


func SaveFile(strFileName string, strData string) (ok bool) {
	f, err := os.Create(strFileName)
	if err != nil {
		fmt.Println("create file faild error:", err)
		return false
	}
	_, err_w := f.Write([]byte(strData))
	if err_w != nil {
		fmt.Println("Server start faild error:", err_w)
		return false
	}
	return true
}

func Config() *config.Config {
	file, _ := os.Getwd()
	c, _ := config.ReadDefault(file + "/common/config.cfg")
	fmt.Println("Config init success ...")
	return c
}

//截取固定位置以前的字符串
func SubstrBefore(s string, l int) string {
	if len(s) <= l {
		return s
	}
	ret, rs := "", []rune(s)

	for i, r := range rs {
		if i >= l {
			break
		}

		ret += string(r)
	}
	return ret
}

//截取固定位置以后的字符串
func SubstrAfter(s string, l int) string {
	if len(s) <= l {
		return s
	}
	ret, rs := "", []rune(s)

	for i, r := range rs {
		if i <= l {
			continue
		}

		ret += string(r)
	}
	return ret
}

//判断文件夹是否存在
func IsDirExists(path string) bool {
	fi, err := os.Stat(path)

	if err != nil {
		return os.IsExist(err)
	} else {
		return fi.IsDir()
	}
}

//判断文件是否存在
func IsFileExists(path string) bool {
	fi, err := os.Stat(path)
	if err != nil {
		return os.IsExist(err)
	} else {
		return !fi.IsDir()
	}
}

//从文件中读取内容
func ReadFile(filePth string) string {
	bytes, err := ioutil.ReadFile(filePth)
	if err != nil {
		fmt.Println("读取文件失败: ", err)
		return ""
	}

	return string(bytes)
}

func HttpGet(strURL string)(strBody string,err error) {
	resp, err := http.Get(strURL)
	if err != nil {
		// handle error
		return "",err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "",err
	}

	return string(body),nil
}

func HttpPost(strURL,strPostData string) (strBody string,err error) {
	resp, err := http.Post(strURL,
		"application/x-www-form-urlencoded", strings.NewReader(strPostData))
	if err != nil {
		fmt.Println(err)
		return "",err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "",err
	}

	return string(body),nil
}

func HttpPostForm(strURL string,uv url.Values) (strBody string,err error) {
	resp, err := http.PostForm(strURL,uv)

	if err != nil {
		return "",err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "",err
	}

	return string(body),nil
}

func Invoker(invoke_type int,invoke_dest string,invoke_data interface{})(strBody string,err error) {
	switch invoke_type {
	case HTTP_GET:
		strBody,err=HttpGet(invoke_dest)
	case HTTP_POST:
		//fmt.Println(invoke_data.(type))
		switch invoke_data.(type) {
		case string:
			strBody,err=HttpPost(invoke_dest,invoke_data.(string))
		case url.Values:
			strBody,err=HttpPostForm(invoke_dest,invoke_data.(url.Values))
		}
	case HTTPS_GET:
		fmt.Println("HTTPS_GET")
	case HTTPS_POST:
		fmt.Println("HTTPS_POST")
	default:
		fmt.Println("is another type not handle yet")
	}
	return strBody,err
}
/*
	strBody,err:=common.Invoker(common.HTTP_GET,"http://127.0.0.1:8090/hello/tomzhao","")
	//strTemp := "request=test" 
	//strBody,err:=common.Invoker(common.HTTP_POST,"http://127.0.0.1:8090/sysinfo",value)
	//value:=url.Values{"user_name": {"tomzhao"}, "password": {"111111"}}
	//strBody,err:=common.Invoker(common.HTTP_POST,"http://127.0.0.1:8090/sysinfo",value)
	if err!=nil {
		fmt.Println(err)
		return
	}
	fmt.Println(strBody)
	return 
*/