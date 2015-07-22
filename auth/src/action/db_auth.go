package action

import (
	"common"
	"fmt"
)

//通过ac_id和app_id查询openid
func GetOpenId(ac_id int, app_id string) (openid string) {
	mydb := common.GetDB()
	if mydb == nil {
		fmt.Println("get db connection error")
		return ""
	}
	defer common.FreeDB(mydb)

	sqlStr := "select openid from openid_tab where ac_id=? and app_id=?"
	rows, err := mydb.Query(sqlStr, ac_id, app_id)
	if err != nil {
		fmt.Println("query openid failure", err)
		return ""
	} else {
		rows.Next()
		rows.Scan(&openid)
		return openid
	}

}

//通过app_id和res_name查询res_id
func GetResId(res_name string) (res_id int) {
	mydb := common.GetDB()
	if mydb == nil {
		fmt.Println("get db connection error")
		return 0
	}
	defer common.FreeDB(mydb)

	sqlStr := "select res_id from resource_tab where res_name=?"
	rows, err := mydb.Query(sqlStr, res_name)
	if err != nil {
		fmt.Println("query res_id failure", err)
		return 0
	} else {
		rows.Next()
		rows.Scan(&res_id)
		return res_id
	}
}

//通过app_id和res_name查询res_id
func GetResCname(res_name string) (res_cname string) {
	mydb := common.GetDB()
	if mydb == nil {
		fmt.Println("get db connection error")
		return ""
	}
	defer common.FreeDB(mydb)

	sqlStr := "select res_cname from resource_tab where res_name=?"
	rows, err := mydb.Query(sqlStr, res_name)
	if err != nil {
		fmt.Println("query res_id failure", err)
		return ""
	} else {
		rows.Next()
		rows.Scan(&res_cname)
		return res_cname
	}
}

//判断资源是否被授予给指定应用的指定用户
func IsPersonConfered(app_id string, openid string, res_id int) bool {
	mydb := common.GetDB()
	if mydb == nil {
		fmt.Println("get db connection error")
		return false
	}
	defer common.FreeDB(mydb)

	sqlStr := "select count(*) from app_confered_person_tab where app_id=? and openid=? and res_id=? and status=0"
	rows, err := mydb.Query(sqlStr, app_id, openid, res_id)
	if err != nil {
		fmt.Println("query IsPersonConfered failure", err)
		return false
	} else {
		num := 0
		rows.Next()
		rows.Scan(&num)
		if num > 0 {
			return true
		}
	}
	return false
}

//查询应用是否被授予访问指定资源的权限
func IsAppConfered(app_id string, res_id int) bool {
	mydb := common.GetDB()
	if mydb == nil {
		fmt.Println("get db connection error")
		return false
	}
	defer common.FreeDB(mydb)

	sqlStr := "select count(*) from app_confered_tab where app_id=? and res_id=? and status=0"
	rows, err := mydb.Query(sqlStr, app_id, res_id)
	if err != nil {
		fmt.Println("query IsAppConfered failure", err)
	} else {
		num := 0
		rows.Next()
		rows.Scan(&num)
		if num > 0 {
			return true
		}
	}
	return false
}

//增加“用户授权表”记录
func InsertPersonConfered(app_id string, openid string, res_id int) bool {
	mydb := common.GetDB()
	if mydb == nil {
		fmt.Println("get db connection error")
		return false
	}
	defer common.FreeDB(mydb)

	sqlStr := "insert into app_confered_person_tab(app_id, openid, res_id, status, create_time) values(?,?,?,?,unix_timestamp())"
	tx, err := mydb.Begin()
	if err != nil {
		fmt.Println(err)
		return false
	}
	stmt, err := tx.Prepare(sqlStr)
	if err != nil {
		fmt.Println(err)
		return false
	}
	defer stmt.Close()

	_, err = stmt.Exec(app_id, openid, res_id, 0)

	if err != nil {
		fmt.Println(err)
		return false
	}
	tx.Commit()
	return true
}

//
func GetAcId(acName string) int {
	strSQL := fmt.Sprintf("select ac_id from account_tab where ac_name='%s' limit 1", acName)
	mydb := common.GetDB()
	if mydb == nil {
		fmt.Println("get db connection error")
		return -1
	}
	defer common.FreeDB(mydb)
	rows, err := mydb.Query(strSQL)
	if err != nil {
		return -1
	} else {
		defer rows.Close()
		var acid int
		for rows.Next() {
			rows.Scan(&acid)
		}
		return acid
	}
	return -1
}

//通过用户名和应用名称获取openId
func GetOpenIdByacName(acName string, appId string) string {
	acId := GetAcId(acName)
	return GetOpenId(acId, appId)
}
