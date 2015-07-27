package action

import (
	"common"
	"fmt"
	 "gopkg.in/mgo.v2/bson" 
	"strings"
)

type PersonRes struct {
	App_id      string
	Openid      string
	Res_id      int
	Status      int
	Create_time int
}

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

//通过openId获取用户资源权限列表
func GetPersonResList(openId string) []PersonRes {
	mydb := common.GetDB()
	if mydb == nil {
		fmt.Println("get db connection error")
		return nil
	}
	defer common.FreeDB(mydb)

	sqlStr := "select app_id,openid,res_id,status,create_time from app_confered_person_tab where openId=?"
	rows, err := mydb.Query(sqlStr, openId)

	var personResList []PersonRes = make([]PersonRes, 0)

	if err != nil {
		fmt.Println("query res_id failure", err)
		return nil
	} else {
		for rows.Next() {
			var personRes PersonRes
			rows.Scan(&personRes.App_id, &personRes.Openid, &personRes.Res_id, &personRes.Status, &personRes.Create_time)
			personResList = append(personResList, personRes)
		}
	}
	return personResList
}


//查询账户是否存在
func isUserExist_i(name string, value int) (UserInfo* ATUserInfo,ok bool) {
	session := common.GetSession()
	if(session==nil){
		return nil,false
	}	
	defer common.FreeSession(session)

	result := ATUserInfo{}
	coll := session.DB("at_db").C("user_tab")

	err:=coll.Find(&bson.M{name:value}).Sort(name).One(&result)
	if err!=nil {
		return nil,false;
	}

	return &result,true
}

func isUserExist(name, value string) (UserInfo* ATUserInfo,ok bool) {
	session := common.GetSession()
	if(session==nil){
		return nil,false
	}	
	defer common.FreeSession(session)

	result := ATUserInfo{}
	coll := session.DB("at_db").C("user_tab")

	err:=coll.Find(&bson.M{name:value}).Sort(name).One(&result)
	if err!=nil {
		return nil,false;
	}

	return &result,true
}

func UpdateUserInfo(UserInfo* ATUserInfo) (ok bool) {
	session := common.GetSession()
	if(session==nil){
		return false
	}	
	defer common.FreeSession(session)

	old_info,ok:=isUserExist_i("ac_id",UserInfo.Ac_id)
	coll := session.DB("at_db").C("user_tab")	
	if ok {
		for k, v := range UserInfo.Info {
			old_info.Info[k]=v
		}		
		condition:=bson.M{"ac_id":UserInfo.Ac_id}
		err := coll.Update(condition, bson.M{"$set": bson.M{"info": old_info.Info}})
		if(err==nil){
			return true
		}
		return false
	}
	
	err:=coll.Insert(UserInfo)
	if(err!=nil){		
		return false
	}
	return true
}

func LoginById(ac_id int,strPassword string) (strName string,ok bool) {
	strSQL := fmt.Sprintf("select ac_name from account_tab where ac_id=%d and ac_password='%s'", ac_id,strPassword)
	mydb := common.GetDB()
	if mydb == nil {
		fmt.Println("get db connection error")
		return "",false
	}
	defer common.FreeDB(mydb)

	rows, err := mydb.Query(strSQL)
	if err != nil {
		return "",false
	} else {
		defer rows.Close()
		for rows.Next() {
			rows.Scan(&strName)
		}
		if len(strName) == 0 {
			return "",false
		}
	}
	return strName,false
}
/*
	ok:=action.LoginMulti("qq,mobile,email","26343637","222222")
	if ok {
		fmt.Println("success")
	} else {

		fmt.Println("faild")
	}
	return;
*/
func LoginMulti(strFieldList,strValue,strPassword string) (strName string,ok bool) {
	FieldList:=strings.Split(strFieldList, ",")
	//fmt.Println("Count=",len(FieldList))
	condition := make([]bson.M, len(FieldList))
	var strTemp string
	for i := 0; i < len(FieldList); i++ {
		strTemp="info."+strings.ToLower(FieldList[i])
		condition[i]=bson.M{strTemp:strValue}
	}
	condition_or:=bson.M{"$or":condition}

	session := common.GetSession()
	if(session==nil){
		return "",false
	}	
	defer common.FreeSession(session)

	result := []ATUserInfo{}
	coll := session.DB("at_db").C("user_tab")
	err:=coll.Find(condition_or).Sort("ac_id").All(&result)
	if(err!=nil){
		fmt.Println("faild")
		return "",false
	}
	
	for _, m := range result {
		strName,ok:=LoginById(m.Ac_id,strPassword)
		if ok {
			return strName,true
		}
	}
	return "",false
}

func MultiRegister(strName,strValue string) (ok bool){
	session := common.GetSession()
	if(session==nil){
		return false
	}	
	defer common.FreeSession(session)

	coll := session.DB("at_db").C("user_tab")
	strNode:="info."+strName
	_,OK:=isUserExist(strNode,strValue)
	if OK {
		return true
	}
	var UserInfo ATUserInfo
	UserInfo.Ac_id=-1;
	UserInfo.Info[strName]=strValue
	err:=coll.Insert(&UserInfo)
	if(err!=nil){		
		return false
	}
	return true
}

func RegisterInsert(ac *Account) (ok bool) {
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
