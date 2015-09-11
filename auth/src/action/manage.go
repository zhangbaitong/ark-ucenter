package action

import (
	_"common"
	"fmt"
	 _"gopkg.in/mgo.v2/bson" 
	"strings"
	"net/http"
	"github.com/julienschmidt/httprouter"
	"strconv" 
)

func GetOnlyCheckList(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Println("GetOnlyCheckList:")
	FieldList:=GetCheckList()
	strList:=""
	for i := 0; i < len(FieldList); i++ {
		if i==0 {
			strList=FieldList[i]
		} else{
			strList=strList+","+FieldList[i]
		}
	}
	strMessage:="{\"Code\":0,\"Message\":\""+strList+"\"}"
	w.Write([]byte(strMessage))
}

func UpdateOnlyCheckList(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	strFieldes := r.FormValue("fieldes")
	FieldList:=strings.Split(strFieldes, ",")

	strBody := []byte("{\"Code\":0,\"Message\":\"ok\"}")	
	ok := SetOnlyCheckList(FieldList)
	if !ok {
		strBody = []byte("{\"Code\":1,\"Message\":\"save data error\"}")
	}
	w.Write(strBody)
}

func PasswordReset(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	req.ParseForm()
	fmt.Println(req.Form)

	acname := req.FormValue("acname")
	password := req.FormValue("password")
	strBody:=[]byte("")
	if acname == "" || password == "" {
		strBody = []byte(fmt.Sprintf("{\"Code\":%d,\"Message\":\"%s\"}",PARAM_ERROR,GetError(PARAM_ERROR)))
		w.Write(strBody)
		return
	}
	//password=common.MD5(password)
	ok := ResetPassword(acname,password)
	if ok {
		strBody = []byte("{\"Code\":0,\"Message\":\"ok\"}")
	} else {
		strBody = []byte(fmt.Sprintf("{\"Code\":%d,\"Message\":\"%s\"}",UPDATE_DB_ERROR,GetError(UPDATE_DB_ERROR)))
	}
	w.Write(strBody)
}

func ExportData(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	req.ParseForm()
	fmt.Println(req.Form)

	data_type := req.FormValue("data_type")
	start_time := req.FormValue("start_time")
	end_time := req.FormValue("end_time")
	strBody:=[]byte("")
	if data_type=="" || start_time == "" || end_time == "" {
		strBody = []byte(fmt.Sprintf("{\"Code\":%d,\"Message\":\"%s\"}",PARAM_ERROR,GetError(PARAM_ERROR)))
		w.Write(strBody)
		return
	}

	StartTime, _ := strconv.Atoi(start_time)
	EndTime, _ := strconv.Atoi(end_time)

	//password=common.MD5(password)
	user_count,ok := ExportMongo(StartTime,EndTime)
	if ok {
		strBody = []byte(fmt.Sprintf("{\"Code\":%d,\"Message\":\"user_count=%d\"}",OK,user_count))
	} else {
		strBody = []byte(fmt.Sprintf("{\"Code\":%d,\"Message\":\"%s\"}",UPDATE_DB_ERROR,GetError(UPDATE_DB_ERROR)))
	}
	w.Write(strBody)
}

func UserStat(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	req.ParseForm()
	fmt.Println(req.Form)

	start_time := req.FormValue("start_time")
	end_time := req.FormValue("end_time")
	strBody:=[]byte("")
	if start_time == "" || end_time == "" {
		strBody = []byte(fmt.Sprintf("{\"Code\":%d,\"Message\":\"%s\"}",PARAM_ERROR,GetError(PARAM_ERROR)))
		w.Write(strBody)
		return
	}

	StartTime, _ := strconv.Atoi(start_time)
	EndTime, _ := strconv.Atoi(end_time)
	
	reg_bat,_:=RegUserStat(StartTime,EndTime,0)
	reg_mobile,_:=RegUserStat(StartTime,EndTime,1)
	strBody = []byte(fmt.Sprintf("{\"Code\":%d,\"Message\":\"reg_bat=%d;reg_mobile=%d\"}",OK,reg_bat,reg_mobile))
	w.Write(strBody)
}