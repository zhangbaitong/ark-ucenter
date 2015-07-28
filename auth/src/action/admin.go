package action

import (
	_"common"
	"fmt"
	 _"gopkg.in/mgo.v2/bson" 
	"strings"
	"net/http"
	"github.com/julienschmidt/httprouter"
)
//申请获取授权码
func GetSearchFieldList(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	FieldList:=GetSearchFieldes()
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
	return 
}

func UpdateSearchFieldList(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	strFieldes := r.FormValue("fieldes")
	fmt.Println("UpdateSearchFieldList:",strFieldes)
	FieldList:=strings.Split(strFieldes, ",")

	strBody := []byte("{\"Code\":0,\"Message\":\"ok\"}")	
	ok := SetSearchFieldes(FieldList)
	if !ok {
		strBody = []byte("{\"Code\":1,\"Message\":\"save data error\"}")
	}
	w.Write(strBody)
}

