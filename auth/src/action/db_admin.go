package action

import (
	"common"
	"fmt"
	"gopkg.in/mgo.v2/bson" 
	"sync"
	_"encoding/json"
)

type Field_List struct {
	List []string
}
var (
	SearchFieldes Field_List
	Mu  sync.Mutex
)

func InitSearchFieldes() (ok bool){
	fmt.Println("InitSearchFieldes begin")
	Mu.Lock()
	defer Mu.Unlock()
	session := common.GetSession()
	if(session==nil){
		return false
	}	
	defer common.FreeSession(session)

	coll := session.DB("at_db").C("search_field_tab")
	err:=coll.Find(bson.M{}).One(&SearchFieldes)
	if err!=nil {
		return false;
	}

	return true
}

func GetSearchFieldes() ([]string){
	Mu.Lock()
	defer Mu.Unlock()
	return SearchFieldes.List;
}

func SetSearchFieldes(Fieldes []string) (ok bool) {	
	Mu.Lock()
	defer Mu.Unlock()
	SearchFieldes.List=Fieldes;	
	session := common.GetSession()
	if(session==nil){
		return false
	}	
	defer common.FreeSession(session)

	coll := session.DB("at_db").C("search_field_tab")	
	List:=Field_List{List:Fieldes}
	coll.Remove(bson.M{})
	err:=coll.Insert(List)
	if(err!=nil){		
		fmt.Println("SetSearchFieldes:",err)
		return false
	}
	return true
}
