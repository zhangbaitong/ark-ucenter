package main

import (
	"fmt"
	"common"
	"gopkg.in/mgo.v2/bson" 
)
type Response struct {
	Code int
	Message string
}

type ATUserInfo struct {
      Id bson.ObjectId "_id"
	Ac_id   int
	Info map[string] string
}

func main() {

	session := common.GetSession()
	if(session==nil){
		return 
	}
	defer common.FreeSession(session)

	result := ATUserInfo{}
	coll := session.DB("at_db").C("user_tab")

	err:=coll.Find(&bson.M{"_id": bson.ObjectIdHex("55c1804be13823298d000001")}).One(&result)
	if err!=nil {
	fmt.Println(err)	
		return 
	}
	fmt.Println(result)	
	return 
}
