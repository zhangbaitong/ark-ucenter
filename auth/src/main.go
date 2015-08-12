package main

import (
	"common"
	"action"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"time"
)

func SayHello(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	strTemp:=fmt.Sprintf("db_pool MaxPoolSize=%d;PoolSize=%d",common.DBpool.MaxPoolSize,common.DBpool.PoolSize)
	w.Write([]byte(strTemp))
}
func main() {
	oauth := action.NewOAuth()
	resource := action.NewResource()
	action.InitSearchFieldes()
	
	router := httprouter.New()
	router.GET("/", SayHello)
	router.POST("/", SayHello)
	router.NotFound = http.FileServer(http.Dir("./static/public")).ServeHTTP

	//Step1：获取Authorization Code
	router.POST("/oauth2/login", oauth.PostLogin)
	router.GET("/oauth2/authorize", oauth.GetAuthorize)
	router.POST("/oauth2/authorize", oauth.PostAuthorize)
	//Step2：通过Authorization Code获取Access Token
	//Step3：（可选）权限自动续期，获取Access Token
	router.GET("/oauth2/token", oauth.Token)
	//	router.POST("/oauth2/token", oauth.Token)
	//Step4:通过Access Token获取用户OpenID
	router.GET("/oauth2/me", oauth.Get)
	router.GET("/oauth2/queryPersonRes", oauth.QueryPersonResList)
	router.GET("/oauth2/privilige", oauth.CheckPrivilige)

	router.GET("/res/queryResId", resource.QueryResId)
	router.GET("/res/queryResCname", resource.QueryResCname)
	router.GET("/res/queryResByAppId", resource.QueryResByAppId)
	router.GET("/res/queryResByResName", resource.QueryResByResName)
	router.GET("/res/addResource", resource.AddResource)
	router.GET("/res/modifyResourceStatus", resource.ModifyResourceStatus)

	router.POST("/user/login", action.LoginCenter)
	router.GET("/user/logout", action.Logout)
	router.POST("/user/logout", action.Logout)
	router.POST("/user/exist", action.Exist)
	router.POST("/user/register", action.Register)
	router.POST("/user/set_user_info", action.SetUserInfo)
	router.POST("/user/get_user_info", action.GetUserInfo)
	router.POST("/user/get_user_list", action.GetUserList)

	router.GET("/admin/get_search_fieldes", action.GetSearchFieldList)
	router.POST("/admin/update_search_fieldes", action.UpdateSearchFieldList)
/*
	go func() {
		//start https server
		fmt.Println("Http Server is start at ", time.Now().String(), " , on port 80")
		log.Fatal(http.ListenAndServe(":80",  router))
	}()
*/		
	//start https server
	fmt.Println("Https Server is start at ", time.Now().String(), " , on port 443")
	log.Fatal(http.ListenAndServeTLS(":443", "./static/pem/servercert.pem", "./static/pem/serverkey.pem", router))
}
