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
