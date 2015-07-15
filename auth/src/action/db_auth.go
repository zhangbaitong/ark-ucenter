package action

import (
	"common"
	"fmt"
)

type Res struct {
	Res_id        string
	Res_name      string
	Owner_acid    int
	Operator_acid int
	Status        int
	Create_time   int
}

//通过resId获取资源信息
func GetResStatus(resId string) (status int) {
	mydb := common.GetDB()
	if mydb == nil {
		fmt.Println("get db connection error")
		return -1
	}
	defer common.FreeDB(mydb)

	sqlStr := "select status from resource_tab where res_id=?"
	rows, err := mydb.Query(sqlStr, resId)
	if err != nil {
		fmt.Println("query status failure", err)
		return -1
	} else {
		rows.Next()
		rows.Scan(&status)
		return status
	}
}

//通过openId查询acId
func GetAcIdByresIdAndopenId(openId string) (acId int) {
	mydb := common.GetDB()
	if mydb == nil {
		fmt.Println("get db connection error")
		return -1
	}
	defer common.FreeDB(mydb)

	sqlStr := "select acid from openid_tab where openid=? "
	rows, err := mydb.Query(sqlStr, openId)
	if err != nil {
		fmt.Println("query acid failure")
		return -1
	} else {
		rows.Next()
		rows.Scan(&acId)
		return acId
	}
}

//通过res_id和acid查询权限
func GetPriviliges(resId string, acId int) (priviliges string) {
	mydb := common.GetDB()
	if mydb == nil {
		fmt.Println("get db connection error")
		return ""
	}
	defer common.FreeDB(mydb)

	sqlStr := "select priviliges from authorization_tab where res_id=? and acid=? "
	rows, err := mydb.Query(sqlStr, resId, acId)
	if err != nil {
		fmt.Println("query acid failure")
		return ""
	} else {
		rows.Next()
		rows.Scan(&priviliges)
		return priviliges
	}
}
