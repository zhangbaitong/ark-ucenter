package action

import (
	"common"
	"fmt"
)

type Res struct {
	Res_id         string
	Res_name       string
	Owner_acid     int
	Operator_acid  int
	Interface_url  string
	Interface_type int
	Status         int
	Create_time    int
}

//通过resId获取资源信息
func GetRes(resId string) *Res {
	mydb := common.GetDB()
	if mydb == nil {
		fmt.Println("get db connection error")
		return nil
	}
	defer common.FreeDB(mydb)

	sqlStr := "select * from resource_tab where res_id=? limit 1"
	rows, err := mydb.Query(sqlStr, resId)
	if err != nil {
		return nil
	} else {
		var r Res
		rows.Next()
		rows.Scan(&r.Res_id, &r.Res_name, &r.Owner_acid, &r.Operator_acid, &r.Interface_url, &r.Interface_type, &r.Status, &r.Create_time)
		return &r
	}
}

//通过appId查询资源Id
//通过resId获取资源信息
func GetResIds(appId string) (resId string) {
	fmt.Println("345")
	mydb := common.GetDB()
	if mydb == nil {
		fmt.Println("get db connection error")
		return ""
	}
	defer common.FreeDB(mydb)

	sqlStr := "select res_id from client_secret_tab where app_id=? limit 1"
	rows, err := mydb.Query(sqlStr, appId)
	if err == nil {
		resId = ""
		for rows.Next() {
			rows.Scan(&resId)
		}
		return resId
	}
	return ""
}

func Test() {
	fmt.Println("test")
}
