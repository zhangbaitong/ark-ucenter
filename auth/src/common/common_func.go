package common

import (
	"fmt"
	"github.com/robfig/config"
	"html/template"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
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

//重定向到页面
func ForwardPage(w http.ResponseWriter, pageName string, data interface{}) {
	t, err := template.ParseFiles(pageName)
	if err != nil {
		fmt.Println("ForwardPage error ", err.Error())
		return
	}
	err = t.Execute(w, data)
	if err != nil {
		fmt.Println("ForwardPage error ", err.Error())
		return
	}
}

//获取http URL参数
func GetUrlParam(r *http.Request) map[string][]string {
	queryForm, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		fmt.Println("url ParseQuery error ", err.Error())
		return nil
	}
	return queryForm
}
