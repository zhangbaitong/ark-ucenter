package action

import (
	"common"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func NewResource() *Resource {
	return new(Resource)
}

func (resource *Resource) QueryResId(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	queryUrl := common.GetUrlParam(r)
	resName := GetResId(queryUrl["res_name"][0])
	strBody := []byte("{\"resName\":" + string(resName) + "}")
	w.Write(strBody)
}
