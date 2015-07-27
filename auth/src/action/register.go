package action

import (
	_"common"
	_"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

type (
	Register struct{}

	Account struct {
		Acid        int
		Ac_name     string
		Ac_password string
		Status      int
	}
)

const (
	INSERT string = "insert into account_tab (ac_name,ac_password,status,create_time) values (?,?,?,unix_timestamp())"
)

func NewRegister() *Register {
	return new(Register)
}

func (register *Register) Post(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	acname := req.FormValue("acname")
	password := req.FormValue("password")
	account := Account{Ac_name: acname, Ac_password: password}
	strBody := []byte("{\"Code\":0,\"Message\":\"ok\"}")	
	ok := RegisterInsert(&account)
	if !ok {
		strBody = []byte("{\"Code\":1,\"Message\":\"insert db faild!!\"}")
	}
	w.Write(strBody)
}

func (register *Register) RegisterMulti(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	req.ParseForm()

	strBody := []byte("{\"Code\":0,\"Message\":\"ok\"}")	
	for k, v := range req.Form {
		ok := MultiRegister(k,v[0])
		if !ok {
			strBody = []byte("{\"Code\":1,\"Message\":\"insert db faild!!\"}")
		}
		break
	}
	w.Write(strBody)
}
