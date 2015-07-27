package action

import (
	"common"
	"fmt"
)

type App struct {
	App_id      string
	App_key     string
	App_name    string
	App_desc    string
	Domain      string
	Status      int
	Create_time int
}

//通过app_id获取app_key
func GetAppKey(appId string) (appKey string) {
	mydb := common.GetDB()
	if mydb == nil {
		fmt.Println("get db connection error")
		return ""
	}
	defer common.FreeDB(mydb)

	sqlStr := "select app_key from app_info_tab where app_id=?"
	rows, err := mydb.Query(sqlStr, appId)
	if err != nil {
		fmt.Println("query openid failure", err)
		return ""
	} else {
		rows.Next()
		rows.Scan(&appKey)
		return appKey
	}
}
