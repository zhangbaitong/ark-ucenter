package action

import (
	"common"
	"fmt"
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
	ok := register_insert(&account)
	if ok {
		w.Write([]byte("0"))
	} else {
		w.Write([]byte("-1"))
	}
}

func register_insert(ac *Account) (ok bool) {
	mydb := common.GetDB()
	if mydb == nil {
		return false
	}
	defer common.FreeDB(mydb)

	tx, err := mydb.Begin()
	if err != nil {
		fmt.Println(err)
		return false
	}

	stmt, err := tx.Prepare(INSERT)
	if err != nil {
		fmt.Println(err)
		return false
	}
	defer stmt.Close()

	_, err = stmt.Exec(ac.Ac_name, ac.Ac_password, 0)
	if err != nil {
		fmt.Println(err)
		return false
	}
	tx.Commit()
	return true
}
