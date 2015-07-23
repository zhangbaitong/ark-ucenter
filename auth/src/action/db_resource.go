package action

import (
	"common"
	"fmt"
)

type Resource struct {
	Res_id      int
	App_id      string
	Res_name    string
	Res_cname   string
	Res_type    int
	Res_target  string
	Res_desc    string
	Status      int
	Create_time int
}

//通过res_name查询res_id
func GetResId(resName string) (resId int) {
	mydb := common.GetDB()
	if mydb == nil {
		fmt.Println("get db connection error")
		return 0
	}
	defer common.FreeDB(mydb)

	sqlStr := "select res_id from resource_tab where res_name=?"
	rows, err := mydb.Query(sqlStr, resName)
	if err != nil {
		fmt.Println("query res_id failure", err)
		return 0
	} else {
		rows.Next()
		rows.Scan(&resId)
		return resId
	}
}

//通过res_name查询res_cname
func GetResCname(resName string) (resCname string) {
	mydb := common.GetDB()
	if mydb == nil {
		fmt.Println("get db connection error")
		return ""
	}
	defer common.FreeDB(mydb)

	sqlStr := "select res_cname from resource_tab where res_name=?"
	rows, err := mydb.Query(sqlStr, resName)
	if err != nil {
		fmt.Println("query res_id failure", err)
		return ""
	} else {
		rows.Next()
		rows.Scan(&resCname)
		return resCname
	}
}

//通过app_id查询资源记录
func GetResByAppId(appId string) []Resource {
	mydb := common.GetDB()
	if mydb == nil {
		fmt.Println("get db connection error")
		return nil
	}
	defer common.FreeDB(mydb)

	sqlStr := "select res_id,app_id,res_name,res_cname,res_type,res_target,res_desc,status,create_time from resource_tab where app_id=?"
	rows, err := mydb.Query(sqlStr, appId)

	var resList []Resource = make([]Resource, 0)

	if err != nil {
		fmt.Println("query res_id failure", err)
		return nil
	} else {
		for rows.Next() {
			var res Resource
			rows.Scan(&res.Res_id, &res.App_id, &res.Res_name, &res.Res_cname, &res.Res_type, &res.Res_target, &res.Res_desc, &res.Status, &res.Create_time)
			resList = append(resList, res)
		}
	}
	return resList
}

//通过res_name查询资源记录
func GetResByResName(resName string) []Resource {
	mydb := common.GetDB()
	if mydb == nil {
		fmt.Println("get db connection error")
		return nil
	}
	defer common.FreeDB(mydb)

	sqlStr := "select res_id,app_id,res_name,res_cname,res_type,res_target,res_desc,status,create_time from resource_tab where res_name=?"
	rows, err := mydb.Query(sqlStr, resName)

	var resList []Resource = make([]Resource, 0)

	if err != nil {
		fmt.Println("query res_id failure", err)
		return nil
	} else {
		for rows.Next() {
			var res Resource
			rows.Scan(&res.Res_id, &res.App_id, &res.Res_name, &res.Res_cname, &res.Res_type, &res.Res_target, &res.Res_desc, &res.Status, &res.Create_time)
			resList = append(resList, res)
		}
	}
	return resList
}

//增加资源记录
func InsertResource(res Resource) bool {
	mydb := common.GetDB()
	if mydb == nil {
		fmt.Println("get db connection error")
		return false
	}
	defer common.FreeDB(mydb)

	sqlStr := "insert into resource_tab(app_id, res_name, res_cname, res_type, res_target, res_desc, status, create_time) values(?,?,?,?,?,?,?,unix_timestamp())"
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

	_, err = stmt.Exec(res.App_id, res.Res_name, res.Res_cname, res.Res_type, res.Res_target, res.Res_desc, res.Status)

	if err != nil {
		fmt.Println(err)
		return false
	}
	tx.Commit()
	return true
}

//修改资源状态
func UpdateResStatus(status int, resId int) bool {
	mydb := common.GetDB()
	if mydb == nil {
		fmt.Println("get db connection error")
		return false
	}
	defer common.FreeDB(mydb)

	sqlStr := "update resource_tab set status=? where res_id=?"
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

	_, err = stmt.Exec(status, resId)

	if err != nil {
		fmt.Println(err)
		return false
	}
	tx.Commit()
	return true
}
