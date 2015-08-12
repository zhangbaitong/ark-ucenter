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
	      Mid string
		Create_time   int
}

type ATUserInfo struct {
      Id bson.ObjectId "_id"
	Ac_id   int
	Info map[string] string
}

type UserInfoResult struct {
      Id string
	Info map[string] string
}

type UserInfoAll struct {
      Id string
	Ac_name   string
	Status   int
	Source   int
	Create_time   int
	Info map[string] string
}

const (
	INSERT string = "insert into account_tab (ac_name,ac_password,status,mid,create_time) values (?,?,?,?,unix_timestamp())"
)
const (
	OK=iota
	GETDB_ERROR
	INSERT_DB_ERROR
	UPDATE_DB_ERROR
	USER_NOT_EX
	USER_EX
	LIST_EMPTY
	PARAM_ERROR
	UNKNOWN_ERROR
)
var (
	error_list=[...]string{"OK","get db connection error","insert db error","update db error",
	"user not exist or password error","user was existed","list empty","param error","unknown error"}
)

func GetError(code int) (strMessage string){
	if code >=len(error_list) {
		strMessage=error_list[UNKNOWN_ERROR]
		return
	}

	return error_list[code]
}

func MultiRegister(Info* map[string]string) (InfoResult UserInfoResult,code int){
	session := common.GetSession()
	if(session==nil){
		return InfoResult,1
	}	
	defer common.FreeSession(session)

	coll := session.DB("at_db").C("user_tab")
	User,OK:=isUserExist_m(Info)
	if OK {
		InfoResult.Id=User.Id.Hex()
		return InfoResult,USER_EX
	}

	var UserInfo ATUserInfo
	UserInfo.Ac_id=-1;
	UserInfo.Info=*Info
	UserInfo.Id=bson.NewObjectId()
	InfoResult.Info=*Info
	InfoResult.Id=UserInfo.Id.Hex()
	err:=coll.Insert(&UserInfo)
	if(err!=nil){		
		return InfoResult,INSERT_DB_ERROR
	}
	return InfoResult,0
}

func RegisterInsert(strACName,strPassword,strID string) (ok bool) {
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

	_, err = stmt.Exec(strACName, strPassword, 0,strID)
	if err != nil {
		fmt.Println(err)
		return false
	}
	tx.Commit()
	return true
}

//登录插入
func LoginQuery(strACName,strPassword string) bool {
	strSQL := fmt.Sprintf("select count(ac_name) from account_tab where ac_name='%s' and ac_password='%s'", strACName,strPassword)
	mydb := common.GetDB()
	if mydb == nil {
		fmt.Println("get db connection error3333")
		return false
	}
	defer common.FreeDB(mydb)

	rows, err := mydb.Query(strSQL)
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

func InsertUserInfo(UserInfo* ATUserInfo) (ok bool) {
	session := common.GetSession()
	if(session==nil){
		return false
	}	
	defer common.FreeSession(session)

	coll := session.DB("at_db").C("user_tab")	
	err:=coll.Insert(UserInfo)	
	if(err!=nil){		
		return false
	}
	return true
}

func UpdateUserInfo(UserInfo* ATUserInfo) (ok bool) {
	session := common.GetSession()
	if(session==nil){
		return false
	}	
	defer common.FreeSession(session)

	if UserInfo.Ac_id>0 {
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
		}
		return false
	}

	old_info,ok:=GetUserInfoM(UserInfo.Id.Hex())
	coll := session.DB("at_db").C("user_tab")	
	if ok {
		for k, v := range UserInfo.Info {
			old_info.Info[k]=v
		}		
		condition:=&bson.M{"_id": UserInfo.Id}
		err := coll.Update(condition, bson.M{"$set": bson.M{"info": old_info.Info}})
		if(err==nil){
			return true
		}
	}
	return false	
}

func GetUserInfoM(id string) (UserInfo ATUserInfo,ok bool) {
	session := common.GetSession()
	if(session==nil){
		return UserInfo,false
	}	
	defer common.FreeSession(session)

	coll := session.DB("at_db").C("user_tab")
	fmt.Println("id=",id)
	//err := coll.Find(&bson.M{"ac_id":ac_id}).Sort("ac_id").One(&UserInfo)
	err:=coll.Find(&bson.M{"_id": bson.ObjectIdHex(id)}).One(&UserInfo)
	if(err==nil){
		return UserInfo,true
	}
	fmt.Println(err)
	return UserInfo,false
}

func GetUser(ac_name string) (UserData* ATUserData,ok bool) {
	UserData=&ATUserData{}
	strSQL := fmt.Sprintf("select ac_name,ac_id,status,source,mid,create_time from account_tab where ac_name='%s' ", ac_name)
	mydb := common.GetDB()
	if mydb == nil {
		fmt.Println("get db connection error3333")
		return UserData,false
	}
	defer common.FreeDB(mydb)

	rows, err := mydb.Query(strSQL)
	if err != nil {
		return UserData,false
	} else {
		defer rows.Close()
		if rows.Next() {
			rows.Scan(&UserData.Ac_name,&UserData.Ac_id,&UserData.Status,&UserData.Source,&UserData.Mid,&UserData.Create_time)
		} else {
			return UserData,false
		}
	}
	return UserData,true
}

func GetUserById(id string) (UserData* ATUserData,ok bool) {
	UserData=&ATUserData{}
	strSQL := fmt.Sprintf("select ac_name,ac_id,status,source,mid,create_time from account_tab where mid='%s' ", id)
	mydb := common.GetDB()
	if mydb == nil {
		fmt.Println("get db connection error2222")
		return UserData,false
	}
	defer common.FreeDB(mydb)

	rows, err := mydb.Query(strSQL)
	if err != nil {
		return UserData,false
	} else {
		defer rows.Close()
		if rows.Next() {
			rows.Scan(&UserData.Ac_name,&UserData.Ac_id,&UserData.Status,&UserData.Source,&UserData.Mid,&UserData.Create_time)
		} else {
			return UserData,false
		}
	}
	return UserData,true
}

func GetUserByAcId(acid int) (UserData* ATUserData,ok bool) {
	UserData=&ATUserData{}
	strSQL := fmt.Sprintf("select ac_name,ac_id,status,source,mid,create_time from account_tab where ac_id=%d ", acid)
	mydb := common.GetDB()
	if mydb == nil {
		fmt.Println("get db connection error2222")
		return UserData,false
	}
	defer common.FreeDB(mydb)

	rows, err := mydb.Query(strSQL)
	if err != nil {
		return UserData,false
	} else {
		defer rows.Close()
		if rows.Next() {
			rows.Scan(&UserData.Ac_name,&UserData.Ac_id,&UserData.Status,&UserData.Source,&UserData.Mid,&UserData.Create_time)
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
		fmt.Println("get db connection error 1111")
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

func in_array(str string,strArray []string) (ok bool){
	for i := 0; i < len(strArray); i++ {
		if 	strArray[i]==str {
			return true
		}
	}
	return false
}
func isUserExist_m( Info* map[string]string) (UserInfo* ATUserInfo,ok bool) {

	FieldList:=GetSearchFieldes();
	//condition := make([]bson.M, 0)	
	var strTemp string
	nCount:=0
	for k, _ := range *Info {
		if in_array(strings.ToLower(k),FieldList) {
			nCount++;
		}
	}

	condition := make([]bson.M, nCount)
	nCount=0
	strNode:=""
	for i := 0; i < len(FieldList); i++ {
		strTemp=strings.ToLower(FieldList[i])
		if len((*Info)[strTemp])>0 {			
			strNode="info."+strTemp
			condition[nCount]=bson.M{strNode:(*Info)[strTemp]}
			nCount++
		}
	}

	condition_or:=bson.M{"$or":condition}
	session := common.GetSession()
	if(session==nil){
		return nil,false
	}	
	defer common.FreeSession(session)

	result := ATUserInfo{}
	coll := session.DB("at_db").C("user_tab")
	err:=coll.Find(condition_or).Sort("ac_id").One(&result)
	if err!=nil {
		return nil,false;
	}

	return &result,true
}
