package main

import (
	"action"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"time"
)

func SayHello(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	w.Write([]byte("Hello"))
}
func main() {
	oauth := action.NewOAuth()
	resource := action.NewResource()
	//me := action.NewMe()

	router := httprouter.New()
	router.GET("/", SayHello)
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

	router.GET("/user/logout", action.Logout)
	router.POST("/user/logout", action.Logout)
	router.POST("/user/register", action.Register)
	router.POST("/user/multi_register", action.RegisterMulti)
	router.POST("/user/set_user_info", action.SetUserInfo)
	router.POST("/user/multi_login", action.MultiLogin)
	router.POST("/user/get_user_info", action.SetUserInfo)


	router.GET("/admin/get_search_fieldes", action.GetSearchFieldList)
	router.POST("/admin/update_search_fieldes", action.UpdateSearchFieldList)

	fmt.Println("Server is start at ", time.Now().String(), " , on port 443")
	//log.Fatal(http.ListenAndServeTLS(":443", "./static/pem/servercert.pem", "./static/pem/serverkey.pem", router))
	log.Fatal(http.ListenAndServe(":8080",  router))
}
