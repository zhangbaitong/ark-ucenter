package action

import (
	"database/sql"
	"fmt"
	_ "github.com/Go-SQL-Driver/MySQL"
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

func initdb() (*sql.DB, error) {
	return sql.Open("mysql", "root:111111@tcp(117.78.19.76:3306)/at_db")
}

//通过resId获取资源信息
func GetRes(resId string) *Res {
	db, err := initdb()
	if err != nil {
		fmt.Println("连接数据库失败")
	}
	defer db.Close()

	sqlStr := "select * from resource_tab where res_id=? limit 1"
	rows, err := db.Query(sqlStr, resId)
	if err != nil {
		return nil
	} else {
		var r Res
		rows.Next()
		rows.Scan(&r.Res_id, &r.Res_name, &r.Owner_acid, &r.Operator_acid, &r.Interface_url, &r.Interface_type, &r.Status, &r.Create_time)
		return &r
	}
}
