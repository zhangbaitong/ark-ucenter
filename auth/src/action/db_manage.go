package action

import (
	"common"
	"fmt"
	"gopkg.in/mgo.v2/bson" 
	"sync"
	_"encoding/json"
)

type Field_List struct {
	Name string
	List []string
}

var (
	OnlyCheckList Field_List
	Mu  sync.Mutex
)

func DeleteUser(nACID int) (ok bool) {
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

	stmt, err := tx.Prepare("delete from account_tab where ac_id=?")
	if err != nil {
		fmt.Println(err)
		return false
	}
	defer stmt.Close()

	_, err = stmt.Exec(nACID)
	if err != nil {
		fmt.Println(err)
		return false
	}
	tx.Commit()
	return true
}

func InitOnlyCheckList() (ok bool){
	Mu.Lock()
	defer Mu.Unlock()
	session := common.GetSession()
	if(session==nil){
		return false
	}	
	defer common.FreeSession(session)

	coll := session.DB("at_db").C("dictionary_tab")
	err:=coll.Find(bson.M{"name": "check_list"}).One(&OnlyCheckList)
	if err!=nil {
		return false;
	}

	return true
}

func GetCheckList() ([]string){
	Mu.Lock()
	defer Mu.Unlock()
	return OnlyCheckList.List;
}

func SetOnlyCheckList(Fieldes []string) (ok bool) {	
	Mu.Lock()
	defer Mu.Unlock()
	OnlyCheckList.List=Fieldes;	
	session := common.GetSession()
	if(session==nil){
		return false
	}	
	defer common.FreeSession(session)

	strName:="check_list"
	coll := session.DB("at_db").C("dictionary_tab")	
	coll.Remove(bson.M{"name": strName})
	err:=coll.Insert(OnlyCheckList)
	if(err!=nil){		
		fmt.Println("SetOnlyCheckList:",err)
		return false
	}
	return true
}

func ResetPassword(strAcName string,strNewPwd string)(ok bool){
	mydb := common.GetDB()
	if(mydb==nil){
		return false
	}	
	defer common.FreeDB(mydb)

	tx, err := mydb.Begin()
	if err != nil {
		fmt.Println(err)
		return false
	}
	stmt, err := tx.Prepare(" UPDATE account_tab SET ac_password=? where ac_name=? ")
	if err != nil {
		fmt.Println(err)
		return false
	}
	defer stmt.Close()

	_, err = stmt.Exec(strNewPwd, strAcName)

	if err != nil {
		fmt.Println(err)
		return false
	}
	tx.Commit()
	return true
}

func ExportMongo( start_time,end_time int) (count int,ok bool) {
	session := common.GetSession()
	if(session==nil){
		return 0,false
	}	
	defer common.FreeSession(session)

	result := []ATUserInfo{}
	coll := session.DB("at_db").C("user_tab")
	//conditions := make([]bson.M, 2)

	condition:=bson.M{"info.name": bson.M{"$exists": true}}
	//conditions[0]=bson.M{"info.name": bson.M{"$exists": true}}
	//conditions[1]=bson.M{"create_time":bson.M{"$gte":start_time,"$lte":end_time}}
	err:=coll.Find(condition).Sort("create_time").All(&result)
	if err!=nil {
		fmt.Println(err)
		return 0,false;
	}

	nCount:=0;
	for i := 0; i < len(result); i++ {
		
		/*
		//delete from mongodb
		conditions :=bson.M{"info.phone":result[i].Info["phone"]}
		err = coll.Remove(conditions)
		if err!=nil {
			fmt.Println(err)
			return 0,false;
		}

		//delete from mysql
		if result[i].Ac_id>0 {
			DeleteUser(result[i].Ac_id)
		}
		*/
		if result[i].Create_time>start_time  && result[i].Create_time<end_time {
			nCount++;
		}
		fmt.Println(result[i].Id.Hex(),";",result[i].Info["name"])
		//fmt.Println(result[i].Info)
	}
	return nCount,true
/*
	//err:=coll.Find(&bson.M{"create_time":{"$gte":start_time,"$lte":end_time}}).Sort("create_time").All(&result)
	//err:=coll.Find(&bson.M{"create_time":bson.M{"$gte":start_time,"$lte":end_time}}).Sort("create_time").All(&result)
	condition := make([]bson.M, 6)
	condition[0]=bson.M{"info.user_type":"10"}
	condition[1]=bson.M{"info.user_type":"20"}
	condition[2]=bson.M{"info.user_type":"30"}
	condition[3]=bson.M{"info.user_type":"40"}
	condition[4]=bson.M{"info.user_type":"50"}
	condition[5]=bson.M{"info.user_type":"60"}
	//condition[6]=bson.M{"create_time":0}
	fmt.Println("start_time=",start_time,"end_time=",end_time)
	err:=coll.Find(&bson.M{"$or":condition}).Sort("create_time").All(&result)
	//err:=coll.Find(&bson.M{"create_time":bson.M{"$gte":start_time,"$lte":end_time}}).Sort("create_time").All(&result)
	if err!=nil {
		fmt.Println(err)
		return 0,false;
	}
	
	for i := 0; i < len(result); i++ {

		if result[i].Create_time==0  && result[i].Info["user_name"]!="" {
			strText:=fmt.Sprintf("%s,%s,%s,%s,%s",result[i].Id.Hex(),result[i].Info["user_name"],result[i].Info["company_name"],
				result[i].Info["job"],result[i].Info["mobile"])
			strText=strings.Trim(strText,"\n")
			strText=strings.Trim(strText,"\r")
			fmt.Println(strText)
			count++
		}
	}
	
	return count,true
*/	
}

func RegUserStat(start_time,end_time ,source int ) (count int,ok bool) {
	mydb := common.GetDB()
	if mydb == nil {
		return 0,false
	}
	defer common.FreeDB(mydb)

	strSQL := fmt.Sprintf("select count(ac_id) from account_tab where create_time>%d and create_time<=%d and  source=%d", start_time,end_time,source)
	fmt.Println(strSQL)
	rows, err := mydb.Query(strSQL)
	if err != nil {
		return 0,false
	} else {
		defer rows.Close()
		if rows.Next() {
			rows.Scan(&count)
		} else {
			return 0,false
		}
	}
	return count,true

}

