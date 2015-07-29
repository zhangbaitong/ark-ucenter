package action

import (
	"common"
	"fmt"
	"gopkg.in/mgo.v2/bson" 
	_"encoding/json"
	"strings"
)

type ATUserData struct {
		Ac_name   string
		Ac_id   int
		Status   int
		Source   int
		Create_time   int
}

type ATUserInfo struct {
	Ac_id   int
	Info map[string] string
}

type Account struct {
	Acid        int
	Ac_name     string
	Ac_password string
	Status      int
}

type User struct {
	Acname   string
	Password string
}

const (
	INSERT string = "insert into account_tab (ac_name,ac_password,status,create_time) values (?,?,?,unix_timestamp())"
)

//登录插入
func LoginQuery(user *User) bool {
	strSQL := fmt.Sprintf("select count(ac_name) from account_tab where ac_name='%s' and ac_password='%s'", user.Acname, user.Password)
	rows, err := common.GetDB().Query(strSQL)
	if err != nil {
		return false
	} else {
		defer rows.Close()
		var nCount int
		for rows.Next() {
			rows.Scan(&nCount)
		}
		if nCount == 0 {
			return false
		}
	}
	return true
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

func GetUser(ac_name string) (UserData* ATUserData,ok bool) {
	UserData=&ATUserData{}
	strSQL := fmt.Sprintf("select ac_name,ac_id,status,source,create_time from account_tab where ac_name='%s' ", ac_name)
	mydb := common.GetDB()
	if mydb == nil {
		fmt.Println("get db connection error")
		return UserData,false
	}
	defer common.FreeDB(mydb)

	rows, err := mydb.Query(strSQL)
	if err != nil {
		return UserData,false
	} else {
		defer rows.Close()
		if rows.Next() {
			rows.Scan(&UserData.Ac_name,&UserData.Ac_id,&UserData.Status,&UserData.Source,&UserData.Create_time)
		} else {
			return UserData,false
		}
	}
	return UserData,true
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
	ok:=action.LoginMulti("26343637","222222")
	if ok {
		fmt.Println("success")
	} else {

		fmt.Println("faild")
	}
	return;
*/
func LoginMulti(strName,strPassword string) (strACName string,ok bool) {
	FieldList:=GetSearchFieldes();
	condition := make([]bson.M, len(FieldList))
	var strTemp string
	for i := 0; i < len(FieldList); i++ {
		strTemp="info."+strings.ToLower(FieldList[i])
		condition[i]=bson.M{strTemp:strName}
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
		if m.Ac_id==-1{
			continue
		}
		strACName,ok:=LoginById(m.Ac_id,strPassword)
		if ok {
			return strACName,true
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
