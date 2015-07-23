package action

import (
	"common"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
)

func NewResource() *Resource {
	return new(Resource)
}

//通过res_name获取res_id
func (resource *Resource) QueryResId(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	queryUrl := common.GetUrlParam(r)
	resId := GetResId(queryUrl["res_name"][0])
	ret := make(map[string]interface{})
	ret["res_id"] = resId
	common.Write(w, ret)
}

//通过res_name查询res_cname
func (resource *Resource) QueryResCname(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	queryUrl := common.GetUrlParam(r)
	resCname := GetResCname(queryUrl["res_name"][0])
	ret := make(map[string]interface{})
	ret["res_cname"] = resCname
	common.Write(w, ret)
}

//通过app_id查询资源记录
func (resource *Resource) QueryResByAppId(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	queryUrl := common.GetUrlParam(r)
	resources := GetResByAppId(queryUrl["app_id"][0])
	ret := make(map[string]interface{})
	ret["resources"] = resources
	common.Write(w, ret)
}

//通过res_name查询资源记录
func (resource *Resource) QueryResByResName(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	queryUrl := common.GetUrlParam(r)
	res := GetResByResName(queryUrl["res_name"][0])
	ret := make(map[string]interface{})
	ret["resource"] = res
	common.Write(w, ret)
}

//增加资源记录
func (resource *Resource) AddResource(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	queryUrl := common.GetUrlParam(r)
	var res Resource
	res.App_id = queryUrl["app_id"][0]
	res.Res_name = queryUrl["res_name"][0]
	res.Res_cname = queryUrl["res_cname"][0]
	resType, err := strconv.Atoi(queryUrl["res_type"][0])
	if err != nil {
		fmt.Println("resType convert int failure")
	}
	res.Res_type = resType
	res.Res_target = queryUrl["res_target"][0]
	res.Res_desc = queryUrl["res_desc"][0]
	status, err := strconv.Atoi(queryUrl["status"][0])
	if err != nil {
		fmt.Println("status convert int failure")
	}
	res.Status = status
	result := InsertResource(res)
	ret := make(map[string]interface{})
	ret["result"] = result
	common.Write(w, ret)
}

//修改资源状态
func (resource *Resource) ModifyResourceStatus(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	queryUrl := common.GetUrlParam(r)
	status, err := strconv.Atoi(queryUrl["status"][0])
	if err != nil {
		fmt.Println("status convert int failure")
	}
	resId, err2 := strconv.Atoi(queryUrl["res_id"][0])
	if err2 != nil {
		fmt.Println("res_id convert int failure")
	}
	result := UpdateResStatus(status, resId)
	ret := make(map[string]interface{})
	ret["result"] = result
	common.Write(w, ret)
}
