package main

import (
	"auth/action"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"runtime"
	"time"
	"common"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Welcome!\n")
}
// ab -c 100 -n 1000 'http://127.0.0.1:8080/hello/tomzhao'
func Hello(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	acname := r.FormValue("acname")
	password := r.FormValue("password")
	fmt.Println("acname=",acname)
	fmt.Println("password=",password)
	fmt.Println("name=",ps.ByName("name"))

	strSQL := "SELECT ac_name FROM account_tab"
	mydb := common.GetDB()
	defer common.FreeDB(mydb)	

	rows, err := mydb.Query(strSQL)
	if err != nil {
	} else {
		defer rows.Close()
		var strAcName string
		for rows.Next() {
			rows.Scan(&strAcName)
		}
		fmt.Fprintf(w, "hello, %s at %s!!!!!!\n", strAcName,time.Now().String())
	} 	
}

func sysinfo(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	user_name := r.FormValue("user_name")
	password := r.FormValue("password")
	fmt.Println("user_name=",user_name)
	fmt.Println("password=",password)
	//strPostData := r.FormValue("request")
	//fmt.Println("strPostData :", strPostData)
	fmt.Fprintf(w, common.GetDBInfo())
}

type HostSwitch map[string]http.Handler

func (hs HostSwitch) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if handler := hs[r.Host]; handler != nil {
		handler.ServeHTTP(w, r)
	} else {
		http.Error(w, "Forbidden", 403) // Or Redirect?
	}
}

func new_router() *httprouter.Router {
	router := httprouter.New()
	router.GET("/", Index)
	router.GET("/hello/:name", Hello)
	router.POST("/sysinfo", sysinfo)
	router.POST("/auth/get_user_info", action.GetUserInfo)
	router.POST("/auth/register", action.RegisterAccount)
	router.POST("/auth/login", action.Login)
	router.POST("/auth/logout", action.Logout)
	router.POST("/auth/getacid", action.GetAcidByOpenid)
	router.POST("/auth/changepw", action.ChangePassword)
	return router
}

func main() {
	fmt.Println("Server is start at ", time.Now().String(), " , on port 8090")
	router := new_router()
	hs := make(HostSwitch)
	hs["127.0.0.1:8090"] = router

	log.Fatal(http.ListenAndServe(":8090", hs))
	//log.Fatal(http.ListenAndServeTLS(":8080", "../connect/static/pem/servercert.pem", "../connect/static/pem/serverkey.pem", hs))
}
